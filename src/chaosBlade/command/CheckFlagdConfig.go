package command

import (
	"context"
	"github.com/open-feature/go-sdk/openfeature"
	"log"
)

func (node *NodeLoss) NodeLossFlagd() {
	// 获取feature flag值
	node.flagdValue = FlagClient.Int(
		context.Background(),
		"NodeLoss",
		0,
		openfeature.EvaluationContext{},
	)
	log.Printf("获取 NodeLossFlagd feature : %v,last flagd: %v", node.flagdValue, node.nodeLossLastConfig)
	if node.flagdValue != node.nodeLossLastConfig {
		node.NodeNetLoss()
		node.nodeLossLastConfig = node.flagdValue
	}
}

// NodeLossAPIFlagd 基于阿里云 API 的节点故障注入 flagd 监控
func (node *NodeLoss) NodeLossAPIFlagd() {
	// 获取feature flag值
	node.flagdValue = FlagClient.Int(
		context.Background(),
		"NodeLoss",
		0,
		openfeature.EvaluationContext{},
	)
	log.Printf("获取 NodeLossAPIFlagd feature : %v,last flagd: %v", node.flagdValue, node.nodeLossLastConfig)

	if node.flagdValue != node.nodeLossLastConfig {
		node.NodeAPILoss()
		node.nodeLossLastConfig = node.flagdValue
	}
}

func (region *RegionLoss) RegionLossFlagd() {
	// 获取feature flag值
	region.flagdValue = FlagClient.Int(
		context.Background(),
		"RegionLoss",
		0,
		openfeature.EvaluationContext{},
	)
	log.Printf("获取 RegionLossFlagd feature : %v,last flagd: %v", region.flagdValue, region.regionlossLastConfig)

	if region.flagdValue != region.regionlossLastConfig {
		region.RegionNetLoss()
		region.regionlossLastConfig = region.flagdValue
	}
}

// RegionLossAPIFlagd 基于阿里云 API 的区域故障注入 flagd 监控
func (region *RegionLoss) RegionLossAPIFlagd() {
	// 获取feature flag值
	region.flagdValue = FlagClient.Int(
		context.Background(),
		"RegionLoss",
		0,
		openfeature.EvaluationContext{},
	)
	log.Printf("获取 RegionLossAPIFlagd feature : %v,last flagd: %v", region.flagdValue, region.regionlossLastConfig)

	if region.flagdValue != region.regionlossLastConfig {
		region.RegionAPILoss()
		region.regionlossLastConfig = region.flagdValue
	}
}

func (podcpu *PodCpu) PodCpuFlagd() {
	// 获取feature flag值
	podcpu.flagdValue = FlagClient.Int(
		context.Background(),
		"PodCPULoad",
		0,
		openfeature.EvaluationContext{},
	)
	log.Printf("获取 PodCpuFlagd feature : %v,last flagd: %v", podcpu.flagdValue, podcpu.podCpuLastConfig)
	if podcpu.flagdValue != podcpu.podCpuLastConfig {
		podcpu.PodCpu()
		podcpu.podCpuLastConfig = podcpu.flagdValue
	}
}

func (podmem *PodMem) PodMemFlagd() {
	// 获取feature flag值
	podmem.flagdValue = FlagClient.Int(
		context.Background(),
		"PodMEMLoad",
		0,
		openfeature.EvaluationContext{},
	)
	log.Printf("获取 PodMemFlagd feature : %v,last flagd: %v", podmem.flagdValue, podmem.podMemLastConfig)
	if podmem.flagdValue != podmem.podMemLastConfig {
		podmem.PodMem()
		podmem.podMemLastConfig = podmem.flagdValue
	}
}

func (podnetdelay *PodNetDelay) PodNetDelayFlagd() {
	// 获取feature flag值
	podnetdelay.flagdValue = FlagClient.String(
		context.Background(),
		"PodNetDelay",
		"",
		openfeature.EvaluationContext{},
	)
	log.Printf("获取 PodNetDelayFlagd feature : \"%v\",last flagd: \"%v\"", podnetdelay.flagdValue, podnetdelay.podNetDelayLastConfig)
	if podnetdelay.flagdValue == "on" {
		podnetdelay.labels = "app.kubernetes.io/name=product"
	} else {
		podnetdelay.labels = ""
	}

	if podnetdelay.flagdValue != podnetdelay.podNetDelayLastConfig {
		podnetdelay.PodNetDelay()
		podnetdelay.podNetDelayLastConfig = podnetdelay.flagdValue
	}
}

func (nodecpu *NodeCpu) NodeCpuFlagd() {
	// 获取feature flag值
	nodecpu.flagdValue = FlagClient.Int(
		context.Background(),
		"NodeCPULoad",
		0,
		openfeature.EvaluationContext{},
	)
	log.Printf("获取 NodeCPULoad feature : %v,last flagd: %v", nodecpu.flagdValue, nodecpu.nodeCpuLastConfig)
	if nodecpu.flagdValue != nodecpu.nodeCpuLastConfig {
		nodecpu.NodeCpu()
		nodecpu.nodeCpuLastConfig = nodecpu.flagdValue
	}
}

func (nodemem *NodeMem) NodeMemFlagd() {
	// 获取feature flag值
	nodemem.flagdValue = FlagClient.Int(
		context.Background(),
		"NodeMemLoad",
		0,
		openfeature.EvaluationContext{},
	)
	log.Printf("获取 NodeMemLoad feature : %v,last flagd: %v", nodemem.flagdValue, nodemem.nodeMemLastConfig)
	if nodemem.flagdValue != nodemem.nodeMemLastConfig {
		nodemem.NodeMem()
		nodemem.nodeMemLastConfig = nodemem.flagdValue
	}
}

func (nodedisk *NodeDisk) NodeDiskFlagd() {
	// 获取feature flag值
	nodedisk.flagdValue = FlagClient.Int(
		context.Background(),
		"NodeDiskLoad",
		0,
		openfeature.EvaluationContext{},
	)
	log.Printf("获取 NodeDiskLoad feature : %v,last flagd: %v", nodedisk.flagdValue, nodedisk.nodeDiskLastConfig)
	if nodedisk.flagdValue != nodedisk.nodeDiskLastConfig {
		nodedisk.NodeDisk()
		nodedisk.nodeDiskLastConfig = nodedisk.flagdValue
	}
}
