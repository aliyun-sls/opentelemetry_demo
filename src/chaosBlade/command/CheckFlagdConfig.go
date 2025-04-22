package command

import (
	"chaosBlade/client"
	"context"
	"fmt"
	"github.com/open-feature/go-sdk/openfeature"
	"log"
	"time"
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
		// 如果配置发生变化，执行相应的操作
		if node.nodeLossLastConfig > 0 {
			fmt.Println("删除node-loss")
			arr := client.ListCRD(Dynamic, Gvr)
			for _, s := range arr {
				if s == "node-loss" {
					DeleteCRD(Dynamic, Gvr, s)
				}
			}
			fmt.Println("删除node-loss完成")
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
			fmt.Println("删除region-loss")
			arr := client.ListCRD(Dynamic, Gvr)
			for _, s := range arr {
				if s == "region-loss" {
					DeleteCRD(Dynamic, Gvr, s)
				}
			}
			fmt.Println("删除region-loss完成")
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
			fmt.Println("删除cpu-load")
			DeleteCRD(Dynamic, Gvr, "cpu-load")
			time.Sleep(5 * time.Second)
			fmt.Println("删除cpu-load完成")
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
			fmt.Println("删除mem-load")
			DeleteCRD(Dynamic, Gvr, "mem-load")
			time.Sleep(5 * time.Second)
			fmt.Println("删除mem-load完成")
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
