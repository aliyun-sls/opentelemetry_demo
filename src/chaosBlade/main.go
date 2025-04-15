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

	for {
		command.RDSLossFlagd()
		command.NodeLossFlagd()
		command.PodNetDelayFlagd()
		command.PodCpuFlagd()
		command.PodMemFlagd()

		time.Sleep(10 * time.Second)
	}
}
