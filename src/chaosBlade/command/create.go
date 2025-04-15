package command

import (
	"fmt"
)

// --------------region网络丢包------------
func NodeNetLoss() {
	// 生成 yaml 文件
	nodeNetLoss := NodeNetLossYaml()
	// 根据 yaml 文件创建资源
	data := createCRD(Dynamic, Gvr, nodeNetLoss)
	fmt.Println("创建资源:", data)
}

// ------------RDS断连---------------
func RDSloss(labels string) {
	// 生成 yaml 文件
	RDSLoss := RDSlossYaml(labels)
	// 根据 yaml 文件创建资源
	data := createCRD(Dynamic, Gvr, RDSLoss)
	fmt.Println("创建资源:", data)
}
