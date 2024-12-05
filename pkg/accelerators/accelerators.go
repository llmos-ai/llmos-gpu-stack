package accelerators

import (
	"fmt"
	"sync"

	"github.com/llmos-ai/llmos-gpu-stack/pkg/accelerators/common"
	"github.com/llmos-ai/llmos-gpu-stack/pkg/accelerators/nvidia"
)

var (
	Accelerators                = make(map[string]common.Accelerator)
	AcceleratorNames            = make(map[string]bool)
	AcceleratorDevicesCheckList = make(map[string]string)
	confMu                      sync.Mutex
)

func GetAccelerators() map[string]common.Accelerator {
	return Accelerators
}

func GetAcceleratorDevicesCheckList() map[string]string {
	return AcceleratorDevicesCheckList
}

func GetAccelerator(name string) (common.Accelerator, error) {
	if accelerator, ok := Accelerators[name]; ok {
		if accelerator != nil {
			return accelerator, nil
		}
	}
	return nil, fmt.Errorf("accelerator %s not found", name)
}

func Configure() {
	confMu.Lock()
	defer confMu.Unlock()

	n := nvidia.Configure(nvidia.Name)
	Accelerators[nvidia.Name] = n
	AcceleratorNames[nvidia.Name] = true
	AcceleratorDevicesCheckList[nvidia.Name] = n.GetPodAssignedDevicesKey()
}
