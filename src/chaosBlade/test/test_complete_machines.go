package main

import (
	"fmt"
	cs "github.com/alibabacloud-go/cs-20151215/v5/client"
	openapi "github.com/alibabacloud-go/darabonba-openapi/v2/client"
	ecs "github.com/alibabacloud-go/ecs-20140526/v4/client"
	"github.com/alibabacloud-go/tea/tea"
	"log"
)

func main() {
	fmt.Println("🔍 完整查询集群所有机器（支持分页）")
	fmt.Println("======================================")

	// 硬编码配置（用于测试）
	// 从环境变量获取配置
	accessKeyId := os.Getenv("ALIBABA_CLOUD_ACCESS_KEY_ID")
	accessKeySecret := os.Getenv("ALIBABA_CLOUD_ACCESS_KEY_SECRET")
	regionId := os.Getenv("ALIBABA_CLOUD_REGION_ID")
	clusterId := os.Getenv("ALIBABA_CLOUD_CLUSTER_ID")
	zoneId := os.Getenv("ALIBABA_CLOUD_ZONE_ID")

	fmt.Printf("📋 配置信息:\n")
	fmt.Printf("   - 区域: %s\n", regionId)
	fmt.Printf("   - 集群ID: %s\n", clusterId)
	fmt.Printf("   - 节点池ID: %s\n", nodePoolId)
	fmt.Printf("   - 可用区: %s\n", zoneId)
	fmt.Println()

	// 创建客户端
	ecsConfig := &openapi.Config{
		AccessKeyId:     tea.String(accessKeyId),
		AccessKeySecret: tea.String(accessKeySecret),
		RegionId:        tea.String(regionId),
		Endpoint:        tea.String("ecs." + regionId + ".aliyuncs.com"),
	}

	csConfig := &openapi.Config{
		AccessKeyId:     tea.String(accessKeyId),
		AccessKeySecret: tea.String(accessKeySecret),
		RegionId:        tea.String(regionId),
		Endpoint:        tea.String("cs." + regionId + ".aliyuncs.com"),
	}

	ecsClient, err := ecs.NewClient(ecsConfig)
	if err != nil {
		log.Fatalf("❌ 创建 ECS 客户端失败: %v", err)
	}

	csClient, err := cs.NewClient(csConfig)
	if err != nil {
		log.Fatalf("❌ 创建 CS 客户端失败: %v", err)
	}

	fmt.Println("✅ 客户端创建成功")

	// 方案1：支持分页的完整集群节点查询
	fmt.Println("\n=== 方案1：分页查询整个集群所有节点 ===")
	allClusterInstances := getAllClusterNodes(csClient, clusterId, "")
	fmt.Printf("✅ 整个集群共有 %d 个节点\n", len(allClusterInstances))
	fmt.Printf("📋 所有实例ID: %v\n", allClusterInstances)

	// 方案2：支持分页的指定节点池查询
	fmt.Println("\n=== 方案2：分页查询指定节点池节点 ===")
	poolInstances := getAllClusterNodes(csClient, clusterId, nodePoolId)
	fmt.Printf("✅ 节点池 %s 共有 %d 个节点\n", nodePoolId, len(poolInstances))
	fmt.Printf("📋 节点池实例ID: %v\n", poolInstances)

	// 对比分析
	fmt.Println("\n=== 对比分析 ===")
	fmt.Printf("🔍 整个集群节点数: %d\n", len(allClusterInstances))
	fmt.Printf("🔍 指定节点池节点数: %d\n", len(poolInstances))
	fmt.Printf("🔍 差异节点数: %d\n", len(allClusterInstances)-len(poolInstances))

	// 找出不在指定节点池中的节点
	poolMap := make(map[string]bool)
	for _, id := range poolInstances {
		poolMap[id] = true
	}

	var missingInstances []string
	for _, id := range allClusterInstances {
		if !poolMap[id] {
			missingInstances = append(missingInstances, id)
		}
	}

	if len(missingInstances) > 0 {
		fmt.Printf("\n❗ 发现 %d 个不在指定节点池中的实例:\n", len(missingInstances))
		for i, id := range missingInstances {
			fmt.Printf("   %d. %s\n", i+1, id)
		}
	}

	// 方案3：按可用区过滤所有集群节点
	if len(allClusterInstances) > 0 {
		fmt.Println("\n=== 方案3：按可用区过滤所有集群节点 ===")
		
		allZoneFilteredIds := filterInstancesByZone(ecsClient, allClusterInstances, zoneId, regionId)
		fmt.Printf("✅ 在可用区 %s 中过滤出 %d 个实例（基于所有集群节点）\n", zoneId, len(allZoneFilteredIds))
		fmt.Printf("📋 可用区过滤后的实例: %v\n", allZoneFilteredIds)

		// 同样过滤节点池的实例
		poolZoneFilteredIds := filterInstancesByZone(ecsClient, poolInstances, zoneId, regionId)
		fmt.Printf("✅ 在可用区 %s 中过滤出 %d 个实例（基于指定节点池）\n", zoneId, len(poolZoneFilteredIds))
		fmt.Printf("📋 节点池+可用区过滤后的实例: %v\n", poolZoneFilteredIds)

		// 对比不同过滤方式的结果
		fmt.Println("\n=== 最终对比 ===")
		fmt.Printf("🔍 所有集群节点: %d 个\n", len(allClusterInstances))
		fmt.Printf("🔍 指定节点池节点: %d 个\n", len(poolInstances))
		fmt.Printf("🔍 指定节点池+可用区过滤: %d 个\n", len(poolZoneFilteredIds))
		fmt.Printf("🔍 所有集群节点+可用区过滤: %d 个\n", len(allZoneFilteredIds))
		
		fmt.Println("\n💡 结论:")
		if len(allClusterInstances) > len(poolInstances) {
			fmt.Printf("   - 集群中确实有 %d 个节点不在指定节点池中\n", len(allClusterInstances)-len(poolInstances))
			fmt.Printf("   - 如果要进行区域级故障注入，应该使用所有集群节点而不是仅限节点池\n")
		}
		
		if len(allZoneFilteredIds) > len(poolZoneFilteredIds) {
			fmt.Printf("   - 在可用区 %s 中实际应该有 %d 个可操作实例，而不是 %d 个\n", 
				zoneId, len(allZoneFilteredIds), len(poolZoneFilteredIds))
			fmt.Printf("   - 这就是为什么您觉得'少停止了一台机器'的原因！\n")
		}
	}

	fmt.Println("\n🎉 分析完成!")
}

// getAllClusterNodes 获取集群所有节点，支持分页
func getAllClusterNodes(csClient *cs.Client, clusterId, nodePoolId string) []string {
	var allInstanceIds []string
	pageNumber := 1
	pageSize := 10 // 设置每页大小

	for {
		request := &cs.DescribeClusterNodesRequest{
			PageNumber: tea.String(fmt.Sprintf("%d", pageNumber)),
			PageSize:   tea.String(fmt.Sprintf("%d", pageSize)),
		}
		
		if nodePoolId != "" {
			request.NodepoolId = tea.String(nodePoolId)
		}

		fmt.Printf("📄 查询第 %d 页（每页 %d 条）...\n", pageNumber, pageSize)
		
		response, err := csClient.DescribeClusterNodes(tea.String(clusterId), request)
		if err != nil {
			log.Printf("❌ 查询第 %d 页失败: %v", pageNumber, err)
			break
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
			if response.Body.Page.PageNumber != nil {
				fmt.Printf("   📊 当前页: %d\n", *response.Body.Page.PageNumber)
			}
			if response.Body.Page.PageSize != nil {
				fmt.Printf("   📊 页大小: %d\n", *response.Body.Page.PageSize)
			}
		}

		// 判断是否还有下一页
		if currentPageNodes < pageSize {
			fmt.Printf("   ✅ 已到最后一页\n")
			break
		}

		pageNumber++
	}

	return allInstanceIds
}

// filterInstancesByZone 根据可用区过滤实例
func filterInstancesByZone(ecsClient *ecs.Client, instanceIds []string, zoneId, regionId string) []string {
	if len(instanceIds) == 0 {
		return []string{}
	}

	request := &ecs.DescribeInstanceStatusRequest{
		InstanceId: tea.StringSlice(instanceIds),
		ZoneId:     tea.String(zoneId),
		RegionId:   tea.String(regionId),
	}

	response, err := ecsClient.DescribeInstanceStatus(request)
	if err != nil {
		log.Printf("❌ 按可用区过滤失败: %v", err)
		return []string{}
	}

	var filteredIds []string
	if response.Body.InstanceStatuses != nil && response.Body.InstanceStatuses.InstanceStatus != nil {
		for _, status := range response.Body.InstanceStatuses.InstanceStatus {
			if status.InstanceId != nil {
				filteredIds = append(filteredIds, *status.InstanceId)
			}
		}
	}

	return filteredIds
} 