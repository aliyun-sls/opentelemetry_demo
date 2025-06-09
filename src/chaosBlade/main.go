package main

import (
	"chaosBlade/client"
	"chaosBlade/command"
	"fmt"
	flagd "github.com/open-feature/go-sdk-contrib/providers/flagd/pkg"
	"github.com/open-feature/go-sdk/openfeature"
	"log"
	"time"
)

func main() {
	command.Dynamic, command.Gvr = client.Client()
	fmt.Println("初始化 Kubernetes 客户端")

	// 初始化阿里云客户端
	command.InitAliyunClients()
	fmt.Println("初始化阿里云客户端")

	time.Sleep(2 * time.Second)

	//清理环境
	arr := client.ListCRD(command.Dynamic, command.Gvr)
	if arr != nil {
		for _, s := range arr {
			command.DeleteCRD(command.Dynamic, command.Gvr, s)
		}
	}
	fmt.Println("清理环境")

	command.Once.Do(func() {
		provider := flagd.NewProvider(
			flagd.WithHost("flagd"),
			flagd.WithPort(8013),
		)
		if err := openfeature.SetProvider(provider); err != nil {
			log.Printf("设置provider失败: %v", err)
			return
		}
		command.FlagClient = openfeature.NewClient("chaosblade-run")
	})

	region := command.RegionLoss{}
	nodeloss := command.NodeLoss{}
	podnetdelay := command.PodNetDelay{}
	podcpu := command.PodCpu{}
	podmem := command.PodMem{}
	nodecpu := command.NodeCpu{}
	nodemem := command.NodeMem{}
	nodedisk := command.NodeDisk{}

	for {
		nodeloss.NodeLossAPIFlagd()
		region.RegionLossAPIFlagd()
		podnetdelay.PodNetDelayFlagd()
		podcpu.PodCpuFlagd()
		podmem.PodMemFlagd()
		nodecpu.NodeCpuFlagd()
		nodemem.NodeMemFlagd()
		nodedisk.NodeDiskFlagd()

		time.Sleep(30 * time.Second)
	}
}
