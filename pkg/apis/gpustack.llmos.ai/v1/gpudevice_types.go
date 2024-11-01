package v1

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	"github.com/llmos-ai/llmos-gpu-stack/pkg/apis/common"
	"github.com/llmos-ai/llmos-gpu-stack/pkg/utils/condition"
)

var (
	// DeviceInitialized indicates whether the GPU device has been initialized
	DeviceInitialized condition.Cond = "Initialized"
)

// +genclient
// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object
// +genclient:nonNamespaced
// +kubebuilder:resource:shortName=gpu,scope=Cluster
// +kubebuilder:subresource:status
// +kubebuilder:printcolumn:name="NODE NAME",type="string",JSONPath=".status.nodeName"
// +kubebuilder:printcolumn:name="VENDOR",type="string",JSONPath=".status.vendor"
// +kubebuilder:printcolumn:name="DEVICE_NAME",type="string",JSONPath=".status.devName"
// +kubebuilder:printcolumn:name="VRAM",type="integer",JSONPath=".status.vram"
// +kubebuilder:printcolumn:name="State",type="string",JSONPath=".status.state"
// +kubebuilder:printcolumn:name="Age",type="date",JSONPath=".metadata.creationTimestamp"

// GPUDevice describes a GPU accelerator device
type GPUDevice struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   GPUDeviceSpec   `json:"spec,omitempty"`
	Status GPUDeviceStatus `json:"status,omitempty"`
}

// GPUDeviceSpec defines the desired state of GPUDevice
type GPUDeviceSpec struct {
}

// GPUDeviceStatus defines the observed state of GPUDevice
type GPUDeviceStatus struct {
	Conditions []common.Condition `json:"conditions,omitempty"`
	// NodeName is the name of the node where the GPU device is located
	NodeName string `json:"nodeName,omitempty"`
	// GPUDeviceInfo is the information of the GPU device
	GPUDeviceInfo `json:",inline"`
	// Pods is the list of pods that are using this GPU device
	Pods []GPUPod `json:"pods,omitempty"`
	// State describes the current state of the GPU device
	State string `json:"state,omitempty"`
}

type GPUDeviceInfo struct {
	// UUID is the GPU Device UUID
	UUID string `json:"uuid,omitempty"`
	// Vendor is the vendor name of the GPU device
	Vendor string `json:"vendor,omitempty"`
	// DevName is the name of the GPU device
	DevName string `json:"devName,omitempty"`
	// MaxCount is the maximum number of splitter instances that can be created from this GPU
	MaxCount int32 `json:"maxCount,omitempty"`
	// VRAM is the amount of video RAM in MB
	VRAM int32 `json:"vram,omitempty"`
	// CUDACores is the number of CUDA cores available on the GPU device
	CUDACores int32 `json:"cudaCores,omitempty"`
	// DevCores is the total percentage number of cores available on the GPU
	DevCores int32 `json:"devCores,omitempty"`
	// Numa is the NUMA node where the GPU device is located
	Numa int `json:"numa,omitempty"`
	// Health indicates whether the GPU device is healthy
	Health bool `json:"health,omitempty"`
}

type GPUPod struct {
	// Name is the namespace:name of the pod, e.g. "default:my-pod"
	Name string `json:"name"`
	// MemReq is the amount of memory requested by the pod
	MemReq int32 `json:"memReq,omitempty"`
	// MemPercentageReq is the percentage of memory requested by the pod
	MemPercentageReq int32 `json:"memPercentageReq,omitempty"`
	// CoresReq is the number of cores requested by the pod
	CoresReq int32 `json:"coresReq,omitempty"`
}
