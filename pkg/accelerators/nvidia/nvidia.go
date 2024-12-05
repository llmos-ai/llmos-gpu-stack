package nvidia

import (
	"fmt"

	"github.com/sirupsen/logrus"
	corev1 "k8s.io/api/core/v1"

	"github.com/llmos-ai/llmos-gpu-stack/pkg/accelerators/common"
	"github.com/llmos-ai/llmos-gpu-stack/pkg/accelerators/utils"
)

const (
	Name       = "nvidia"
	CommonName = "gpu"

	gpuPresentLabel = "nvidia.com/gpu.present"

	NodeHandshakeAnno        = "volcano.sh/node-vgpu-handshake"
	NodeDeviceRegisteredAnno = "volcano.sh/node-vgpu-register"
	DeviceAssignedIDsAnno    = "volcano.sh/vgpu-ids-new"
)

type nvidiaAccelerator struct {
	acceleratorName       string
	commonName            string
	nodeHandshakeAnno     string
	nodeRegisteredAnno    string
	deviceAssignedIdsAnno string
}

func Configure(name string) common.Accelerator {
	return &nvidiaAccelerator{
		acceleratorName:       name,
		commonName:            CommonName,
		nodeHandshakeAnno:     NodeHandshakeAnno,
		nodeRegisteredAnno:    NodeDeviceRegisteredAnno,
		deviceAssignedIdsAnno: DeviceAssignedIDsAnno,
	}
}

func (a *nvidiaAccelerator) GetName() string {
	return a.acceleratorName
}

func (a *nvidiaAccelerator) GetCommonName() string {
	return a.commonName
}

func (a *nvidiaAccelerator) GetPodAssignedDevicesKey() string {
	return a.deviceAssignedIdsAnno
}

func (a *nvidiaAccelerator) GetNodeDevices(node corev1.Node) ([]*utils.DeviceInfo, error) {
	devEncoded, ok := node.Annotations[NodeDeviceRegisteredAnno]
	if !ok {
		logrus.Debugf("annos %s not found for node %s", NodeDeviceRegisteredAnno, node.Name)
		return []*utils.DeviceInfo{}, nil
	}

	nodeDevices, err := utils.DecodeNodeDevices(devEncoded)
	if err != nil {
		return []*utils.DeviceInfo{}, fmt.Errorf("failed to decode node %s devices annotation: %s", node.Name, devEncoded)
	}

	if len(nodeDevices) == 0 {
		logrus.Debugf("no gpu device found for node: %s, anno: %s", node.Name, devEncoded)
		return []*utils.DeviceInfo{}, fmt.Errorf("no gpu found on node")
	}

	return nodeDevices, nil
}

func (a *nvidiaAccelerator) HasGPUPresent(node *corev1.Node) bool {
	return node.Labels[gpuPresentLabel] == "true"
}
