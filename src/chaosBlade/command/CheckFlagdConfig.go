package command

import (
	"context"
	"github.com/open-feature/go-sdk/openfeature"
	"log"
	"time"
)

// 将 lastConfig 提升为包级别的全局变量
var (
	rdsLossLastConfig     string
	nodeLossLastConfig    string
	podCpuLastConfig      int64
	podMemLastConfig      int64
	podNetDelayLastConfig string
)

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
	var currentConfig string

	// 获取feature flag值
	istrue := FlagClient.String(
		context.Background(),
		"NodeLoss",
		"",
		openfeature.EvaluationContext{},
	)
	log.Printf("获取 NodeLossFlagd feature : %v,last flagd: %v", istrue, nodeLossLastConfig)
	if istrue == "on" {
		currentConfig = "topology.kubernetes.io/zone=cn-guangzhou-a"
	} else {
		currentConfig = ""
	}

	if istrue != nodeLossLastConfig {
		// 如果配置发生变化，执行相应的操作
		NodeNetLoss(currentConfig)
		nodeLossLastConfig = istrue
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
		if podCpuLastConfig > 0 {
			DeleteCRD(Dynamic, Gvr, "cpu-load")
			time.Sleep(5 * time.Second)
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
		if podCpuLastConfig > 0 {
			DeleteCRD(Dynamic, Gvr, "mem-load")
			time.Sleep(5 * time.Second)
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
