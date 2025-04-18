package command

import (
	"context"
	"fmt"
	"github.com/open-feature/go-sdk/openfeature"
	"log"
	"time"
)

// 将 lastConfig 提升为包级别的全局变量
var (
	rdsLossLastConfig     string
	nodeLossLastConfig    int64
	podCpuLastConfig      int64
	podMemLastConfig      int64
	podNetDelayLastConfig string
)

type RegionLoss struct {
	regionlossLastConfig int64
	flagdValue           int64
}

func RDSLossFlagd() {
	var currentConfig string

	// 获取feature flag值
	istrue := FlagClient.String(
		context.Background(),
		"RDSLoss",
		"",
		openfeature.EvaluationContext{},
	)
	log.Printf("获取 RDSLossFlagd feature : %v,last flagd: %v", istrue, rdsLossLastConfig)
	if istrue == "on" {
		currentConfig = "app.kubernetes.io/name=product"
	} else {
		currentConfig = ""
	}

	if istrue != rdsLossLastConfig {
		// 如果配置发生变化，执行相应的操作
		RDSLoss(currentConfig)
		rdsLossLastConfig = istrue
	}
}

func NodeLossFlagd() {
	// 获取feature flag值
	istrue := FlagClient.Int(
		context.Background(),
		"NodeLoss",
		0,
		openfeature.EvaluationContext{},
	)
	log.Printf("获取 NodeLossFlagd feature : %v,last flagd: %v", istrue, nodeLossLastConfig)
	if istrue != nodeLossLastConfig {
		// 如果配置发生变化，执行相应的操作
		fmt.Println("删除node-loss")
		if nodeLossLastConfig > 0 {
			DeleteCRD(Dynamic, Gvr, "node-loss")
			time.Sleep(5 * time.Second)
			fmt.Println("删除node-loss完成")
		}
		NodeNetLoss(istrue)
		nodeLossLastConfig = istrue
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
	log.Printf("获取 NodeLossFlagd feature : %v,last flagd: %v", region.flagdValue, region.regionlossLastConfig)

	if region.flagdValue != region.regionlossLastConfig {
		// 如果配置发生变化，执行相应的操作
		fmt.Println("删除region-loss")
		if region.regionlossLastConfig > 0 {
			DeleteCRD(Dynamic, Gvr, "region-loss")
			time.Sleep(5 * time.Second)
			fmt.Println("删除region-loss完成")
		}
		region.RegionNetLoss()
		region.regionlossLastConfig = region.flagdValue
	}
}

func PodCpuFlagd() {
	// 获取feature flag值
	istrue := FlagClient.Int(
		context.Background(),
		"PodCPULoad",
		0,
		openfeature.EvaluationContext{},
	)
	log.Printf("获取 PodCpuFlagd feature : %v,last flagd: %v", istrue, podCpuLastConfig)
	if istrue != podCpuLastConfig {
		// 如果配置发生变化，执行相应的操作
		fmt.Println("删除cpu-load")
		if podCpuLastConfig > 0 {
			DeleteCRD(Dynamic, Gvr, "cpu-load")
			time.Sleep(5 * time.Second)
			fmt.Println("删除cpu-load完成")
		}
		PodCpu(istrue)
		podCpuLastConfig = istrue
	}
}

func PodMemFlagd() {
	// 获取feature flag值
	istrue := FlagClient.Int(
		context.Background(),
		"PodMEMLoad",
		0,
		openfeature.EvaluationContext{},
	)
	log.Printf("获取 PodMemFlagd feature : %v,last flagd: %v", istrue, podMemLastConfig)
	if istrue != podMemLastConfig {
		if podMemLastConfig > 0 {
			fmt.Println("删除mem-load")
			DeleteCRD(Dynamic, Gvr, "mem-load")
			time.Sleep(5 * time.Second)
			fmt.Println("删除mem-load完成")
		}
		// 如果配置发生变化，执行相应的操作
		PodMem(istrue)
		podMemLastConfig = istrue
	}
}

func PodNetDelayFlagd() {
	var currentConfig string

	// 获取feature flag值
	istrue := FlagClient.String(
		context.Background(),
		"PodNetDelay",
		"",
		openfeature.EvaluationContext{},
	)
	log.Printf("获取 PodNetDelayFlagd feature : %v,last flagd: %v", istrue, podNetDelayLastConfig)
	if istrue == "on" {
		currentConfig = "app.kubernetes.io/name=product"
	} else {
		currentConfig = ""
	}

	if istrue != podNetDelayLastConfig {
		// 如果配置发生变化，执行相应的操作
		PodNetDelay(currentConfig)
		podNetDelayLastConfig = istrue
	}
}
