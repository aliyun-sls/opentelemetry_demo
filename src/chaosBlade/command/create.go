package command

import (
	"fmt"
)

// --------------region网络丢包------------
func NodeNetLoss(istrue, labels string) {
	// 生成 yaml 文件
	nodeNetLoss := NodeNetLossYaml(istrue, labels)
	if nodeNetLoss == "" {
		return
	}
	// 根据 yaml 文件创建资源
	data := createCRD(Dynamic, Gvr, nodeNetLoss)
	fmt.Println("创建NodeNetLoss资源:", data)
}

// ------------RDS断连---------------
func RDSLoss(labels string) {
	// 生成 yaml 文件
	RDSloss := RDSlossYaml(labels)
	if RDSloss == "" {
		return
	}
	// 根据 yaml 文件创建资源
	data := createCRD(Dynamic, Gvr, RDSloss)
	fmt.Println("创建RDSLoss资源:", data)
}

func PodNetDelay(labels string) {
	// 生成 yaml 文件
	podNetDelay := PodNetDelayYaml(labels)
	if podNetDelay == "" {
		return
	}
	// 根据 yaml 文件创建资源
	data := createCRD(Dynamic, Gvr, podNetDelay)
	fmt.Println("创建PodNetDelay资源:", data)
}

func PodCpu(percent int64) {
	podCpu := PodCpuYaml(percent)
	if podCpu == "" {
		return
	}
	data := createCRD(Dynamic, Gvr, podCpu)
	fmt.Println("创建PodCpu资源:", data)
}

func PodMem(percent int64) {
	podMem := PodMemYaml(percent)
	if podMem == "" {
		return
	}
	data := createCRD(Dynamic, Gvr, podMem)
	fmt.Println("创建PodMem资源:", data)
}
