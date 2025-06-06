package main

import (
	"fmt"
	"log"
	"os"
	"time"

	cs "github.com/alibabacloud-go/cs-20151215/v5/client"
	openapi "github.com/alibabacloud-go/darabonba-openapi/v2/client"
	ecs "github.com/alibabacloud-go/ecs-20140526/v4/client"
	"github.com/alibabacloud-go/tea/tea"
)

func main() {
	fmt.Println("🔄 NodeLoss 重启已停止实例测试")
	fmt.Println("===============================")

	// 从环境变量获取配置
	accessKeyId := os.Getenv("ALIBABA_CLOUD_ACCESS_KEY_ID")
	accessKeySecret := os.Getenv("ALIBABA_CLOUD_ACCESS_KEY_SECRET")
	regionId := os.Getenv("ALIBABA_CLOUD_REGION_ID")
	clusterId := os.Getenv("ALIBABA_CLOUD_CLUSTER_ID")
	zoneId := os.Getenv("ALIBABA_CLOUD_ZONE_ID")

	// 检查必要的环境变量
	if accessKeyId == "" || accessKeySecret == "" || regionId == "" || clusterId == "" || zoneId == "" {
		log.Fatalf("❌ 缺少必要的环境变量:\n"+
			"   - ALIBABA_CLOUD_ACCESS_KEY_ID: %s\n"+
			"   - ALIBABA_CLOUD_ACCESS_KEY_SECRET: %s\n"+
			"   - ALIBABA_CLOUD_REGION_ID: %s\n"+
			"   - ALIBABA_CLOUD_CLUSTER_ID: %s\n"+
			"   - ALIBABA_CLOUD_ZONE_ID: %s\n",
			maskString(accessKeyId), maskString(accessKeySecret), regionId, clusterId, zoneId)
	}

	fmt.Printf("📋 测试配置:\n")
	fmt.Printf("   - 区域: %s\n", regionId)
	fmt.Printf("   - 集群: %s\n", clusterId)
	fmt.Printf("   - 可用区: %s\n", zoneId)
	fmt.Println()

	// 初始化阿里云客户端
	ecsClient, csClient, err := initClients(accessKeyId, accessKeySecret, regionId)
	if err != nil {
		log.Fatalf("❌ 初始化客户端失败: %v", err)
	}

	fmt.Println("✅ 阿里云客户端初始化成功")

	// 1. 获取当前AZ的实例列表
	fmt.Println("🔍 步骤1: 获取当前AZ的实例列表...")
	instanceIds, err := getCurrentZoneInstances(csClient, ecsClient, clusterId, zoneId, regionId)
	if err != nil {
		log.Fatalf("❌ 获取实例列表失败: %v", err)
	}

	if len(instanceIds) == 0 {
		fmt.Println("⚠️  当前AZ没有可操作的实例")
		return
	}

	fmt.Printf("✅ 找到 %d 个实例: %v\n", len(instanceIds), instanceIds)

	// 2. 检查实例状态，找出已停止的实例
	fmt.Println("🔍 步骤2: 检查实例状态，寻找已停止的实例...")
	stoppedInstances, runningInstances, err := categorizeInstancesByStatus(ecsClient, instanceIds, regionId)
	if err != nil {
		log.Fatalf("❌ 检查实例状态失败: %v", err)
	}

	fmt.Printf("📊 实例状态统计:\n")
	fmt.Printf("   - 运行中: %d 个 %v\n", len(runningInstances), runningInstances)
	fmt.Printf("   - 已停止: %d 个 %v\n", len(stoppedInstances), stoppedInstances)

	if len(stoppedInstances) == 0 {
		fmt.Println("ℹ️  没有找到已停止的实例")
		fmt.Println("💡 提示: 请先运行 test_random_single_node_stop.go 来停止一些实例")
		return
	}

	// 3. 重启已停止的实例
	fmt.Printf("🚀 步骤3: 重启 %d 个已停止的实例...\n", len(stoppedInstances))
	err = rebootInstances(ecsClient, stoppedInstances, regionId)
	if err != nil {
		log.Fatalf("❌ 重启实例失败: %v", err)
	}

	fmt.Printf("✅ 成功发送重启请求，实例: %v\n", stoppedInstances)
	fmt.Println("⏱️  预计1-2分钟内重启完成")

	// 4. 等待几秒后再次检查状态
	fmt.Println("⏳ 等待10秒后检查状态...")
	time.Sleep(10 * time.Second)

	fmt.Println("📊 重启后实例状态:")
	for _, instanceId := range stoppedInstances {
		status, err := getInstanceStatus(ecsClient, instanceId, regionId)
		if err != nil {
			log.Printf("⚠️  无法获取实例 %s 状态: %v", instanceId, err)
		} else {
			statusIcon := "🟡"
			if status == "Running" {
				statusIcon = "🟢"
			} else if status == "Stopped" {
				statusIcon = "🔴"
			}
			fmt.Printf("   %s %s: %s\n", statusIcon, instanceId, status)
		}
	}

	fmt.Println("🎉 重启实例测试完成！")
}

// initClients 初始化阿里云客户端
func initClients(accessKeyId, accessKeySecret, regionId string) (*ecs.Client, *cs.Client, error) {
	// 创建 ECS 客户端配置
	ecsConfig := &openapi.Config{
		AccessKeyId:     tea.String(accessKeyId),
		AccessKeySecret: tea.String(accessKeySecret),
		RegionId:        tea.String(regionId),
		Endpoint:        tea.String("ecs." + regionId + ".aliyuncs.com"),
	}

	// 创建 CS 客户端配置
	csConfig := &openapi.Config{
		AccessKeyId:     tea.String(accessKeyId),
		AccessKeySecret: tea.String(accessKeySecret),
		RegionId:        tea.String(regionId),
		Endpoint:        tea.String("cs." + regionId + ".aliyuncs.com"),
	}

	// 初始化 ECS 客户端
	ecsClient, err := ecs.NewClient(ecsConfig)
	if err != nil {
		return nil, nil, fmt.Errorf("初始化 ECS 客户端失败: %v", err)
	}

	// 初始化 CS 客户端
	csClient, err := cs.NewClient(csConfig)
	if err != nil {
		return nil, nil, fmt.Errorf("初始化 CS 客户端失败: %v", err)
	}

	return ecsClient, csClient, nil
}

// getCurrentZoneInstances 获取当前可用区的实例列表
func getCurrentZoneInstances(csClient *cs.Client, ecsClient *ecs.Client, clusterId, zoneId, regionId string) ([]string, error) {
	// 1. 调用 DescribeClusterNodes 获取集群节点
	request := &cs.DescribeClusterNodesRequest{}
	response, err := csClient.DescribeClusterNodes(tea.String(clusterId), request)
	if err != nil {
		return nil, fmt.Errorf("调用 DescribeClusterNodes 失败: %v", err)
	}

	var allInstanceIds []string
	if response.Body.Nodes != nil {
		for _, node := range response.Body.Nodes {
			if node.InstanceId != nil {
				allInstanceIds = append(allInstanceIds, *node.InstanceId)
			}
		}
	}

	if len(allInstanceIds) == 0 {
		return []string{}, nil
	}

	// 2. 使用 DescribeInstanceStatus 过滤特定区域的实例
	return filterInstancesByZone(ecsClient, allInstanceIds, zoneId, regionId)
}

// filterInstancesByZone 根据可用区过滤实例
func filterInstancesByZone(ecsClient *ecs.Client, instanceIds []string, zoneId, regionId string) ([]string, error) {
	if len(instanceIds) == 0 {
		return []string{}, nil
	}

	// 调用 DescribeInstanceStatus
	request := &ecs.DescribeInstanceStatusRequest{
		InstanceId: tea.StringSlice(instanceIds),
		ZoneId:     tea.String(zoneId),
		RegionId:   tea.String(regionId),
	}

	response, err := ecsClient.DescribeInstanceStatus(request)
	if err != nil {
		return nil, fmt.Errorf("调用 DescribeInstanceStatus 失败: %v", err)
	}

	var filteredIds []string
	if response.Body.InstanceStatuses != nil && response.Body.InstanceStatuses.InstanceStatus != nil {
		for _, status := range response.Body.InstanceStatuses.InstanceStatus {
			if status.InstanceId != nil {
				filteredIds = append(filteredIds, *status.InstanceId)
			}
		}
	}

	return filteredIds, nil
}

// categorizeInstancesByStatus 根据状态分类实例
func categorizeInstancesByStatus(ecsClient *ecs.Client, instanceIds []string, regionId string) ([]string, []string, error) {
	request := &ecs.DescribeInstanceStatusRequest{
		InstanceId: tea.StringSlice(instanceIds),
		RegionId:   tea.String(regionId),
	}

	response, err := ecsClient.DescribeInstanceStatus(request)
	if err != nil {
		return nil, nil, fmt.Errorf("检查实例状态失败: %v", err)
	}

	var stoppedInstances []string // 已停止的实例
	var runningInstances []string // 运行中的实例

	if response.Body.InstanceStatuses != nil && response.Body.InstanceStatuses.InstanceStatus != nil {
		for _, status := range response.Body.InstanceStatuses.InstanceStatus {
			if status.InstanceId != nil && status.Status != nil {
				instanceId := *status.InstanceId
				instanceStatus := *status.Status

				switch instanceStatus {
				case "Stopped", "Shutted":
					stoppedInstances = append(stoppedInstances, instanceId)
				case "Running":
					runningInstances = append(runningInstances, instanceId)
				default:
					fmt.Printf("⚠️  实例 %s 状态异常: %s\n", instanceId, instanceStatus)
				}
			}
		}
	}

	return stoppedInstances, runningInstances, nil
}

// getInstanceStatus 获取实例状态
func getInstanceStatus(ecsClient *ecs.Client, instanceId, regionId string) (string, error) {
	request := &ecs.DescribeInstanceStatusRequest{
		InstanceId: tea.StringSlice([]string{instanceId}),
		RegionId:   tea.String(regionId),
	}

	response, err := ecsClient.DescribeInstanceStatus(request)
	if err != nil {
		return "", fmt.Errorf("获取实例状态失败: %v", err)
	}

	if response.Body.InstanceStatuses != nil &&
		response.Body.InstanceStatuses.InstanceStatus != nil &&
		len(response.Body.InstanceStatuses.InstanceStatus) > 0 &&
		response.Body.InstanceStatuses.InstanceStatus[0].Status != nil {
		return *response.Body.InstanceStatuses.InstanceStatus[0].Status, nil
	}

	return "Unknown", nil
}

// rebootInstances 批量重启实例（智能选择启动或重启）
func rebootInstances(ecsClient *ecs.Client, instanceIds []string, regionId string) error {
	if len(instanceIds) == 0 {
		return nil
	}

	fmt.Printf("📋 重启实例参数:\n")
	fmt.Printf("   - 区域: %s\n", regionId)
	fmt.Printf("   - 实例: %v\n", instanceIds)

	// 检查实例状态，决定使用 StartInstances 还是 RebootInstances
	fmt.Println("🔍 检查实例当前状态...")
	statusRequest := &ecs.DescribeInstanceStatusRequest{
		InstanceId: tea.StringSlice(instanceIds),
		RegionId:   tea.String(regionId),
	}

	statusResponse, err := ecsClient.DescribeInstanceStatus(statusRequest)
	if err != nil {
		return fmt.Errorf("检查实例状态失败: %v", err)
	}

	var stoppedInstances []string // 需要启动的实例
	var runningInstances []string // 需要重启的实例

	if statusResponse.Body.InstanceStatuses != nil && statusResponse.Body.InstanceStatuses.InstanceStatus != nil {
		for _, status := range statusResponse.Body.InstanceStatuses.InstanceStatus {
			if status.InstanceId != nil && status.Status != nil {
				instanceId := *status.InstanceId
				instanceStatus := *status.Status

				switch instanceStatus {
				case "Stopped", "Shutted":
					stoppedInstances = append(stoppedInstances, instanceId)
				case "Running":
					runningInstances = append(runningInstances, instanceId)
				default:
					fmt.Printf("⚠️  实例 %s 状态异常: %s\n", instanceId, instanceStatus)
				}
			}
		}
	}

	fmt.Printf("🎯 操作策略:\n")
	fmt.Printf("   - 需要启动的实例: %d 个 %v\n", len(stoppedInstances), stoppedInstances)
	fmt.Printf("   - 需要重启的实例: %d 个 %v\n", len(runningInstances), runningInstances)

	// 启动已停止的实例
	if len(stoppedInstances) > 0 {
		fmt.Printf("🚀 启动已停止的实例...\n")

		startRequest := &ecs.StartInstancesRequest{
			InstanceId: tea.StringSlice(stoppedInstances),
			RegionId:   tea.String(regionId),
		}

		startResponse, err := ecsClient.StartInstances(startRequest)
		if err != nil {
			return fmt.Errorf("调用 StartInstances 失败: %v", err)
		}

		if startResponse.Body.RequestId != nil {
			fmt.Printf("📝 启动请求ID: %s\n", *startResponse.Body.RequestId)
		}
	}

	// 重启正在运行的实例
	if len(runningInstances) > 0 {
		fmt.Printf("🔄 重启运行中的实例...\n")

		rebootRequest := &ecs.RebootInstancesRequest{
			InstanceId: tea.StringSlice(runningInstances),
			RegionId:   tea.String(regionId),
		}

		rebootResponse, err := ecsClient.RebootInstances(rebootRequest)
		if err != nil {
			return fmt.Errorf("调用 RebootInstances 失败: %v", err)
		}

		if rebootResponse.Body.RequestId != nil {
			fmt.Printf("📝 重启请求ID: %s\n", *rebootResponse.Body.RequestId)
		}
	}

	totalProcessed := len(stoppedInstances) + len(runningInstances)
	fmt.Printf("🎉 批量重启操作完成，处理了 %d 个实例 (启动: %d, 重启: %d)\n",
		totalProcessed, len(stoppedInstances), len(runningInstances))

	return nil
}

// maskString 掩码字符串，用于安全显示敏感信息
func maskString(s string) string {
	if s == "" {
		return "未设置"
	}
	if len(s) <= 8 {
		return "***"
	}
	return s[:4] + "****" + s[len(s)-4:]
}
