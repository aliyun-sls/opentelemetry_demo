package command

import (
	"fmt"
)

// --------------region网络丢包------------
func NodeNetLoss(labels string) {
	// 生成 yaml 文件
	nodeNetLoss := NodeNetLossYaml(labels)
	// 根据 yaml 文件创建资源
	data := createCRD(Dynamic, Gvr, nodeNetLoss)
	fmt.Println("创建资源:", data)
}

// ------------RDS断连---------------
func RDSLoss(labels string) {
	// 生成 yaml 文件
	RDSloss := RDSlossYaml(labels)
	// 根据 yaml 文件创建资源
	data := createCRD(Dynamic, Gvr, RDSloss)
	fmt.Println("创建资源:", data)
}

func PodNetDelay(labels string) {
	// 生成 yaml 文件
	podNetDelay := PodNetDelayYaml(labels)
	// 根据 yaml 文件创建资源
	data := createCRD(Dynamic, Gvr, podNetDelay)
	fmt.Println("创建资源:", data)
}

func PodCpu(percent int64) {
	podCpu := PodCpuYaml(percent)
	data := createCRD(Dynamic, Gvr, podCpu)
	fmt.Println("创建资源:", data)
}

func PodMem(percent int64) {
	podCpu := PodMemYaml(percent)
	data := createCRD(Dynamic, Gvr, podCpu)
	fmt.Println("创建资源:", data)
}
