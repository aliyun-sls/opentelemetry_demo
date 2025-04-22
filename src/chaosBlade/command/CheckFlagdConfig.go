package command

import (
	"chaosBlade/client"
	"context"
	"fmt"
	"github.com/open-feature/go-sdk/openfeature"
	"log"
	"time"
)

// 提取重复的删除CRD逻辑到单独的函数
func deleteCRDIfExists(crdName string) {
	fmt.Printf("删除 %s\n", crdName)
	arr := client.ListCRD(Dynamic, Gvr)
	for _, s := range arr {
		if s == crdName {
			DeleteCRD(Dynamic, Gvr, s)
		}
	}
	fmt.Printf("删除 %s 完成\n", crdName)
}

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
		// 如果配置发生变化，执行相应的操作
		if node.nodeLossLastConfig > 0 {
			deleteCRDIfExists("node-loss")
		}
		node.NodeNetLoss()
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
		// 如果配置发生变化，执行相应的操作
		if region.regionlossLastConfig > 0 {
			deleteCRDIfExists("region-loss")
		}
		region.RegionNetLoss()
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
		// 如果配置发生变化，执行相应的操作
		if podcpu.podCpuLastConfig > 0 {
			deleteCRDIfExists("cpu-load")
			time.Sleep(5 * time.Second)
		}
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
		if podmem.podMemLastConfig > 0 {
			deleteCRDIfExists("mem-load")
			time.Sleep(5 * time.Second)
		}
		// 如果配置发生变化，执行相应的操作
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
	log.Printf("获取 PodNetDelayFlagd feature : %v,last flagd: %v", podnetdelay.flagdValue, podnetdelay.podNetDelayLastConfig)
	if podnetdelay.flagdValue == "on" {
		podnetdelay.labels = "app.kubernetes.io/name=product"
	} else {
		podnetdelay.labels = ""
	}

	if podnetdelay.flagdValue != podnetdelay.podNetDelayLastConfig {
		// 如果配置发生变化，执行相应的操作
		podnetdelay.PodNetDelay()
		podnetdelay.podNetDelayLastConfig = podnetdelay.flagdValue
	}
}
