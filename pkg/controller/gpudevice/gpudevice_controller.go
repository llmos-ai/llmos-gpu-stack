package gpudevice

import (
	"context"
	"fmt"
	"reflect"

	ctlcorev1 "github.com/rancher/wrangler/v3/pkg/generated/controllers/core/v1"
	"k8s.io/apimachinery/pkg/labels"

	"github.com/llmos-ai/llmos-gpu-stack/pkg/accelerators"
	gpustackv1 "github.com/llmos-ai/llmos-gpu-stack/pkg/apis/gpustack.llmos.ai/v1"
	"github.com/llmos-ai/llmos-gpu-stack/pkg/config"
	ctlgpustackv1 "github.com/llmos-ai/llmos-gpu-stack/pkg/generated/controllers/gpustack.llmos.ai/v1"
	"github.com/llmos-ai/llmos-gpu-stack/pkg/utils"
)

const (
	gpuDeviceNodeOnChange = "gpuDevice.nodeOnChange"
	gpuDeviceOnChange     = "gpuDevice.onChange"
	gpuDeviceNodeOnDelete = "gpuDevice.nodeOnDelete"
	gpuDevicePodOnChange  = "gpuDevice.PodOnChange"
)

type gpuHandler struct {
	gpuDevices           ctlgpustackv1.GPUDeviceClient
	gpuDeviceCache       ctlgpustackv1.GPUDeviceCache
	podCache             ctlcorev1.PodCache
	acceleratorChecklist map[string]string
}

func Register(_ context.Context, mgmt *config.Management) error {
	gpuDevices := mgmt.GPUStackFactory.Gpustack().V1().GPUDevice()
	nodes := mgmt.CoreFactory.Core().V1().Node()
	pods := mgmt.CoreFactory.Core().V1().Pod()
	acceleratorCheckList := accelerators.GetAcceleratorDevicesCheckList()

	nodeHandler := &nodeHandler{
		nodeClient:      nodes,
		nodeCache:       nodes.Cache(),
		gpuDevices:      gpuDevices,
		gpuDeviceCache:  gpuDevices.Cache(),
		accelerators:    mgmt.Accelerators,
		nodeDeviceCache: NewThreadSafeCache(),
	}
	nodes.OnChange(mgmt.Ctx, gpuDeviceNodeOnChange, nodeHandler.nodeGPUDevicesOnChange)
	nodes.OnRemove(mgmt.Ctx, gpuDeviceNodeOnDelete, nodeHandler.nodeGPUDevicesOnRemove)

	gpuHandler := &gpuHandler{
		gpuDevices:           gpuDevices,
		gpuDeviceCache:       gpuDevices.Cache(),
		podCache:             pods.Cache(),
		acceleratorChecklist: acceleratorCheckList,
	}
	gpuDevices.OnChange(mgmt.Ctx, gpuDeviceOnChange, gpuHandler.gpuDeviceOnChange)

	podHandler := &podHandler{
		gpuDevices:           gpuDevices,
		gpuDeviceController:  gpuDevices,
		gpuDeviceCache:       gpuDevices.Cache(),
		pods:                 pods,
		acceleratorCheckList: acceleratorCheckList,
	}

	pods.OnChange(mgmt.Ctx, gpuDevicePodOnChange, podHandler.onGpuPodChange)

	return nil
}

func (h *gpuHandler) gpuDeviceOnChange(_ string, gpuDevice *gpustackv1.GPUDevice) (*gpustackv1.GPUDevice, error) {
	if gpuDevice == nil || gpuDevice.DeletionTimestamp != nil {
		return nil, nil
	}

	selector := labels.SelectorFromSet(map[string]string{
		getPodDeviceNameLabelKey(gpuDevice.Name): "true",
	})
	pods, err := h.podCache.List("", selector)
	if err != nil {
		return gpuDevice, fmt.Errorf("failed to list GPU device %s pods, error: %s", gpuDevice.Name, err)
	}

	if len(pods) == 0 {
		return gpuDevice, nil
	}
	var deviceList []GPUDevice
	for _, pod := range pods {
		// Do not count pods that are being deleted
		if pod.DeletionTimestamp != nil {
			continue
		}

		devices, err := getPodAllocatedDevices(h.acceleratorChecklist, pod)
		if err != nil {
			return gpuDevice, err
		}
		deviceList = append(deviceList, devices...)
	}

	var podList = make([]gpustackv1.GPUPod, 0)

	for _, device := range deviceList {
		podList = append(podList, gpustackv1.GPUPod{
			Name:             device.PodName,
			CoresReq:         device.UsedCores,
			MemReq:           device.UsedMem,
			MemPercentageReq: utils.RoundToInt((float64(device.UsedMem)/float64(gpuDevice.Status.VRAM))*100, 2),
		})
	}

	toUpdate := gpuDevice.DeepCopy()
	toUpdate.Status.Pods = podList
	if !reflect.DeepEqual(gpuDevice.Status, toUpdate.Status) {
		return h.gpuDevices.UpdateStatus(toUpdate)
	}

	return gpuDevice, nil
}
