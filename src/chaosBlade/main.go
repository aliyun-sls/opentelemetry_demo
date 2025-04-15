package main

import (
	"chaosBlade/client"
	"chaosBlade/command"
	"fmt"
	"time"
)

func main() {
	command.Dynamic, command.Gvr = client.Client()
	fmt.Println("初始化客户端")
	time.Sleep(2 * time.Second)

	//清理环境
	arr := client.ListCRD(command.Dynamic, command.Gvr)
	if arr != nil {
		for _, s := range arr {
			command.DeleteCRD(command.Dynamic, command.Gvr, s)
		}
	}
	fmt.Println("清理环境")

	// 持续运行并等待flagd配置触发操作
	var lastConfig string
	for {
		// 检查flagd配置是否发生变化
		currentConfig := command.CheckRDSFlagdConfig()
		if currentConfig != "" && currentConfig != lastConfig {
			// 如果配置发生变化，执行相应的操作
			command.RDSloss(currentConfig)
			lastConfig = currentConfig
		} else if currentConfig == "" && currentConfig != lastConfig {
			command.RDSloss(currentConfig)
			lastConfig = currentConfig
		}
	}

	// 等待一段时间后再次检查
	time.Sleep(10 * time.Second)
}

// checkFlagdConfig 检查flagd配置是否发生变化
