package utils

const (
	// OneContainerMultiDeviceSplitSymbol this is when one container use multi device, use : symbol to join device info.
	OneContainerMultiDeviceSplitSymbol = ":"

	// OnePodMultiContainerSplitSymbol this is when one pod having multi container and
	// more than one container use device, use ; symbol to join device info.
	OnePodMultiContainerSplitSymbol = ";"
)

type DeviceInfo struct {
	Index                int        `json:"index,omitempty"`
	Id                   string     `protobuf:"bytes,1,opt,name=id,proto3" json:"id,omitempty"`
	Count                int32      `protobuf:"varint,2,opt,name=count,proto3" json:"count,omitempty"`
	Devmem               int32      `protobuf:"varint,3,opt,name=devmem,proto3" json:"devmem,omitempty"`
	Type                 string     `protobuf:"bytes,4,opt,name=type,proto3" json:"type,omitempty"`
	Health               bool       `protobuf:"varint,5,opt,name=health,proto3" json:"health,omitempty"`
	Mode                 string     `json:"mode,omitempty"`
	MIGTemplate          []Geometry `json:"migtemplate,omitempty"`
	XXX_NoUnkeyedLiteral struct{}   `json:"-"`
	XXX_unrecognized     []byte     `json:"-"`
	XXX_sizecache        int32      `json:"-"`
}

type ContainerDevices []ContainerDevice

type PodSingleDevice []ContainerDevices

type PodDevices map[string]PodSingleDevice

type ContainerDevice struct {
	UUID      string
	Type      string
	Usedmem   int32
	Usedcores int32
}

type Geometry struct {
	Group     string        `yaml:"group"`
	Instances []MigTemplate `yaml:"geometries"`
}

type MigTemplate struct {
	Name   string `yaml:"name"`
	Memory int32  `yaml:"memory"`
	Count  int32  `yaml:"count"`
}
