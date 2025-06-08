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
	fmt.Println("🔍 查询集群所有机器（不限制节点池）")
	fmt.Println("=====================================")

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

	// 方案1：查询整个集群所有节点（不指定节点池）
	fmt.Println("\n=== 方案1：查询整个集群所有节点 ===")
	allNodesRequest := &cs.DescribeClusterNodesRequest{}
	// 不设置 NodepoolId，获取整个集群的所有节点
	allNodesResponse, err := csClient.DescribeClusterNodes(tea.String(clusterId), allNodesRequest)

	var allClusterInstanceIds []string

	if err != nil {
		log.Printf("❌ 查询集群所有节点失败: %v", err)
	} else {
		nodeCount := 0
		if allNodesResponse.Body.Nodes != nil {
			nodeCount = len(allNodesResponse.Body.Nodes)
			for _, node := range allNodesResponse.Body.Nodes {
				if node.InstanceId != nil {
					allClusterInstanceIds = append(allClusterInstanceIds, *node.InstanceId)
				}
			}
		}
		fmt.Printf("✅ 整个集群共有 %d 个节点\n", nodeCount)
		fmt.Printf("📋 所有实例ID: %v\n", allClusterInstanceIds)
	}

	// 方案2：查询指定节点池的节点
	fmt.Println("\n=== 方案2：查询指定节点池的节点 ===")
	poolNodesRequest := &cs.DescribeClusterNodesRequest{
		NodepoolId: tea.String(nodePoolId),
	}
	poolNodesResponse, err := csClient.DescribeClusterNodes(tea.String(clusterId), poolNodesRequest)

	var poolInstanceIds []string

	if err != nil {
		log.Printf("❌ 查询节点池节点失败: %v", err)
	} else {
		nodeCount := 0
		if poolNodesResponse.Body.Nodes != nil {
			nodeCount = len(poolNodesResponse.Body.Nodes)
			for _, node := range poolNodesResponse.Body.Nodes {
				if node.InstanceId != nil {
					poolInstanceIds = append(poolInstanceIds, *node.InstanceId)
				}
			}
		}
		fmt.Printf("✅ 节点池 %s 共有 %d 个节点\n", nodePoolId, nodeCount)
		fmt.Printf("📋 节点池实例ID: %v\n", poolInstanceIds)
	}

	// 对比分析
	fmt.Println("\n=== 对比分析 ===")
	fmt.Printf("🔍 整个集群节点数: %d\n", len(allClusterInstanceIds))
	fmt.Printf("🔍 指定节点池节点数: %d\n", len(poolInstanceIds))
	fmt.Printf("🔍 差异节点数: %d\n", len(allClusterInstanceIds)-len(poolInstanceIds))

	// 找出不在指定节点池中的节点
	poolMap := make(map[string]bool)
	for _, id := range poolInstanceIds {
		poolMap[id] = true
	}

	var missingInstances []string
	for _, id := range allClusterInstanceIds {
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
	if len(allClusterInstanceIds) > 0 {
		fmt.Println("\n=== 方案3：按可用区过滤所有集群节点 ===")

		// 使用所有集群实例进行可用区过滤
		filterRequest := &ecs.DescribeInstanceStatusRequest{
			InstanceId: tea.StringSlice(allClusterInstanceIds),
			ZoneId:     tea.String(zoneId),
			RegionId:   tea.String(regionId),
		}

		filterResponse, err := ecsClient.DescribeInstanceStatus(filterRequest)
		var allZoneFilteredIds []string

		if err != nil {
			log.Printf("❌ 按可用区过滤失败: %v", err)
		} else {
			filteredCount := 0
			if filterResponse.Body.InstanceStatuses != nil && filterResponse.Body.InstanceStatuses.InstanceStatus != nil {
				filteredCount = len(filterResponse.Body.InstanceStatuses.InstanceStatus)
				for _, status := range filterResponse.Body.InstanceStatuses.InstanceStatus {
					if status.InstanceId != nil {
						allZoneFilteredIds = append(allZoneFilteredIds, *status.InstanceId)
					}
				}
			}
			fmt.Printf("✅ 在可用区 %s 中过滤出 %d 个实例（基于所有集群节点）\n", zoneId, filteredCount)
			fmt.Printf("📋 可用区过滤后的实例: %v\n", allZoneFilteredIds)
		}

		// 对比不同过滤方式的结果
		fmt.Println("\n=== 最终对比 ===")
		fmt.Printf("🔍 所有集群节点: %d 个\n", len(allClusterInstanceIds))
		fmt.Printf("🔍 指定节点池节点: %d 个\n", len(poolInstanceIds))
		fmt.Printf("🔍 指定节点池+可用区过滤: %d 个\n", 6) // 从日志中得知
		fmt.Printf("🔍 所有集群节点+可用区过滤: %d 个\n", len(allZoneFilteredIds))

		fmt.Println("\n💡 结论:")
		if len(allClusterInstanceIds) > len(poolInstanceIds) {
			fmt.Printf("   - 集群中确实有 %d 个节点不在指定节点池中\n", len(allClusterInstanceIds)-len(poolInstanceIds))
			fmt.Printf("   - 如果要进行区域级故障注入，应该使用所有集群节点而不是仅限节点池\n")
		}

		if len(allZoneFilteredIds) > 6 {
			fmt.Printf("   - 在可用区 %s 中实际应该有 %d 个可操作实例，而不是 6 个\n", zoneId, len(allZoneFilteredIds))
		}
	}

	fmt.Println("\n🎉 分析完成!")
}
