package gpudevice

import (
	"fmt"
	"reflect"

	ctlcorev1 "github.com/rancher/wrangler/v3/pkg/generated/controllers/core/v1"
	"github.com/sirupsen/logrus"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/errors"

	"github.com/llmos-ai/llmos-gpu-stack/pkg/accelerators/utils"
	ctlgpustackv1 "github.com/llmos-ai/llmos-gpu-stack/pkg/generated/controllers/gpustack.llmos.ai/v1"
)

type podHandler struct {
	pods                 ctlcorev1.PodClient
	gpuDevices           ctlgpustackv1.GPUDeviceClient
	gpuDeviceController  ctlgpustackv1.GPUDeviceController
	gpuDeviceCache       ctlgpustackv1.GPUDeviceCache
	acceleratorCheckList map[string]string
}

func (h *podHandler) onGpuPodChange(_ string, pod *corev1.Pod) (*corev1.Pod, error) {
	if pod == nil {
		return nil, nil
	}

	// Note: handle gpu pods syncing in updates to avoid finalizer blocking on all other pods
	if pod.DeletionTimestamp != nil {
		return h.onGpuPodDelete(pod)
	}

	if !hasVGPUDevice(pod) {
		return nil, nil
	}

	devices, err := getPodAllocatedDevices(h.acceleratorCheckList, pod)
	if err != nil {
		return pod, fmt.Errorf("get pod gpu devices error: %v", err)
	}

	if len(devices) == 0 {
		return pod, nil
	}

	if err = h.syncGPUDeviceStatus(devices); err != nil {
		return pod, nil
	}

	toUpdate := pod.DeepCopy()
	if toUpdate.Labels == nil {
		toUpdate.Labels = make(map[string]string)
	}

	for _, dev := range devices {
		devNameLabel := fmt.Sprintf("%s/%s", GPUStackPrefix, getDeviceName(dev.UUID))
		toUpdate.Labels[devNameLabel] = "true"
	}

	if !reflect.DeepEqual(toUpdate.Labels, pod.Labels) {
		return h.pods.Update(toUpdate)
	}

	return pod, nil
}

func (h *podHandler) syncGPUDeviceStatus(devices []GPUDevice) error {
	for _, dev := range devices {
		deviceName := getDeviceName(dev.UUID)
		_, err := h.gpuDeviceCache.Get(deviceName)
		if err != nil && errors.IsNotFound(err) {
			logrus.Warnf("gpu device %s not found, skip syncing device status", deviceName)
			h.gpuDeviceController.Enqueue(deviceName)
			continue
		} else if err != nil {
			return fmt.Errorf("get gpu device %s error: %v", deviceName, err)
		}

		h.gpuDeviceController.Enqueue(deviceName)
	}
	return nil
}

func getPodAllocatedDevices(checkList map[string]string, pod *corev1.Pod) ([]GPUDevice, error) {
	var gpuDevices []GPUDevice
	podDevices, err := utils.DecodePodDevices(checkList, pod.Annotations)
	if err != nil {
		return nil, err
	}

	for _, pDevice := range podDevices {
		for _, cDevices := range pDevice {
			for _, dev := range cDevices {
				gpuDevices = append(gpuDevices, GPUDevice{
					PodName:   getPodNamespaceName(pod),
					Vendor:    dev.Type,
					UUID:      dev.UUID,
					UsedMem:   dev.Usedmem,
					UsedCores: dev.Usedcores,
				})
			}
		}
	}

	return gpuDevices, nil
}

func (h *podHandler) onGpuPodDelete(pod *corev1.Pod) (*corev1.Pod, error) {
	if pod == nil || pod.DeletionTimestamp == nil || !hasVGPUDevice(pod) {
		return nil, nil
	}

	devices, err := getPodAllocatedDevices(h.acceleratorCheckList, pod)
	if err != nil {
		return pod, fmt.Errorf("get pod gpu devices error: %v", err)
	}

	if len(devices) == 0 {
		return pod, nil
	}

	for _, dev := range devices {
		deviceName := getDeviceName(dev.UUID)
		gpuDevice, err := h.gpuDeviceCache.Get(deviceName)
		if err != nil && errors.IsNotFound(err) {
			logrus.Warnf("gpu device %s not found, skip syncing device status", deviceName)
			h.gpuDeviceController.Enqueue(deviceName)
			continue
		} else if err != nil {
			return pod, fmt.Errorf("get gpu device %s error: %v", dev.UUID, err)
		}

		toUpdate := gpuDevice.DeepCopy()
		// remove the pod from the device pods status
		for i, p := range toUpdate.Status.Pods {
			if p.Name == getPodNamespaceName(pod) {
				toUpdate.Status.Pods = append(toUpdate.Status.Pods[:i], toUpdate.Status.Pods[i+1:]...)
				break
			}
		}
		if _, err = h.gpuDevices.UpdateStatus(toUpdate); err != nil {
			return pod, fmt.Errorf("update GPU device status error: %v", err)
		}
	}

	return nil, nil
}

func getPodNamespaceName(pod *corev1.Pod) string {
	return fmt.Sprintf("%s/%s", pod.Namespace, pod.Name)
}
