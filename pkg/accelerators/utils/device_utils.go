package utils

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/sirupsen/logrus"
)

func DecodeNodeDevices(str string) ([]*DeviceInfo, error) {
	if !strings.Contains(str, ":") {
		return []*DeviceInfo{}, fmt.Errorf("node annotations not decode successfully")
	}
	tmp := strings.Split(str, ":")
	var retval []*DeviceInfo
	for i, val := range tmp {
		if strings.Contains(val, ",") {
			items := strings.Split(val, ",")
			if len(items) == 5 {
				count, _ := strconv.Atoi(items[1])
				devmem, _ := strconv.Atoi(items[2])
				//health, _ := strconv.ParseBool(items[4])
				i := DeviceInfo{
					Index:  i,
					Id:     items[0],
					Count:  int32(count),
					Devmem: int32(devmem),
					Type:   items[3],
					Health: true, //TODO, fix the health status, is always false for now
				}
				retval = append(retval, &i)
			} else {
				return []*DeviceInfo{}, fmt.Errorf("node annotations not decode successfully")
			}
		}
	}
	return retval, nil
}

func EncodeNodeDevices(dlist []*DeviceInfo) string {
	tmp := ""
	for _, val := range dlist {
		tmp += val.Id + "," + strconv.FormatInt(int64(val.Count), 10) + "," +
			strconv.Itoa(int(val.Devmem)) + "," + val.Type + "," + strconv.FormatBool(val.Health) + ":"
	}
	logrus.Debugf("encoded node devices %s", tmp)
	return tmp
}

func DecodePodDevices(checkList map[string]string, annos map[string]string) (PodDevices, error) {
	logrus.Debugf("validate pod annotations [%+v] with accelerator checklist [%+v]", annos, checkList)
	if len(annos) == 0 {
		return PodDevices{}, nil
	}
	pd := make(PodDevices)

	for devID, ac := range checkList {
		str, ok := annos[ac]
		if !ok {
			continue
		}
		pd[devID] = make(PodSingleDevice, 0)
		for _, s := range strings.Split(str, OnePodMultiContainerSplitSymbol) {
			cd, err := DecodeContainerDevices(s)
			if err != nil {
				return PodDevices{}, nil
			}
			if len(cd) == 0 {
				continue
			}
			pd[devID] = append(pd[devID], cd)
		}
	}
	return pd, nil
}

func DecodeContainerDevices(str string) (ContainerDevices, error) {
	if len(str) == 0 {
		return ContainerDevices{}, nil
	}
	cd := strings.Split(str, OneContainerMultiDeviceSplitSymbol)
	contdev := ContainerDevices{}
	tmpdev := ContainerDevice{}
	if len(str) == 0 {
		return contdev, nil
	}
	for _, val := range cd {
		if strings.Contains(val, ",") {
			tmpstr := strings.Split(val, ",")
			if len(tmpstr) < 4 {
				return ContainerDevices{}, fmt.Errorf("pod annotation format error; required split count is 4, val:[%s]", val)
			}
			tmpdev.UUID = tmpstr[0]
			tmpdev.Type = tmpstr[1]
			devmem, _ := strconv.ParseInt(tmpstr[2], 10, 32)
			tmpdev.Usedmem = int32(devmem)
			devcores, _ := strconv.ParseInt(tmpstr[3], 10, 32)
			tmpdev.Usedcores = int32(devcores)
			contdev = append(contdev, tmpdev)
		}
	}
	return contdev, nil
}
