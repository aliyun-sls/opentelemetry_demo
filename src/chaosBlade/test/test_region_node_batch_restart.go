package main

import (
	"fmt"
	openapi "github.com/alibabacloud-go/darabonba-openapi/v2/client"
	ecs "github.com/alibabacloud-go/ecs-20140526/v4/client"
	"github.com/alibabacloud-go/tea/tea"
	"log"
	"os"
	"strings"
	"time"
)

func main() {
	// 从环境变量获取配置
	accessKeyId := os.Getenv("ALIBABA_CLOUD_ACCESS_KEY_ID")
	accessKeySecret := os.Getenv("ALIBABA_CLOUD_ACCESS_KEY_SECRET")
	regionId := os.Getenv("ALIBABA_CLOUD_REGION_ID")
	clusterId := os.Getenv("ALIBABA_CLOUD_CLUSTER_ID")
	zoneId := os.Getenv("ALIBABA_CLOUD_ZONE_ID")

	// 检查必需的环境变量
	if accessKeyId == "" || accessKeySecret == "" || regionId == "" {
		log.Fatal("❌ 请设置必需的环境变量: ALIBABA_CLOUD_ACCESS_KEY_ID, ALIBABA_CLOUD_ACCESS_KEY_SECRET, ALIBABA_CLOUD_REGION_ID")
	}

	// 要重启的实例列表 (从环境变量获取，用逗号分隔)
	instancesEnv := os.Getenv("RESTART_INSTANCES")
	var instancesToRestart []string

	if instancesEnv != "" {
		// 如果环境变量有设置，解析实例列表
		instancesToRestart = strings.Split(strings.TrimSpace(instancesEnv), ",")
		for i, id := range instancesToRestart {
			instancesToRestart[i] = strings.TrimSpace(id)
		}
	} else {
		// 如果没有设置环境变量，给出提示
		fmt.Println("⚠️  未设置 RESTART_INSTANCES 环境变量")
		fmt.Println("💡 使用方法:")
		fmt.Println("   export RESTART_INSTANCES=\"i-bp1xxx,i-bp2xxx\"")
		fmt.Println("   或单个实例: export RESTART_INSTANCES=\"i-bp1xxx\"")
		fmt.Println("")
		fmt.Println("🔍 正在检查是否有已停止的实例...")

		// 尝试通过可用区查找已停止的实例
		if zoneId != "" {
			foundInstances, err := findStoppedInstances(accessKeyId, accessKeySecret, regionId, zoneId)
			if err != nil {
				log.Fatalf("查找已停止实例失败: %v", err)
			}
			if len(foundInstances) > 0 {
				fmt.Printf("📋 在可用区 %s 发现 %d 个已停止的实例:\n", zoneId, len(foundInstances))
				for i, id := range foundInstances {
					fmt.Printf("   %d. %s\n", i+1, id)
				}
				instancesToRestart = foundInstances
			} else {
				fmt.Printf("❌ 在可用区 %s 未发现已停止的实例\n", zoneId)
				return
			}
		} else {
			fmt.Println("❌ 请设置 RESTART_INSTANCES 或 ZONE_ID 环境变量")
			return
		}
	}

	fmt.Printf("🔄 ECS 实例批量启动/重启程序\n")
	fmt.Printf("Region: %s\n", regionId)
	if zoneId != "" {
		fmt.Printf("Zone ID: %s\n", zoneId)
	}
	fmt.Printf("要处理的实例数量: %d\n", len(instancesToRestart))
	fmt.Printf("实例列表: %v\n", instancesToRestart)

	// 创建 ECS 客户端
	fmt.Println("\n=== 初始化 ECS 客户端 ===")
	ecsConfig := &openapi.Config{
		AccessKeyId:     tea.String(accessKeyId),
		AccessKeySecret: tea.String(accessKeySecret),
		RegionId:        tea.String(regionId),
		Endpoint:        tea.String("ecs." + regionId + ".aliyuncs.com"),
	}

	ecsClient, err := ecs.NewClient(ecsConfig)
	if err != nil {
		log.Fatalf("创建 ECS 客户端失败: %v", err)
	}
	fmt.Println("✅ ECS 客户端创建成功")

	// 检查实例状态
	fmt.Println("\n=== 检查实例状态 ===")
	statusRequest := &ecs.DescribeInstanceStatusRequest{
		InstanceId: tea.StringSlice(instancesToRestart),
		RegionId:   tea.String(regionId),
	}

	statusResponse, err := ecsClient.DescribeInstanceStatus(statusRequest)
	if err != nil {
		log.Fatalf("获取实例状态失败: %v", err)
	}

	// 分析实例状态
	var stoppedInstances []string     // 需要启动的实例（Stopped/Shutted）
	var runningInstances []string     // 需要重启的实例（Running）
	var otherStatusInstances []string // 其他状态的实例

	if statusResponse.Body.InstanceStatuses != nil && statusResponse.Body.InstanceStatuses.InstanceStatus != nil {
		for _, status := range statusResponse.Body.InstanceStatuses.InstanceStatus {
			if status.InstanceId != nil && status.Status != nil {
				instanceId := *status.InstanceId
				instanceStatus := *status.Status

				fmt.Printf("📊 实例 %s 状态: %s\n", instanceId, instanceStatus)

				switch instanceStatus {
				case "Stopped", "Shutted": // 已停止状态，需要启动
					stoppedInstances = append(stoppedInstances, instanceId)
				case "Running": // 运行中，需要重启
					runningInstances = append(runningInstances, instanceId)
				default:
					otherStatusInstances = append(otherStatusInstances, instanceId)
				}
			}
		}
	}

	fmt.Printf("\n📈 状态统计:\n")
	fmt.Printf("   - 已停止(需启动): %d 个\n", len(stoppedInstances))
	fmt.Printf("   - 正在运行(需重启): %d 个\n", len(runningInstances))
	fmt.Printf("   - 其他状态: %d 个\n", len(otherStatusInstances))

	if len(otherStatusInstances) > 0 {
		fmt.Printf("⚠️  以下实例状态异常，请检查: %v\n", otherStatusInstances)
	}

	// 启动已停止的实例
	if len(stoppedInstances) > 0 {
		fmt.Printf("\n=== 🚀 开始启动已停止的实例 ===\n")
		fmt.Printf("🎯 准备启动 %d 个已停止的实例: %v\n", len(stoppedInstances), stoppedInstances)
		fmt.Printf("📅 启动时间: %s\n", time.Now().Format("2006-01-02 15:04:05"))

		// 执行启动操作 (使用 StartInstances 而不是 RebootInstances)
		startRequest := &ecs.StartInstancesRequest{
			InstanceId: tea.StringSlice(stoppedInstances),
			RegionId:   tea.String(regionId),
		}

		startResponse, err := ecsClient.StartInstances(startRequest)
		if err != nil {
			log.Printf("❌ 启动失败: %v", err)
		} else {
			fmt.Printf("✅ 启动请求发送成功!\n")
			fmt.Printf("📋 启动详细信息:\n")
			fmt.Printf("   - 实例数量: %d\n", len(stoppedInstances))
			fmt.Printf("   - 实例列表: %v\n", stoppedInstances)
			fmt.Printf("   - 地域: %s\n", regionId)
			fmt.Printf("   - 操作类型: StartInstances (启动已停止的实例)\n")

			if startResponse.Body.RequestId != nil {
				fmt.Printf("   - 请求ID: %s\n", *startResponse.Body.RequestId)
			}

			fmt.Printf("⏱️  实例启动中，通常需要60-120秒完成\n")

			// 等待并检查启动进度
			fmt.Println("\n⏳ 等待10秒后开始检查启动进度...")
			time.Sleep(10 * time.Second)

			for i := 0; i < 6; i++ { // 检查6次，共约60秒
				fmt.Printf("\n🔍 第 %d 次状态检查 (%s):\n", i+1, time.Now().Format("15:04:05"))

				checkRequest := &ecs.DescribeInstanceStatusRequest{
					InstanceId: tea.StringSlice(stoppedInstances),
					RegionId:   tea.String(regionId),
				}

				checkResponse, err := ecsClient.DescribeInstanceStatus(checkRequest)
				if err != nil {
					log.Printf("❌ 检查状态失败: %v", err)
					continue
				}

				var runningCount int
				var startingCount int
				var stoppedCount int

				if checkResponse.Body.InstanceStatuses != nil && checkResponse.Body.InstanceStatuses.InstanceStatus != nil {
					for _, status := range checkResponse.Body.InstanceStatuses.InstanceStatus {
						if status.InstanceId != nil && status.Status != nil {
							instanceId := *status.InstanceId
							instanceStatus := *status.Status

							fmt.Printf("   📊 %s: %s\n", instanceId, instanceStatus)

							switch instanceStatus {
							case "Running":
								runningCount++
							case "Starting":
								startingCount++
							case "Stopped", "Shutted":
								stoppedCount++
							}
						}
					}
				}

				fmt.Printf("   ✅ 已启动: %d, 🔄 启动中: %d, ❌ 仍停止: %d\n", runningCount, startingCount, stoppedCount)

				if runningCount == len(stoppedInstances) {
					fmt.Printf("\n🎉 所有实例已成功启动!\n")
					break
				}

				if i < 5 { // 不是最后一次检查
					fmt.Printf("⏳ 等待10秒后继续检查...\n")
					time.Sleep(10 * time.Second)
				}
			}
		}
	}

	// 重启正在运行的实例
	if len(runningInstances) > 0 {
		fmt.Printf("\n=== 🔄 开始重启正在运行的实例 ===\n")
		fmt.Printf("🎯 准备重启 %d 个正在运行的实例: %v\n", len(runningInstances), runningInstances)
		fmt.Printf("📅 重启时间: %s\n", time.Now().Format("2006-01-02 15:04:05"))

		// 执行重启操作 (使用 RebootInstances)
		rebootRequest := &ecs.RebootInstancesRequest{
			InstanceId: tea.StringSlice(runningInstances),
			RegionId:   tea.String(regionId),
		}

		rebootResponse, err := ecsClient.RebootInstances(rebootRequest)
		if err != nil {
			log.Printf("❌ 重启失败: %v", err)
		} else {
			fmt.Printf("✅ 重启请求发送成功!\n")
			fmt.Printf("📋 重启详细信息:\n")
			fmt.Printf("   - 实例数量: %d\n", len(runningInstances))
			fmt.Printf("   - 实例列表: %v\n", runningInstances)
			fmt.Printf("   - 地域: %s\n", regionId)
			fmt.Printf("   - 操作类型: RebootInstances (重启正在运行的实例)\n")

			if rebootResponse.Body.RequestId != nil {
				fmt.Printf("   - 请求ID: %s\n", *rebootResponse.Body.RequestId)
			}

			fmt.Printf("⏱️  实例重启中，通常需要60-120秒完成\n")
		}
	}

	if len(stoppedInstances) == 0 && len(runningInstances) == 0 {
		fmt.Println("\n⚠️  没有需要处理的实例")
		fmt.Println("💡 可能的原因:")
		fmt.Println("   - 实例ID不存在")
		fmt.Println("   - 实例状态异常")
		fmt.Println("   - 网络连接问题")
	}

	fmt.Println("\n🎉 批量启动/重启操作完成!")
	fmt.Println("📝 操作后建议:")
	fmt.Println("   - 检查应用服务是否正常运行")
	fmt.Println("   - 验证集群节点是否重新加入")
	fmt.Println("   - 确认相关业务功能正常")

	fmt.Println("\n=== 程序执行完成 ===")
}

// findStoppedInstances 查找已停止的实例
func findStoppedInstances(accessKeyId, accessKeySecret, regionId, zoneId string) ([]string, error) {
	ecsConfig := &openapi.Config{
		AccessKeyId:     tea.String(accessKeyId),
		AccessKeySecret: tea.String(accessKeySecret),
		RegionId:        tea.String(regionId),
		Endpoint:        tea.String("ecs." + regionId + ".aliyuncs.com"),
	}

	ecsClient, err := ecs.NewClient(ecsConfig)
	if err != nil {
		return nil, err
	}

	statusRequest := &ecs.DescribeInstanceStatusRequest{
		ZoneId:   tea.String(zoneId),
		RegionId: tea.String(regionId),
	}

	statusResponse, err := ecsClient.DescribeInstanceStatus(statusRequest)
	if err != nil {
		return nil, err
	}

	var stoppedInstances []string
	if statusResponse.Body.InstanceStatuses != nil && statusResponse.Body.InstanceStatuses.InstanceStatus != nil {
		for _, status := range statusResponse.Body.InstanceStatuses.InstanceStatus {
			if status.InstanceId != nil && status.Status != nil {
				instanceStatus := *status.Status
				if instanceStatus == "Stopped" || instanceStatus == "Shutted" {
					stoppedInstances = append(stoppedInstances, *status.InstanceId)
				}
			}
		}
	}

	return stoppedInstances, nil
}
