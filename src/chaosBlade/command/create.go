package command

import (
	"fmt"
)

// --------------region网络丢包------------
func (node *NodeLoss) NodeNetLoss() {
	// 生成 yaml 文件
	nodeNetLoss := node.NodeNetLossYaml()
	if nodeNetLoss == "" {
		return
	}
	// 根据 yaml 文件创建资源
	data := createCRD(Dynamic, Gvr, nodeNetLoss)
	fmt.Println("创建NodeNetLoss资源:", data)
}

func (region *RegionLoss) RegionNetLoss() {
	// 生成 yaml 文件
	nodeNetLoss := region.RegionNetLossYaml()
	if nodeNetLoss == "" {
		return
	}
	// 根据 yaml 文件创建资源
	data := createCRD(Dynamic, Gvr, nodeNetLoss)
	fmt.Println("创建RegionNetLoss资源:", data)
}

// RegionAPILoss 基于阿里云 API 的区域故障注入
func (region *RegionLoss) RegionAPILoss() {
	AliyunRegionChaosOnce.Do(func() {
		AliyunRegionChaosInstance = NewAliyunRegionChaos()
	})

	AliyunRegionChaosInstance.ExecuteChaos(region.flagdValue)
}

func (podnetdelay *PodNetDelay) PodNetDelay() {
	// 生成 yaml 文件
	podNetDelay := podnetdelay.PodNetDelayYaml()
	if podNetDelay == "" {
		return
	}
	// 根据 yaml 文件创建资源
	data := createCRD(Dynamic, Gvr, podNetDelay)
	fmt.Println("创建PodNetDelay资源:", data)
}

func (podcpu *PodCpu) PodCpu() {
	podCpu := podcpu.PodCpuYaml()
	if podCpu == "" {
		return
	}
	data := createCRD(Dynamic, Gvr, podCpu)
	fmt.Println("创建PodCpu资源:", data)
}

func (podmem *PodMem) PodMem() {
	podMem := podmem.PodMemYaml()
	if podMem == "" {
		return
	}
	data := createCRD(Dynamic, Gvr, podMem)
	fmt.Println("创建PodMem资源:", data)
}

func (nodecpu *NodeCpu) NodeCpu() {
	nodeCpu := nodecpu.NodeCpuYaml()
	if nodeCpu == "" {
		return
	}
	data := createCRD(Dynamic, Gvr, nodeCpu)
	fmt.Println("创建NodeCpu资源:", data)
}

func (nodemem *NodeMem) NodeMem() {
	nodeMem := nodemem.NodeMemYaml()
	if nodeMem == "" {
		return
	}
	data := createCRD(Dynamic, Gvr, nodeMem)
	fmt.Println("创建NodeCpu资源:", data)
}

func (nodedisk *NodeDisk) NodeDisk() {
	nodeDisk := nodedisk.NodeDiskYaml()
	if nodeDisk == "" {
		return
	}
	data := createCRD(Dynamic, Gvr, nodeDisk)
	fmt.Println("创建NodeDisk资源:", data)
}
