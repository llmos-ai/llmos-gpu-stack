package gpudevice

import (
	"fmt"
	"strings"

	"github.com/sirupsen/logrus"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/utils/ptr"

	autils "github.com/llmos-ai/llmos-gpu-stack/pkg/accelerators/utils"
	gpustackv1 "github.com/llmos-ai/llmos-gpu-stack/pkg/apis/gpustack.llmos.ai/v1"
	"github.com/llmos-ai/llmos-gpu-stack/pkg/utils"
	"github.com/llmos-ai/llmos-gpu-stack/pkg/utils/condition"
)

type GPUDevice struct {
	PodName   string `json:"podName"`
	Index     int    `json:"index"`
	Vendor    string `json:"vendor,omitempty"`
	UUID      string `json:"uuid,omitempty"`
	UsedMem   int32  `json:"usedMem,omitempty"`
	UsedCores int32  `json:"usedCores,omitempty"`
}

const (
	LLMOSPrefix         = "llmos.ai"
	LabelGPUNodeRoleKey = LLMOSPrefix + "/gpu-node"

	GPUStackPrefix   = "gpustack.llmos.ai"
	LabelNodeNameKey = GPUStackPrefix + "/node-name"

	NodeHandshakeAnnotation = "volcano.sh/node-vgpu-handshake"
	AssignedNodeAnnotations = "volcano.sh/vgpu-node"
)

func constructGPUDevice(device *autils.DeviceInfo, node *corev1.Node) *gpustackv1.GPUDevice {
	var internalIp string
	for _, address := range node.Status.Addresses {
		if address.Type == corev1.NodeInternalIP {
			internalIp = address.Address
			break
		}
	}
	logrus.Debugf("construct gpu device %+v for node %s", ParseDeviceInfo(device), node.Name)
	return &gpustackv1.GPUDevice{
		ObjectMeta: metav1.ObjectMeta{
			Name: getDeviceName(device.Id),
			Labels: map[string]string{
				LabelNodeNameKey: node.Name,
			},
			OwnerReferences: []metav1.OwnerReference{
				*metav1.NewControllerRef(node, node.GroupVersionKind()),
			},
		},
		Status: gpustackv1.GPUDeviceStatus{
			GPUDeviceInfo: ParseDeviceInfo(device),
			NodeName:      node.Name,
			InternalIP:    internalIp,
			State:         getDeviceState(device, node),
		},
	}
}

func ParseDeviceInfo(devInfo *autils.DeviceInfo) gpustackv1.GPUDeviceInfo {
	return gpustackv1.GPUDeviceInfo{
		UUID:     devInfo.Id,
		Index:    ptr.To(devInfo.Index),
		Vendor:   devInfo.Type[:strings.IndexByte(devInfo.Type, '-')],
		DevName:  strings.TrimPrefix(devInfo.Type[strings.IndexByte(devInfo.Type, '-'):], "-"),
		MaxCount: devInfo.Count,
		VRAM:     devInfo.Devmem,
		Health:   true, //TODO, need to get the actual health status from the device
	}
}

func getDeviceState(device *autils.DeviceInfo, node *corev1.Node) string {
	nodeIsReady := utils.IsNodeReady(node)

	logrus.Debugf("check status: device %s is %v, node %s is %v", device.Id, device.Health, node.Name, nodeIsReady)
	if !nodeIsReady {
		return condition.StateOffline
	}
	if !device.Health {
		return condition.StateUnhealthy
	}
	if device.Health && nodeIsReady {
		return condition.StateHealthy
	}

	return condition.StatePending
}

func getDeviceName(name string) string {
	return strings.ToLower(name)
}

func getPodDeviceNameLabelKey(deviceId string) string {
	return fmt.Sprintf("%s/%s", GPUStackPrefix, getDeviceName(deviceId))
}

func getNodeDeviceNameLabelKey(commonName string) string {
	return fmt.Sprintf("%s/%s-node", GPUStackPrefix, strings.ToLower(commonName))
}

func hasVGPUDevice(pod *corev1.Pod) bool {
	return pod.Annotations[AssignedNodeAnnotations] != ""
}
