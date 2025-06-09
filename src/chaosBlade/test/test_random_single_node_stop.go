package main

import (
	"fmt"
	"log"
	"os"
	"time"

	cs "github.com/alibabacloud-go/cs-20151215/v5/client"
	ecs "github.com/alibabacloud-go/ecs-20140526/v4/client"
	openapi "github.com/alibabacloud-go/darabonba-openapi/v2/client"
	"github.com/alibabacloud-go/tea/tea"
)

func main() {
	fmt.Println("🔥 NodeLoss 停止任意实例测试")
	fmt.Println("=============================")

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

	// 2. 随机选择一台机器
	selectedInstance := selectRandomInstance(instanceIds)
	fmt.Printf("🎯 随机选择实例: %s\n", selectedInstance)

	// 3. 检查实例当前状态
	fmt.Println("🔍 步骤2: 检查实例当前状态...")
	status, err := getInstanceStatus(ecsClient, selectedInstance, regionId)
	if err != nil {
		log.Printf("⚠️  无法获取实例状态: %v", err)
	} else {
		fmt.Printf("📊 实例 %s 当前状态: %s\n", selectedInstance, status)
	}

	// 4. 停止选中的实例
	fmt.Println("🔄 步骤3: 停止选中的实例...")
	selectedInstances := []string{selectedInstance}
	err = stopInstances(ecsClient, selectedInstances, regionId)
	if err != nil {
		log.Fatalf("❌ 停止实例失败: %v", err)
	}

	fmt.Printf("✅ 成功发送停止请求，实例: %s\n", selectedInstance)
	fmt.Println("⏱️  预计1分钟内停机完成")
	fmt.Println("📝 请记住这个实例ID，用于后续重启测试")
	
	// 5. 等待几秒后再次检查状态
	fmt.Println("⏳ 等待5秒后检查状态...")
	time.Sleep(5 * time.Second)
	
	newStatus, err := getInstanceStatus(ecsClient, selectedInstance, regionId)
	if err != nil {
		log.Printf("⚠️  无法获取实例状态: %v", err)
	} else {
		fmt.Printf("📊 实例 %s 更新后状态: %s\n", selectedInstance, newStatus)
	}

	fmt.Println("🎉 停止实例测试完成！")
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

// getCurrentZoneInstances 获取当前可用区的实例列表（与主代码逻辑完全一致）
func getCurrentZoneInstances(csClient *cs.Client, ecsClient *ecs.Client, clusterId, zoneId, regionId string) ([]string, error) {
	// 1. 调用 DescribeClusterNodes 获取集群节点（支持分页）
	fmt.Println("🔍 开始获取当前可用区实例列表（支持分页查询）")
	allInstanceIds, err := getAllClusterNodesWithPagination(csClient, clusterId, "")
	if err != nil {
		return nil, fmt.Errorf("获取集群节点失败: %v", err)
	}

	fmt.Printf("✅ 获取到 %d 个集群实例\n", len(allInstanceIds))

	if len(allInstanceIds) == 0 {
		fmt.Println("⚠️  未找到任何实例")
		return []string{}, nil
	}

	// 2. 使用 DescribeInstanceStatus 过滤特定区域的实例
	fmt.Printf("🔍 按可用区 %s 过滤实例\n", zoneId)
	filteredIds, err := filterInstancesByZone(ecsClient, allInstanceIds, zoneId, regionId)
	if err != nil {
		return nil, err
	}
	fmt.Printf("✅ 可用区过滤后剩余 %d 个实例\n", len(filteredIds))
	return filteredIds, nil
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

// selectRandomInstance 从实例列表中随机选择一台机器
func selectRandomInstance(instanceIds []string) string {
	if len(instanceIds) == 0 {
		return ""
	}
	
	if len(instanceIds) == 1 {
		return instanceIds[0]
	}

	// 使用当前时间作为随机种子
	currentTime := time.Now()
	index := int(currentTime.UnixNano()) % len(instanceIds)
	
	fmt.Printf("🎲 随机选择逻辑:\n")
	fmt.Printf("   - 可选实例数量: %d\n", len(instanceIds))
	fmt.Printf("   - 时间种子: %d\n", currentTime.UnixNano())
	fmt.Printf("   - 计算索引: %d\n", index)
	
	return instanceIds[index]
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

// stopInstances 批量停止实例
func stopInstances(ecsClient *ecs.Client, instanceIds []string, regionId string) error {
	if len(instanceIds) == 0 {
		return nil
	}

	fmt.Printf("📋 停止实例参数:\n")
	fmt.Printf("   - 区域: %s\n", regionId)
	fmt.Printf("   - 实例: %v\n", instanceIds)
	fmt.Printf("   - 强制停止: true\n")

	request := &ecs.StopInstancesRequest{
		InstanceId: tea.StringSlice(instanceIds),
		ForceStop:  tea.Bool(true),
		RegionId:   tea.String(regionId),
	}

	response, err := ecsClient.StopInstances(request)
	if err != nil {
		return fmt.Errorf("调用 StopInstances 失败: %v", err)
	}

	if response.Body.RequestId != nil {
		fmt.Printf("📝 请求ID: %s\n", *response.Body.RequestId)
	}

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

// getAllClusterNodesWithPagination 获取集群所有节点，支持分页（与主代码逻辑完全一致）
func getAllClusterNodesWithPagination(csClient *cs.Client, clusterId, nodePoolId string) ([]string, error) {
	var allInstanceIds []string
	pageNumber := 1
	pageSize := 10 // 设置每页大小

	fmt.Printf("🔍 [getAllClusterNodesWithPagination] 开始分页查询集群节点\n")
	if nodePoolId != "" {
		fmt.Printf("🎯 [getAllClusterNodesWithPagination] 限定节点池: %s\n", nodePoolId)
	} else {
		fmt.Printf("🌍 [getAllClusterNodesWithPagination] 查询整个集群节点\n")
	}

	for {
		request := &cs.DescribeClusterNodesRequest{
			PageNumber: tea.String(fmt.Sprintf("%d", pageNumber)),
			PageSize:   tea.String(fmt.Sprintf("%d", pageSize)),
		}
		
		if nodePoolId != "" {
			request.NodepoolId = tea.String(nodePoolId)
		}

		fmt.Printf("📄 [getAllClusterNodesWithPagination] 查询第 %d 页（每页 %d 条）\n", pageNumber, pageSize)
		
		response, err := csClient.DescribeClusterNodes(tea.String(clusterId), request)
		if err != nil {
			return nil, fmt.Errorf("查询第 %d 页失败: %v", pageNumber, err)
		}

		// 处理当前页的数据
		currentPageNodes := 0
		if response.Body.Nodes != nil {
			currentPageNodes = len(response.Body.Nodes)
			for _, node := range response.Body.Nodes {
				if node.InstanceId != nil {
					allInstanceIds = append(allInstanceIds, *node.InstanceId)
				}
			}
		}

		fmt.Printf("   ✅ 第 %d 页返回 %d 个节点\n", pageNumber, currentPageNodes)

		// 检查分页信息
		if response.Body.Page != nil {
			if response.Body.Page.TotalCount != nil {
				fmt.Printf("   📊 总记录数: %d\n", *response.Body.Page.TotalCount)
			}
		}

		// 判断是否还有下一页
		if currentPageNodes < pageSize {
			fmt.Printf("   ✅ 已到最后一页\n")
			break
		}

		pageNumber++
	}

	fmt.Printf("🎉 [getAllClusterNodesWithPagination] 分页查询完成，共获取 %d 个实例\n", len(allInstanceIds))
	return allInstanceIds, nil
} 