package command

import (
	"context"
	"github.com/open-feature/go-sdk/openfeature"
	"log"
)

func RDSLossFlagd() {
	var lastConfig, currentConfig string

	// 获取feature flag值
	istrue := FlagClient.String(
		context.Background(),
		"RDSLoss",
		"",
		openfeature.EvaluationContext{},
	)
	log.Printf("获取 RDSLossFlagd feature : %v,last flagd: %v", istrue, lastConfig)
	if istrue == "on" {
		currentConfig = "app.kubernetes.io/name=product"
	} else {
		currentConfig = ""
	}

	if istrue != lastConfig {
		// 如果配置发生变化，执行相应的操作
		RDSLoss(currentConfig)
		lastConfig = istrue
	}
}

func NodeLossFlagd() {
	var lastConfig, currentConfig string

	// 获取feature flag值
	istrue := FlagClient.String(
		context.Background(),
		"NodeLoss",
		"",
		openfeature.EvaluationContext{},
	)
	log.Printf("获取 NodeLossFlagd feature : %v,last flagd: %v", istrue, lastConfig)
	if istrue == "on" {
		currentConfig = "topology.kubernetes.io/zone=cn-guangzhou-a"
	} else {
		currentConfig = ""
	}

	if istrue != lastConfig {
		// 如果配置发生变化，执行相应的操作
		NodeNetLoss(currentConfig)
		lastConfig = istrue
	}
}

func PodCpuFlagd() {
	var lastConfig int64

	// 获取feature flag值
	istrue := FlagClient.Int(
		context.Background(),
		"PodCPULoad",
		0,
		openfeature.EvaluationContext{},
	)
	log.Printf("获取 PodCpuFlagd feature : %v,last flagd: %v", istrue, lastConfig)
	if istrue != lastConfig {
		// 如果配置发生变化，执行相应的操作
		PodCpu(istrue)
		lastConfig = istrue
	}
}

func PodMemFlagd() {
	var lastConfig int64

	// 获取feature flag值
	istrue := FlagClient.Int(
		context.Background(),
		"PodMEMLoad",
		0,
		openfeature.EvaluationContext{},
	)
	log.Printf("获取 PodMemFlagd feature : %v,last flagd: %v", istrue, lastConfig)
	if istrue != lastConfig {
		// 如果配置发生变化，执行相应的操作
		PodMem(istrue)
		lastConfig = istrue
	}
}

func PodNetDelayFlagd() {
	var lastConfig, currentConfig string

	// 获取feature flag值
	istrue := FlagClient.String(
		context.Background(),
		"PodNetDelay",
		"",
		openfeature.EvaluationContext{},
	)
	log.Printf("获取 PodNetDelayFlagd feature : %v,last flagd: %v", istrue, lastConfig)
	if istrue == "on" {
		currentConfig = "app.kubernetes.io/name=product"
	} else {
		currentConfig = ""
	}

	if istrue != lastConfig {
		// 如果配置发生变化，执行相应的操作
		PodNetDelay(currentConfig)
		lastConfig = istrue
	}
}
