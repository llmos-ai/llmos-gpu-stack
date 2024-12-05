package common

import (
	corev1 "k8s.io/api/core/v1"

	"github.com/llmos-ai/llmos-gpu-stack/pkg/accelerators/utils"
)

type Accelerator interface {
	GetName() string
	GetCommonName() string
	GetPodAssignedDevicesKey() string
	GetNodeDevices(node corev1.Node) ([]*utils.DeviceInfo, error)
	HasGPUPresent(node *corev1.Node) bool
}
