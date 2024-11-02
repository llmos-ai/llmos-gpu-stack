package gpudevice

import (
	"fmt"
	"reflect"
	"strconv"
	"strings"

	hapi "github.com/Project-HAMi/HAMi/pkg/api"
	"github.com/Project-HAMi/HAMi/pkg/device"
	ctlcorev1 "github.com/rancher/wrangler/v3/pkg/generated/controllers/core/v1"
	"github.com/sirupsen/logrus"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/labels"

	gpustackv1 "github.com/llmos-ai/llmos-gpu-stack/pkg/apis/gpustack.llmos.ai/v1"
	ctlgpustackv1 "github.com/llmos-ai/llmos-gpu-stack/pkg/generated/controllers/gpustack.llmos.ai/v1"
)

type nodeHandler struct {
	nodeClient      ctlcorev1.NodeClient
	nodeCache       ctlcorev1.NodeCache
	gpuDevices      ctlgpustackv1.GPUDeviceClient
	gpuDeviceCache  ctlgpustackv1.GPUDeviceCache
	devices         map[string]device.Devices
	nodeDeviceCache *NodeDeviceThreadSafeCache
}

// nodeGPUDevicesOnChange helps to reconcile the node gpu devices when node obj has changed
func (h *nodeHandler) nodeGPUDevicesOnChange(_ string, node *corev1.Node) (*corev1.Node, error) {
	if node == nil || node.DeletionTimestamp != nil || node.Annotations == nil {
		return nil, nil
	}

	if cacheNode, found := h.nodeDeviceCache.Get(node.Name); found {
		// Skip validate node handshake time
		cacheNode.Annotations[HamiNodeHandshakeAnnotation] = node.Annotations[HamiNodeHandshakeAnnotation]
		if reflect.DeepEqual(node.Annotations, cacheNode.Annotations) {
			logrus.Debugf("node %s devices has not changed, skip updating", node.Name)
			return nil, nil
		}
	}

	logrus.Debugf("node %s has changed, check GPU devices", node.Name)
	var gpuDevices = make([]*gpustackv1.GPUDevice, 0)
	var gpuDeviceLabels = make(map[string]string)

	if !h.hasGPUDevices(node) {
		logrus.Debugf("node %s has no gpu devices, check if need to clean up old device obj", node.Name)
		if err := h.cleanNodeNotReadyDevices(node, gpuDevices); err != nil {
			return node, err
		}

		return h.updateGPUNodeLabel(node, gpuDeviceLabels, false)
	}
	// Reconcile all device types
	for _, dt := range h.devices {
		nodeDevices, err := dt.GetNodeDevices(*node)
		if err != nil && strings.Contains(err.Error(), deviceAnnoNotFound) {
			continue
		} else if err != nil {
			return node, fmt.Errorf("get node devices error: %v", err)
		}

		hasDevices := strconv.FormatBool(len(nodeDevices) > 0)
		gpuDeviceLabels[getNodeDeviceNameLabelKey(dt.CommonWord())] = hasDevices

		for _, device := range nodeDevices {
			gpuDevice, err := h.reconcileNodeGPUDevice(device, node)
			if err != nil {
				return node, err
			}
			gpuDevices = append(gpuDevices, gpuDevice)
		}
	}

	if err := h.cleanNodeNotReadyDevices(node, gpuDevices); err != nil {
		return node, err
	}

	return h.updateGPUNodeLabel(node, gpuDeviceLabels, true)
}

func (h *nodeHandler) hasGPUDevices(node *corev1.Node) bool {
	for _, dt := range h.devices {
		if devices, err := dt.GetNodeDevices(*node); err == nil {
			if len(devices) > 0 {
				return true
			}
		}
	}
	return false
}

func (h *nodeHandler) cleanNodeNotReadyDevices(node *corev1.Node, gpuDevices []*gpustackv1.GPUDevice) error {
	deviceList, err := h.gpuDeviceCache.List(labels.SelectorFromSet(map[string]string{
		LabelNodeNameKey: node.Name,
	}))
	if err != nil {
		return fmt.Errorf("failed to list gpu devices: %s", err)
	}

	for _, device := range deviceList {
		found := false
		for _, gpuDevice := range gpuDevices {
			if device.Name == gpuDevice.Name {
				found = true
				break
			}
		}
		if !found {
			if err = h.gpuDevices.Delete(device.Name, &metav1.DeleteOptions{}); err != nil {
				return fmt.Errorf("failed to delete gpu device: %s, %s", device.Name, err)
			}
		}
	}
	return nil
}
func (h *nodeHandler) reconcileNodeGPUDevice(device *hapi.DeviceInfo,
	node *corev1.Node) (*gpustackv1.GPUDevice, error) {
	gpuDevice := constructGPUDevice(device, node)
	foundDevice, err := h.gpuDeviceCache.Get(gpuDevice.Name)
	if err != nil && !errors.IsNotFound(err) {
		return nil, fmt.Errorf("failed to get gpu device by id: %s, %s", device.ID, err)
	}

	if foundDevice == nil {
		foundDevice, err = h.gpuDevices.Create(gpuDevice)
		if err != nil {
			return nil, fmt.Errorf("failed to create gpu device: %s, %s", gpuDevice.Name, err)
		}
		logrus.Debugf("Created gpu device: %s", gpuDevice.Name)
	}

	toUpdate := foundDevice.DeepCopy()
	// Init new gpu device status
	if gpustackv1.DeviceInitialized.GetStatus(foundDevice) == "" {
		logrus.Debugf("initializing gpu device %s status", gpuDevice.Name)
		toUpdate.Status = gpuDevice.Status
		gpustackv1.DeviceInitialized.SetStatusBool(toUpdate, true)
		if _, err = h.gpuDevices.UpdateStatus(toUpdate); err != nil {
			return toUpdate, fmt.Errorf("failed to update gpu device status: %s, %s", toUpdate.Name, err)
		}
		return toUpdate, nil
	}

	if !reflect.DeepEqual(foundDevice.Status.GPUDeviceInfo, gpuDevice.Status.GPUDeviceInfo) {
		logrus.Debugf("updating gpu device %s info of node %s", gpuDevice.Name, node.Name)
		toUpdate.Status.GPUDeviceInfo = gpuDevice.Status.GPUDeviceInfo
		if _, err = h.gpuDevices.UpdateStatus(toUpdate); err != nil {
			return toUpdate, fmt.Errorf("update gpu device status error: %v", err)
		}
	}

	return toUpdate, nil
}

func (h *nodeHandler) nodeGPUDevicesOnRemove(_ string, node *corev1.Node) (*corev1.Node, error) {
	if node.DeletionTimestamp == nil {
		return nil, nil
	}

	gpuDevices, err := h.gpuDeviceCache.List(labels.SelectorFromSet(map[string]string{
		LabelNodeNameKey: node.Name,
	}))

	if err != nil {
		return nil, fmt.Errorf("failed to list gpu devices: %s", err)
	}

	for _, gpuDevice := range gpuDevices {
		if err = h.gpuDevices.Delete(gpuDevice.Name, &metav1.DeleteOptions{}); err != nil {
			return node, err
		}
	}

	return node, nil
}

func (h *nodeHandler) updateGPUNodeLabel(node *corev1.Node, deviceLabels map[string]string,
	hasGPUDevices bool) (*corev1.Node, error) {
	toUpdate := node.DeepCopy()
	toUpdate.Labels[LabelGPUNodeRoleKey] = strconv.FormatBool(hasGPUDevices)

	for k, v := range deviceLabels {
		logrus.Debugf("updating gpu node %s labels %s:%s", node.Name, k, v)
		toUpdate.Labels[k] = v
	}

	if !reflect.DeepEqual(toUpdate.Labels, node.Labels) {
		return h.nodeClient.Update(toUpdate)
	}

	h.nodeDeviceCache.Set(node.Name, NodeDeviceInfo{Annotations: node.Annotations})
	return node, nil
}
