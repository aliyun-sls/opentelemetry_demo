package command

import (
	"fmt"
	cs "github.com/alibabacloud-go/cs-20151215/v5/client"
	ecs "github.com/alibabacloud-go/ecs-20140526/v4/client"
	"github.com/alibabacloud-go/tea/tea"
	"log"
	"os"
	"strconv"
	"time"
)

type AliyunRegionChaos struct {
	lastFlagdValue      int64
	flagdValue          int64
	instanceIds         []string
	stoppedInstances    []string
	lastOperationTime   time.Time
	minOperationInterval time.Duration
}

func NewAliyunRegionChaos() *AliyunRegionChaos {
	// 🔧 初始化阿里云客户端（一次性，在创建实例时初始化，避免每次ExecuteChaos都调用）
	InitAliyunClients()
	
	// 从环境变量获取最小操作间隔（单位：分钟），默认3分钟
	minIntervalMinutes := 3 // 默认值
	if intervalEnv := os.Getenv("MIN_OPERATION_INTERVAL_MINUTES"); intervalEnv != "" {
		if parsed, err := strconv.Atoi(intervalEnv); err == nil && parsed > 0 {
			minIntervalMinutes = parsed
			log.Printf("📝 [NewAliyunRegionChaos] 从环境变量读取最小操作间隔: %d 分钟", minIntervalMinutes)
		} else {
			log.Printf("⚠️  [NewAliyunRegionChaos] 环境变量 MIN_OPERATION_INTERVAL_MINUTES 无效 (%s)，使用默认值: %d 分钟", 
				intervalEnv, minIntervalMinutes)
		}
	} else {
		log.Printf("ℹ️  [NewAliyunRegionChaos] 未设置 MIN_OPERATION_INTERVAL_MINUTES，使用默认值: %d 分钟", minIntervalMinutes)
	}

	minInterval := time.Duration(minIntervalMinutes) * time.Minute
	
	return &AliyunRegionChaos{
		lastFlagdValue:       0,
		flagdValue:           0,
		instanceIds:          []string{},
		stoppedInstances:     []string{},
		lastOperationTime:    time.Time{}, // 零值表示没有操作过
		minOperationInterval: minInterval, // 从环境变量配置的最小操作间隔
	}
}

func (chaos *AliyunRegionChaos) ExecuteChaos(flagdValue int64) {
	// 客户端已在 NewAliyunRegionChaos() 中初始化，无需重复调用
	currentTime := time.Now()
	
	log.Printf("🔍 [ExecuteChaos] 收到请求 - 当前时间: %s", currentTime.Format("2006-01-02 15:04:05"))
	log.Printf("📊 [ExecuteChaos] 当前状态 - flagd值: %d → %d, 已停止实例数: %d", 
		chaos.lastFlagdValue, flagdValue, len(chaos.stoppedInstances))
	
	chaos.flagdValue = flagdValue
	if chaos.flagdValue == chaos.lastFlagdValue {
		log.Printf("⏭️  [ExecuteChaos] flagd值未变化 (%d)，跳过操作", flagdValue)
		return
	}

	// 检查时间间隔保护：如果不是第一次操作且距离上次操作不足3分钟，则忽略此次请求
	if !chaos.lastOperationTime.IsZero() {
		timeSinceLastOp := currentTime.Sub(chaos.lastOperationTime)
		nextAllowedTime := chaos.lastOperationTime.Add(chaos.minOperationInterval)
		
		log.Printf("🕐 [时间检查] 上次操作时间: %s", chaos.lastOperationTime.Format("2006-01-02 15:04:05"))
		log.Printf("⏰ [时间检查] 距离上次操作: %.1f 秒, 要求间隔: %.0f 秒 (%.1f 分钟)", 
			timeSinceLastOp.Seconds(), chaos.minOperationInterval.Seconds(), chaos.minOperationInterval.Minutes())
		log.Printf("🚀 [时间检查] 下次可操作时间: %s", nextAllowedTime.Format("2006-01-02 15:04:05"))
		
		if timeSinceLastOp < chaos.minOperationInterval {
			remainingTime := chaos.minOperationInterval - timeSinceLastOp
			log.Printf("❌ [操作拒绝] 时间间隔不足，拒绝执行")
			log.Printf("⚠️  [操作拒绝] flagd变更: %d → %d", chaos.lastFlagdValue, chaos.flagdValue)
			log.Printf("⌛ [操作拒绝] 还需等待: %.1f 秒 (到 %s)", 
				remainingTime.Seconds(), nextAllowedTime.Format("15:04:05"))
			log.Printf("🔒 [操作拒绝] 实例状态保持不变，当前已停止: %v", chaos.stoppedInstances)
			return
		} else {
			log.Printf("✅ [时间检查] 间隔足够，允许执行操作")
		}
	} else {
		log.Printf("🎯 [时间检查] 首次操作，无需检查时间间隔")
	}

	log.Printf("🚀 [开始执行] AliyunRegionChaos 执行操作: flagd值 %d → %d", chaos.lastFlagdValue, chaos.flagdValue)

	var operationSuccess bool
	if chaos.flagdValue > 0 {
		log.Printf("🔥 [故障注入] 开始故障注入，持续时间: %d 秒", chaos.flagdValue)
		operationSuccess = chaos.startRegionChaos()
	} else {
		log.Printf("🔄 [故障恢复] 停止故障注入，恢复实例")
		operationSuccess = chaos.stopRegionChaos()
	}

	if operationSuccess {
		// 更新最后操作时间和flagd值
		chaos.lastOperationTime = currentTime
		chaos.lastFlagdValue = chaos.flagdValue
		nextAllowedTime := currentTime.Add(chaos.minOperationInterval)
		
		log.Printf("✅ [操作完成] 成功执行操作，状态已更新")
		log.Printf("📝 [操作完成] 新的lastFlagdValue: %d", chaos.lastFlagdValue)
		log.Printf("🕐 [操作完成] 操作时间已记录: %s", currentTime.Format("2006-01-02 15:04:05"))
		log.Printf("⏰ [操作完成] 下次可操作时间: %s", nextAllowedTime.Format("2006-01-02 15:04:05"))
		log.Printf("📊 [操作完成] 当前已停止实例: %v", chaos.stoppedInstances)
	} else {
		log.Printf("❌ [操作失败] 操作执行失败，状态未更新")
		log.Printf("🔍 [操作失败] 保持原有状态: lastFlagdValue=%d, 已停止实例=%v", 
			chaos.lastFlagdValue, chaos.stoppedInstances)
	}
}

// startRegionChaos 开始区域故障注入
func (chaos *AliyunRegionChaos) startRegionChaos() bool {
	operationTime := time.Now()
	log.Printf("🔥 [startRegionChaos] 开始阿里云区域故障注入")
	log.Printf("🕐 [startRegionChaos] 操作开始时间: %s", operationTime.Format("2006-01-02 15:04:05"))
	log.Printf("📋 [startRegionChaos] 当前实例状态 - instanceIds: %v, stoppedInstances: %v", 
		chaos.instanceIds, chaos.stoppedInstances)

	// 获取集群实例列表
	log.Printf("🔍 [startRegionChaos] 正在获取集群实例列表...")
	instanceIds, err := chaos.getClusterInstances()
	if err != nil {
		log.Printf("❌ [startRegionChaos] 获取集群实例失败: %v", err)
		return false
	}

	if len(instanceIds) == 0 {
		log.Printf("⚠️  [startRegionChaos] 未找到可操作的实例")
		log.Printf("🔍 [startRegionChaos] 请检查环境变量: CLUSTER_ID, NODEPOOL_ID, ZONE_ID")
		return false
	}

	chaos.instanceIds = instanceIds
	log.Printf("✅ [startRegionChaos] 找到 %d 个实例: %v", len(instanceIds), instanceIds)

	// 停止实例
	log.Printf("🔄 [startRegionChaos] 开始停止实例...")
	err = chaos.stopInstances(instanceIds)
	if err != nil {
		log.Printf("❌ [startRegionChaos] 停止实例失败: %v", err)
		log.Printf("🔍 [startRegionChaos] 实例状态未改变，保持原状态")
		return false
	}

	chaos.stoppedInstances = instanceIds
	estimatedCompleteTime := operationTime.Add(1 * time.Minute)
	log.Printf("✅ [startRegionChaos] 成功发送停止请求，影响 %d 个实例", len(instanceIds))
	log.Printf("⏱️  [startRegionChaos] 预计停机完成时间: %s (约1分钟后)", estimatedCompleteTime.Format("15:04:05"))
	log.Printf("📊 [startRegionChaos] 已停止实例列表: %v", chaos.stoppedInstances)
	
	return true
}

// stopRegionChaos 停止区域故障注入
func (chaos *AliyunRegionChaos) stopRegionChaos() bool {
	operationTime := time.Now()
	log.Printf("🔄 [stopRegionChaos] 开始停止区域故障注入")
	log.Printf("🕐 [stopRegionChaos] 操作开始时间: %s", operationTime.Format("2006-01-02 15:04:05"))
	log.Printf("📋 [stopRegionChaos] 当前已停止实例: %v", chaos.stoppedInstances)
	
	if len(chaos.stoppedInstances) == 0 {
		log.Printf("ℹ️  [stopRegionChaos] 没有需要重启的实例")
		log.Printf("✅ [stopRegionChaos] 操作完成（无需重启）")
		return true
	}

	log.Printf("🚀 [stopRegionChaos] 准备重启 %d 个实例: %v", 
		len(chaos.stoppedInstances), chaos.stoppedInstances)

	// 重启实例
	log.Printf("🔄 [stopRegionChaos] 开始重启实例...")
	err := chaos.rebootInstances(chaos.stoppedInstances)
	if err != nil {
		log.Printf("❌ [stopRegionChaos] 重启实例失败: %v", err)
		log.Printf("🔍 [stopRegionChaos] 实例状态未改变，仍保持停止状态: %v", chaos.stoppedInstances)
		return false
	}

	estimatedCompleteTime := operationTime.Add(1 * time.Minute)
	log.Printf("✅ [stopRegionChaos] 成功发送重启请求")
	log.Printf("⏱️  [stopRegionChaos] 预计重启完成时间: %s (约1分钟后)", estimatedCompleteTime.Format("15:04:05"))
	log.Printf("📊 [stopRegionChaos] 清空已停止实例列表")
	
	chaos.stoppedInstances = []string{}
	return true
}

// getClusterInstances 获取集群实例列表
func (chaos *AliyunRegionChaos) getClusterInstances() ([]string, error) {
	clusterId := os.Getenv("CLUSTER_ID")
	nodePoolId := os.Getenv("NODEPOOL_ID")
	zoneId := os.Getenv("ZONE_ID")

	if clusterId == "" {
		return nil, fmt.Errorf("CLUSTER_ID 环境变量未设置")
	}

	// 1. 调用 DescribeClusterNodes 获取集群节点（支持分页）
	log.Printf("🔍 [getClusterInstances] 开始获取集群实例列表（支持分页查询）")
	allInstanceIds, err := chaos.getAllClusterNodesWithPagination(clusterId, nodePoolId)
	if err != nil {
		return nil, fmt.Errorf("获取集群节点失败: %v", err)
	}

	log.Printf("✅ [getClusterInstances] 获取到 %d 个集群实例", len(allInstanceIds))
	log.Printf("📋 [getClusterInstances] 原始实例列表: %v", allInstanceIds)

	if len(allInstanceIds) == 0 {
		log.Printf("⚠️  [getClusterInstances] 未找到任何实例")
		return []string{}, nil
	}

	// 2. 如果指定了 ZONE_ID，使用 DescribeInstanceStatus 过滤特定区域的实例
	if zoneId != "" {
		log.Printf("🔍 [getClusterInstances] 按可用区 %s 过滤实例", zoneId)
		filteredIds, err := chaos.filterInstancesByZone(allInstanceIds, zoneId)
		if err != nil {
			return nil, err
		}
		log.Printf("✅ [getClusterInstances] 可用区过滤后剩余 %d 个实例", len(filteredIds))
		log.Printf("📋 [getClusterInstances] 可用区过滤后实例列表: %v", filteredIds)
		return filteredIds, nil
	}

	return allInstanceIds, nil
}

// getAllClusterNodesWithPagination 获取集群所有节点，支持分页
func (chaos *AliyunRegionChaos) getAllClusterNodesWithPagination(clusterId, nodePoolId string) ([]string, error) {
	var allInstanceIds []string
	pageNumber := 1
	pageSize := 10 // 设置每页大小

	log.Printf("🔍 [getAllClusterNodesWithPagination] 开始分页查询集群节点")
	if nodePoolId != "" {
		log.Printf("🎯 [getAllClusterNodesWithPagination] 限定节点池: %s", nodePoolId)
	} else {
		log.Printf("🌍 [getAllClusterNodesWithPagination] 查询整个集群节点")
	}

	for {
		request := &cs.DescribeClusterNodesRequest{
			PageNumber: tea.String(fmt.Sprintf("%d", pageNumber)),
			PageSize:   tea.String(fmt.Sprintf("%d", pageSize)),
		}
		
		if nodePoolId != "" {
			request.NodepoolId = tea.String(nodePoolId)
		}

		log.Printf("📄 [getAllClusterNodesWithPagination] 查询第 %d 页（每页 %d 条）", pageNumber, pageSize)
		
		response, err := CsClient.DescribeClusterNodes(tea.String(clusterId), request)
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

		log.Printf("   ✅ 第 %d 页返回 %d 个节点", pageNumber, currentPageNodes)

		// 检查分页信息
		if response.Body.Page != nil {
			if response.Body.Page.TotalCount != nil {
				log.Printf("   📊 总记录数: %d", *response.Body.Page.TotalCount)
			}
		}

		// 判断是否还有下一页
		if currentPageNodes < pageSize {
			log.Printf("   ✅ 已到最后一页")
			break
		}

		pageNumber++
	}

	log.Printf("🎉 [getAllClusterNodesWithPagination] 分页查询完成")
	log.Printf("📊 [getAllClusterNodesWithPagination] 共获取 %d 个实例", len(allInstanceIds))
	log.Printf("📋 [getAllClusterNodesWithPagination] 实例列表: %v", allInstanceIds)
	return allInstanceIds, nil
}

// filterInstancesByZone 根据可用区过滤实例
func (chaos *AliyunRegionChaos) filterInstancesByZone(instanceIds []string, zoneId string) ([]string, error) {
	if len(instanceIds) == 0 {
		return []string{}, nil
	}

	regionId := os.Getenv("ALIBABA_CLOUD_REGION_ID")
	if regionId == "" {
		regionId = "cn-heyuan" // 默认值
	}

	log.Printf("🔍 [filterInstancesByZone] 开始按可用区过滤实例:")
	log.Printf("   - 输入实例数: %d", len(instanceIds))
	log.Printf("   - 目标可用区: %s", zoneId)
	log.Printf("   - 输入实例列表: %v", instanceIds)

	// 调用 DescribeInstanceStatus
	request := &ecs.DescribeInstanceStatusRequest{
		InstanceId: tea.StringSlice(instanceIds),
		ZoneId:     tea.String(zoneId),
		RegionId:   tea.String(regionId),
	}

	response, err := EcsClient.DescribeInstanceStatus(request)
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

	log.Printf("✅ [filterInstancesByZone] 可用区过滤完成:")
	log.Printf("   - 输出实例数: %d", len(filteredIds))
	log.Printf("   - 过滤掉的实例数: %d", len(instanceIds)-len(filteredIds))
	log.Printf("   - 输出实例列表: %v", filteredIds)

	return filteredIds, nil
}

// stopInstances 批量停止实例
func (chaos *AliyunRegionChaos) stopInstances(instanceIds []string) error {
	log.Printf("🔥 [stopInstances] 开始批量停止实例")
	
	if len(instanceIds) == 0 {
		log.Printf("⚠️  [stopInstances] 实例列表为空，跳过操作")
		return nil
	}

	regionId := os.Getenv("ALIBABA_CLOUD_REGION_ID")
	if regionId == "" {
		regionId = "cn-heyuan" // 默认值
	}

	log.Printf("📋 [stopInstances] 请求参数:")
	log.Printf("   - 区域: %s", regionId)
	log.Printf("   - 实例数量: %d", len(instanceIds))
	log.Printf("   - 实例列表: %v", instanceIds)
	log.Printf("   - 强制停止: true")

	request := &ecs.StopInstancesRequest{
		InstanceId: tea.StringSlice(instanceIds),
		ForceStop:  tea.Bool(true),
		RegionId:   tea.String(regionId),
	}

	log.Printf("🚀 [stopInstances] 发送停止实例API请求...")
	response, err := EcsClient.StopInstances(request)
	if err != nil {
		log.Printf("❌ [stopInstances] API调用失败: %v", err)
		return fmt.Errorf("调用 StopInstances 失败: %v", err)
	}

	if response.Body.RequestId != nil {
		log.Printf("✅ [stopInstances] API调用成功")
		log.Printf("📝 [stopInstances] 请求ID: %s", *response.Body.RequestId)
		log.Printf("⏳ [stopInstances] 实例停止中，预计1分钟内完成")
	}

	return nil
}

// rebootInstances 批量重启实例（智能选择启动或重启）
func (chaos *AliyunRegionChaos) rebootInstances(instanceIds []string) error {
	log.Printf("🔄 [rebootInstances] 开始批量重启实例")
	
	if len(instanceIds) == 0 {
		log.Printf("⚠️  [rebootInstances] 实例列表为空，跳过操作")
		return nil
	}

	regionId := os.Getenv("ALIBABA_CLOUD_REGION_ID")
	if regionId == "" {
		regionId = "cn-heyuan" // 默认值
	}

	log.Printf("📋 [rebootInstances] 重启参数:")
	log.Printf("   - 区域: %s", regionId)
	log.Printf("   - 实例数量: %d", len(instanceIds))
	log.Printf("   - 实例列表: %v", instanceIds)

	// 首先检查实例状态，决定使用 StartInstances 还是 RebootInstances
	log.Printf("🔍 [rebootInstances] 检查实例状态...")
	statusRequest := &ecs.DescribeInstanceStatusRequest{
		InstanceId: tea.StringSlice(instanceIds),
		RegionId:   tea.String(regionId),
	}

	statusResponse, err := EcsClient.DescribeInstanceStatus(statusRequest)
	if err != nil {
		log.Printf("❌ [rebootInstances] 检查实例状态失败: %v", err)
		return fmt.Errorf("检查实例状态失败: %v", err)
	}

	var stoppedInstances []string  // 需要启动的实例
	var runningInstances []string  // 需要重启的实例

	log.Printf("📊 [rebootInstances] 实例状态分析:")
	if statusResponse.Body.InstanceStatuses != nil && statusResponse.Body.InstanceStatuses.InstanceStatus != nil {
		for _, status := range statusResponse.Body.InstanceStatuses.InstanceStatus {
			if status.InstanceId != nil && status.Status != nil {
				instanceId := *status.InstanceId
				instanceStatus := *status.Status
				
				log.Printf("   - 实例 %s: %s", instanceId, instanceStatus)
				
				switch instanceStatus {
				case "Stopped", "Shutted":
					stoppedInstances = append(stoppedInstances, instanceId)
				case "Running":
					runningInstances = append(runningInstances, instanceId)
				default:
					log.Printf("⚠️  [rebootInstances] 实例 %s 状态异常: %s", instanceId, instanceStatus)
				}
			}
		}
	}

	log.Printf("🎯 [rebootInstances] 操作策略:")
	log.Printf("   - 需要启动的实例: %d 个 %v", len(stoppedInstances), stoppedInstances)
	log.Printf("   - 需要重启的实例: %d 个 %v", len(runningInstances), runningInstances)

	// 启动已停止的实例
	if len(stoppedInstances) > 0 {
		log.Printf("🚀 [rebootInstances] 启动已停止的实例...")
		log.Printf("📋 [rebootInstances] 启动 %d 个实例: %v", len(stoppedInstances), stoppedInstances)
		
		startRequest := &ecs.StartInstancesRequest{
			InstanceId: tea.StringSlice(stoppedInstances),
			RegionId:   tea.String(regionId),
		}

		startResponse, err := EcsClient.StartInstances(startRequest)
		if err != nil {
			log.Printf("❌ [rebootInstances] 启动实例失败: %v", err)
			return fmt.Errorf("调用 StartInstances 失败: %v", err)
		}

		if startResponse.Body.RequestId != nil {
			log.Printf("✅ [rebootInstances] 启动请求成功")
			log.Printf("📝 [rebootInstances] 启动请求ID: %s", *startResponse.Body.RequestId)
		}
	}

	// 重启正在运行的实例
	if len(runningInstances) > 0 {
		log.Printf("🔄 [rebootInstances] 重启运行中的实例...")
		log.Printf("📋 [rebootInstances] 重启 %d 个实例: %v", len(runningInstances), runningInstances)
		
		rebootRequest := &ecs.RebootInstancesRequest{
			InstanceId: tea.StringSlice(runningInstances),
			RegionId:   tea.String(regionId),
		}

		rebootResponse, err := EcsClient.RebootInstances(rebootRequest)
		if err != nil {
			log.Printf("❌ [rebootInstances] 重启实例失败: %v", err)
			return fmt.Errorf("调用 RebootInstances 失败: %v", err)
		}

		if rebootResponse.Body.RequestId != nil {
			log.Printf("✅ [rebootInstances] 重启请求成功")
			log.Printf("📝 [rebootInstances] 重启请求ID: %s", *rebootResponse.Body.RequestId)
		}
	}

	totalProcessed := len(stoppedInstances) + len(runningInstances)
	log.Printf("🎉 [rebootInstances] 批量重启操作完成")
	log.Printf("📊 [rebootInstances] 处理统计: 总计 %d 个实例 (启动: %d, 重启: %d)", 
		totalProcessed, len(stoppedInstances), len(runningInstances))
	log.Printf("⏳ [rebootInstances] 预计实例恢复时间: 1-2分钟")

	return nil
}

// ============================================================================
// 节点故障注入实现 - 从当前AZ任意挑选一台机器进行操作
// ============================================================================

type AliyunNodeChaos struct {
	lastFlagdValue      int64
	flagdValue          int64
	instanceIds         []string
	stoppedInstances    []string
	lastOperationTime   time.Time
	minOperationInterval time.Duration
}

func NewAliyunNodeChaos() *AliyunNodeChaos {
	// 🔧 初始化阿里云客户端（一次性，在创建实例时初始化，避免每次ExecuteChaos都调用）
	InitAliyunClients()
	
	// 从环境变量获取最小操作间隔（单位：分钟），默认3分钟
	minIntervalMinutes := 3 // 默认值
	if intervalEnv := os.Getenv("MIN_OPERATION_INTERVAL_MINUTES"); intervalEnv != "" {
		if parsed, err := strconv.Atoi(intervalEnv); err == nil && parsed > 0 {
			minIntervalMinutes = parsed
			log.Printf("📝 [NewAliyunNodeChaos] 从环境变量读取最小操作间隔: %d 分钟", minIntervalMinutes)
		} else {
			log.Printf("⚠️  [NewAliyunNodeChaos] 环境变量 MIN_OPERATION_INTERVAL_MINUTES 无效 (%s)，使用默认值: %d 分钟", 
				intervalEnv, minIntervalMinutes)
		}
	} else {
		log.Printf("ℹ️  [NewAliyunNodeChaos] 未设置 MIN_OPERATION_INTERVAL_MINUTES，使用默认值: %d 分钟", minIntervalMinutes)
	}

	minInterval := time.Duration(minIntervalMinutes) * time.Minute
	
	return &AliyunNodeChaos{
		lastFlagdValue:       0,
		flagdValue:           0,
		instanceIds:          []string{},
		stoppedInstances:     []string{},
		lastOperationTime:    time.Time{}, // 零值表示没有操作过
		minOperationInterval: minInterval, // 从环境变量配置的最小操作间隔
	}
}

func (chaos *AliyunNodeChaos) ExecuteChaos(flagdValue int64) {
	// 客户端已在 NewAliyunNodeChaos() 中初始化，无需重复调用
	currentTime := time.Now()
	
	log.Printf("🔍 [NodeChaos ExecuteChaos] 收到请求 - 当前时间: %s", currentTime.Format("2006-01-02 15:04:05"))
	log.Printf("📊 [NodeChaos ExecuteChaos] 当前状态 - flagd值: %d → %d, 已停止实例数: %d", 
		chaos.lastFlagdValue, flagdValue, len(chaos.stoppedInstances))
	
	chaos.flagdValue = flagdValue
	if chaos.flagdValue == chaos.lastFlagdValue {
		log.Printf("⏭️  [NodeChaos ExecuteChaos] flagd值未变化 (%d)，跳过操作", flagdValue)
		return
	}

	// 检查时间间隔保护：如果不是第一次操作且距离上次操作不足3分钟，则忽略此次请求
	if !chaos.lastOperationTime.IsZero() {
		timeSinceLastOp := currentTime.Sub(chaos.lastOperationTime)
		nextAllowedTime := chaos.lastOperationTime.Add(chaos.minOperationInterval)
		
		log.Printf("🕐 [NodeChaos 时间检查] 上次操作时间: %s", chaos.lastOperationTime.Format("2006-01-02 15:04:05"))
		log.Printf("⏰ [NodeChaos 时间检查] 距离上次操作: %.1f 秒, 要求间隔: %.0f 秒 (%.1f 分钟)", 
			timeSinceLastOp.Seconds(), chaos.minOperationInterval.Seconds(), chaos.minOperationInterval.Minutes())
		log.Printf("🚀 [NodeChaos 时间检查] 下次可操作时间: %s", nextAllowedTime.Format("2006-01-02 15:04:05"))
		
		if timeSinceLastOp < chaos.minOperationInterval {
			remainingTime := chaos.minOperationInterval - timeSinceLastOp
			log.Printf("❌ [NodeChaos 操作拒绝] 时间间隔不足，拒绝执行")
			log.Printf("⚠️  [NodeChaos 操作拒绝] flagd变更: %d → %d", chaos.lastFlagdValue, chaos.flagdValue)
			log.Printf("⌛ [NodeChaos 操作拒绝] 还需等待: %.1f 秒 (到 %s)", 
				remainingTime.Seconds(), nextAllowedTime.Format("15:04:05"))
			log.Printf("🔒 [NodeChaos 操作拒绝] 实例状态保持不变，当前已停止: %v", chaos.stoppedInstances)
			return
		} else {
			log.Printf("✅ [NodeChaos 时间检查] 间隔足够，允许执行操作")
		}
	} else {
		log.Printf("🎯 [NodeChaos 时间检查] 首次操作，无需检查时间间隔")
	}

	log.Printf("🚀 [NodeChaos 开始执行] AliyunNodeChaos 执行操作: flagd值 %d → %d", chaos.lastFlagdValue, chaos.flagdValue)

	var operationSuccess bool
	if chaos.flagdValue > 0 {
		log.Printf("🔥 [NodeChaos 故障注入] 开始故障注入，持续时间: %d 秒", chaos.flagdValue)
		operationSuccess = chaos.startNodeChaos()
	} else {
		log.Printf("🔄 [NodeChaos 故障恢复] 停止故障注入，恢复实例")
		operationSuccess = chaos.stopNodeChaos()
	}

	if operationSuccess {
		// 更新最后操作时间和flagd值
		chaos.lastOperationTime = currentTime
		chaos.lastFlagdValue = chaos.flagdValue
		nextAllowedTime := currentTime.Add(chaos.minOperationInterval)
		
		log.Printf("✅ [NodeChaos 操作完成] 成功执行操作，状态已更新")
		log.Printf("📝 [NodeChaos 操作完成] 新的lastFlagdValue: %d", chaos.lastFlagdValue)
		log.Printf("🕐 [NodeChaos 操作完成] 操作时间已记录: %s", currentTime.Format("2006-01-02 15:04:05"))
		log.Printf("⏰ [NodeChaos 操作完成] 下次可操作时间: %s", nextAllowedTime.Format("2006-01-02 15:04:05"))
		log.Printf("📊 [NodeChaos 操作完成] 当前已停止实例: %v", chaos.stoppedInstances)
	} else {
		log.Printf("❌ [NodeChaos 操作失败] 操作执行失败，状态未更新")
		log.Printf("🔍 [NodeChaos 操作失败] 保持原有状态: lastFlagdValue=%d, 已停止实例=%v", 
			chaos.lastFlagdValue, chaos.stoppedInstances)
	}
}

// startNodeChaos 开始节点故障注入 - 从当前AZ任意挑选一台机器
func (chaos *AliyunNodeChaos) startNodeChaos() bool {
	operationTime := time.Now()
	log.Printf("🔥 [startNodeChaos] 开始阿里云节点故障注入")
	log.Printf("🕐 [startNodeChaos] 操作开始时间: %s", operationTime.Format("2006-01-02 15:04:05"))
	log.Printf("📋 [startNodeChaos] 当前实例状态 - instanceIds: %v, stoppedInstances: %v", 
		chaos.instanceIds, chaos.stoppedInstances)

	// 获取当前AZ的实例列表
	log.Printf("🔍 [startNodeChaos] 正在获取当前AZ的实例列表...")
	instanceIds, err := chaos.getCurrentZoneInstances()
	if err != nil {
		log.Printf("❌ [startNodeChaos] 获取当前AZ实例失败: %v", err)
		return false
	}

	if len(instanceIds) == 0 {
		log.Printf("⚠️  [startNodeChaos] 未找到当前AZ的可操作实例")
		log.Printf("🔍 [startNodeChaos] 请检查环境变量: CLUSTER_ID, ZONE_ID")
		return false
	}

	// 从实例列表中随机选择一台机器
	selectedInstance := chaos.selectRandomInstance(instanceIds)
	log.Printf("✅ [startNodeChaos] 从 %d 个实例中随机选择: %s", len(instanceIds), selectedInstance)
	log.Printf("📋 [startNodeChaos] 可选实例列表: %v", instanceIds)

	selectedInstances := []string{selectedInstance}
	chaos.instanceIds = selectedInstances

	// 停止选中的实例
	log.Printf("🔄 [startNodeChaos] 开始停止选中的实例...")
	err = chaos.stopInstances(selectedInstances)
	if err != nil {
		log.Printf("❌ [startNodeChaos] 停止实例失败: %v", err)
		log.Printf("🔍 [startNodeChaos] 实例状态未改变，保持原状态")
		return false
	}

	chaos.stoppedInstances = selectedInstances
	estimatedCompleteTime := operationTime.Add(1 * time.Minute)
	log.Printf("✅ [startNodeChaos] 成功发送停止请求，影响实例: %s", selectedInstance)
	log.Printf("⏱️  [startNodeChaos] 预计停机完成时间: %s (约1分钟后)", estimatedCompleteTime.Format("15:04:05"))
	log.Printf("📊 [startNodeChaos] 已停止实例列表: %v", chaos.stoppedInstances)
	
	return true
}

// stopNodeChaos 停止节点故障注入
func (chaos *AliyunNodeChaos) stopNodeChaos() bool {
	operationTime := time.Now()
	log.Printf("🔄 [stopNodeChaos] 开始停止节点故障注入")
	log.Printf("🕐 [stopNodeChaos] 操作开始时间: %s", operationTime.Format("2006-01-02 15:04:05"))
	log.Printf("📋 [stopNodeChaos] 当前已停止实例: %v", chaos.stoppedInstances)
	
	if len(chaos.stoppedInstances) == 0 {
		log.Printf("ℹ️  [stopNodeChaos] 没有需要重启的实例")
		log.Printf("✅ [stopNodeChaos] 操作完成（无需重启）")
		return true
	}

	log.Printf("🚀 [stopNodeChaos] 准备重启 %d 个实例: %v", 
		len(chaos.stoppedInstances), chaos.stoppedInstances)

	// 重启实例
	log.Printf("🔄 [stopNodeChaos] 开始重启实例...")
	err := chaos.rebootInstances(chaos.stoppedInstances)
	if err != nil {
		log.Printf("❌ [stopNodeChaos] 重启实例失败: %v", err)
		log.Printf("🔍 [stopNodeChaos] 实例状态未改变，仍保持停止状态: %v", chaos.stoppedInstances)
		return false
	}

	estimatedCompleteTime := operationTime.Add(1 * time.Minute)
	log.Printf("✅ [stopNodeChaos] 成功发送重启请求")
	log.Printf("⏱️  [stopNodeChaos] 预计重启完成时间: %s (约1分钟后)", estimatedCompleteTime.Format("15:04:05"))
	log.Printf("📊 [stopNodeChaos] 清空已停止实例列表")
	
	chaos.stoppedInstances = []string{}
	return true
}

// getCurrentZoneInstances 获取当前可用区的实例列表
func (chaos *AliyunNodeChaos) getCurrentZoneInstances() ([]string, error) {
	clusterId := os.Getenv("CLUSTER_ID")
	zoneId := os.Getenv("ZONE_ID")

	if clusterId == "" {
		return nil, fmt.Errorf("CLUSTER_ID 环境变量未设置")
	}

	if zoneId == "" {
		return nil, fmt.Errorf("ZONE_ID 环境变量未设置")
	}

	// 1. 调用 DescribeClusterNodes 获取集群节点（支持分页）
	log.Printf("🔍 [getCurrentZoneInstances] 开始获取当前可用区实例列表（支持分页查询）")
	allInstanceIds, err := chaos.getAllClusterNodesWithPagination(clusterId, "")
	if err != nil {
		return nil, fmt.Errorf("获取集群节点失败: %v", err)
	}

	log.Printf("✅ [getCurrentZoneInstances] 获取到 %d 个集群实例", len(allInstanceIds))
	log.Printf("📋 [getCurrentZoneInstances] 原始实例列表: %v", allInstanceIds)

	if len(allInstanceIds) == 0 {
		log.Printf("⚠️  [getCurrentZoneInstances] 未找到任何实例")
		return []string{}, nil
	}

	// 2. 使用 DescribeInstanceStatus 过滤特定区域的实例
	log.Printf("🔍 [getCurrentZoneInstances] 按可用区 %s 过滤实例", zoneId)
	filteredIds, err := chaos.filterInstancesByZone(allInstanceIds, zoneId)
	if err != nil {
		return nil, err
	}
	log.Printf("✅ [getCurrentZoneInstances] 可用区过滤后剩余 %d 个实例", len(filteredIds))
	log.Printf("📋 [getCurrentZoneInstances] 可用区过滤后实例列表: %v", filteredIds)
	return filteredIds, nil
}

// getAllClusterNodesWithPagination 获取集群所有节点，支持分页（NodeChaos版本）
func (chaos *AliyunNodeChaos) getAllClusterNodesWithPagination(clusterId, nodePoolId string) ([]string, error) {
	var allInstanceIds []string
	pageNumber := 1
	pageSize := 10 // 设置每页大小

	log.Printf("🔍 [NodeChaos getAllClusterNodesWithPagination] 开始分页查询集群节点")
	if nodePoolId != "" {
		log.Printf("🎯 [NodeChaos getAllClusterNodesWithPagination] 限定节点池: %s", nodePoolId)
	} else {
		log.Printf("🌍 [NodeChaos getAllClusterNodesWithPagination] 查询整个集群节点")
	}

	for {
		request := &cs.DescribeClusterNodesRequest{
			PageNumber: tea.String(fmt.Sprintf("%d", pageNumber)),
			PageSize:   tea.String(fmt.Sprintf("%d", pageSize)),
		}
		
		if nodePoolId != "" {
			request.NodepoolId = tea.String(nodePoolId)
		}

		log.Printf("📄 [NodeChaos getAllClusterNodesWithPagination] 查询第 %d 页（每页 %d 条）", pageNumber, pageSize)
		
		response, err := CsClient.DescribeClusterNodes(tea.String(clusterId), request)
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

		log.Printf("   ✅ 第 %d 页返回 %d 个节点", pageNumber, currentPageNodes)

		// 检查分页信息
		if response.Body.Page != nil {
			if response.Body.Page.TotalCount != nil {
				log.Printf("   📊 总记录数: %d", *response.Body.Page.TotalCount)
			}
		}

		// 判断是否还有下一页
		if currentPageNodes < pageSize {
			log.Printf("   ✅ 已到最后一页")
			break
		}

		pageNumber++
	}

	log.Printf("🎉 [NodeChaos getAllClusterNodesWithPagination] 分页查询完成")
	log.Printf("📊 [NodeChaos getAllClusterNodesWithPagination] 共获取 %d 个实例", len(allInstanceIds))
	log.Printf("📋 [NodeChaos getAllClusterNodesWithPagination] 实例列表: %v", allInstanceIds)
	return allInstanceIds, nil
}

// selectRandomInstance 从实例列表中随机选择一台机器
func (chaos *AliyunNodeChaos) selectRandomInstance(instanceIds []string) string {
	if len(instanceIds) == 0 {
		return ""
	}
	
	if len(instanceIds) == 1 {
		log.Printf("🎯 [selectRandomInstance] 只有一台实例可选: %s", instanceIds[0])
		return instanceIds[0]
	}

	// 使用当前时间作为随机种子
	currentTime := time.Now()
	index := int(currentTime.UnixNano()) % len(instanceIds)
	selectedInstance := instanceIds[index]
	
	log.Printf("🎲 [selectRandomInstance] 随机选择逻辑:")
	log.Printf("   - 可选实例数量: %d", len(instanceIds))
	log.Printf("   - 时间种子: %d", currentTime.UnixNano())
	log.Printf("   - 计算索引: %d", index)
	log.Printf("   - 选中实例: %s", selectedInstance)
	
	return selectedInstance
}

// 复用 RegionChaos 的方法
func (chaos *AliyunNodeChaos) filterInstancesByZone(instanceIds []string, zoneId string) ([]string, error) {
	if len(instanceIds) == 0 {
		return []string{}, nil
	}

	regionId := os.Getenv("ALIBABA_CLOUD_REGION_ID")
	if regionId == "" {
		regionId = "cn-heyuan" // 默认值
	}

	log.Printf("🔍 [NodeChaos filterInstancesByZone] 开始按可用区过滤实例:")
	log.Printf("   - 输入实例数: %d", len(instanceIds))
	log.Printf("   - 目标可用区: %s", zoneId)
	log.Printf("   - 输入实例列表: %v", instanceIds)

	// 调用 DescribeInstanceStatus
	request := &ecs.DescribeInstanceStatusRequest{
		InstanceId: tea.StringSlice(instanceIds),
		ZoneId:     tea.String(zoneId),
		RegionId:   tea.String(regionId),
	}

	response, err := EcsClient.DescribeInstanceStatus(request)
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

	log.Printf("✅ [NodeChaos filterInstancesByZone] 可用区过滤完成:")
	log.Printf("   - 输出实例数: %d", len(filteredIds))
	log.Printf("   - 过滤掉的实例数: %d", len(instanceIds)-len(filteredIds))
	log.Printf("   - 输出实例列表: %v", filteredIds)

	return filteredIds, nil
}

// 复用 RegionChaos 的方法
func (chaos *AliyunNodeChaos) stopInstances(instanceIds []string) error {
	log.Printf("🔥 [NodeChaos stopInstances] 开始批量停止实例")
	
	if len(instanceIds) == 0 {
		log.Printf("⚠️  [NodeChaos stopInstances] 实例列表为空，跳过操作")
		return nil
	}

	regionId := os.Getenv("ALIBABA_CLOUD_REGION_ID")
	if regionId == "" {
		regionId = "cn-heyuan" // 默认值
	}

	log.Printf("📋 [NodeChaos stopInstances] 请求参数:")
	log.Printf("   - 区域: %s", regionId)
	log.Printf("   - 实例数量: %d", len(instanceIds))
	log.Printf("   - 实例列表: %v", instanceIds)
	log.Printf("   - 强制停止: true")

	request := &ecs.StopInstancesRequest{
		InstanceId: tea.StringSlice(instanceIds),
		ForceStop:  tea.Bool(true),
		RegionId:   tea.String(regionId),
	}

	log.Printf("🚀 [NodeChaos stopInstances] 发送停止实例API请求...")
	response, err := EcsClient.StopInstances(request)
	if err != nil {
		log.Printf("❌ [NodeChaos stopInstances] API调用失败: %v", err)
		return fmt.Errorf("调用 StopInstances 失败: %v", err)
	}

	if response.Body.RequestId != nil {
		log.Printf("✅ [NodeChaos stopInstances] API调用成功")
		log.Printf("📝 [NodeChaos stopInstances] 请求ID: %s", *response.Body.RequestId)
		log.Printf("⏳ [NodeChaos stopInstances] 实例停止中，预计1分钟内完成")
	}

	return nil
}

// 复用 RegionChaos 的方法
func (chaos *AliyunNodeChaos) rebootInstances(instanceIds []string) error {
	log.Printf("🔄 [NodeChaos rebootInstances] 开始批量重启实例")
	
	if len(instanceIds) == 0 {
		log.Printf("⚠️  [NodeChaos rebootInstances] 实例列表为空，跳过操作")
		return nil
	}

	regionId := os.Getenv("ALIBABA_CLOUD_REGION_ID")
	if regionId == "" {
		regionId = "cn-heyuan" // 默认值
	}

	log.Printf("📋 [NodeChaos rebootInstances] 重启参数:")
	log.Printf("   - 区域: %s", regionId)
	log.Printf("   - 实例数量: %d", len(instanceIds))
	log.Printf("   - 实例列表: %v", instanceIds)

	// 首先检查实例状态，决定使用 StartInstances 还是 RebootInstances
	log.Printf("🔍 [NodeChaos rebootInstances] 检查实例状态...")
	statusRequest := &ecs.DescribeInstanceStatusRequest{
		InstanceId: tea.StringSlice(instanceIds),
		RegionId:   tea.String(regionId),
	}

	statusResponse, err := EcsClient.DescribeInstanceStatus(statusRequest)
	if err != nil {
		log.Printf("❌ [NodeChaos rebootInstances] 检查实例状态失败: %v", err)
		return fmt.Errorf("检查实例状态失败: %v", err)
	}

	var stoppedInstances []string  // 需要启动的实例
	var runningInstances []string  // 需要重启的实例

	log.Printf("📊 [NodeChaos rebootInstances] 实例状态分析:")
	if statusResponse.Body.InstanceStatuses != nil && statusResponse.Body.InstanceStatuses.InstanceStatus != nil {
		for _, status := range statusResponse.Body.InstanceStatuses.InstanceStatus {
			if status.InstanceId != nil && status.Status != nil {
				instanceId := *status.InstanceId
				instanceStatus := *status.Status
				
				log.Printf("   - 实例 %s: %s", instanceId, instanceStatus)
				
				switch instanceStatus {
				case "Stopped", "Shutted":
					stoppedInstances = append(stoppedInstances, instanceId)
				case "Running":
					runningInstances = append(runningInstances, instanceId)
				default:
					log.Printf("⚠️  [NodeChaos rebootInstances] 实例 %s 状态异常: %s", instanceId, instanceStatus)
				}
			}
		}
	}

	log.Printf("🎯 [NodeChaos rebootInstances] 操作策略:")
	log.Printf("   - 需要启动的实例: %d 个 %v", len(stoppedInstances), stoppedInstances)
	log.Printf("   - 需要重启的实例: %d 个 %v", len(runningInstances), runningInstances)

	// 启动已停止的实例
	if len(stoppedInstances) > 0 {
		log.Printf("🚀 [NodeChaos rebootInstances] 启动已停止的实例...")
		log.Printf("📋 [NodeChaos rebootInstances] 启动 %d 个实例: %v", len(stoppedInstances), stoppedInstances)
		
		startRequest := &ecs.StartInstancesRequest{
			InstanceId: tea.StringSlice(stoppedInstances),
			RegionId:   tea.String(regionId),
		}

		startResponse, err := EcsClient.StartInstances(startRequest)
		if err != nil {
			log.Printf("❌ [NodeChaos rebootInstances] 启动实例失败: %v", err)
			return fmt.Errorf("调用 StartInstances 失败: %v", err)
		}

		if startResponse.Body.RequestId != nil {
			log.Printf("✅ [NodeChaos rebootInstances] 启动请求成功")
			log.Printf("📝 [NodeChaos rebootInstances] 启动请求ID: %s", *startResponse.Body.RequestId)
		}
	}

	// 重启正在运行的实例
	if len(runningInstances) > 0 {
		log.Printf("🔄 [NodeChaos rebootInstances] 重启运行中的实例...")
		log.Printf("📋 [NodeChaos rebootInstances] 重启 %d 个实例: %v", len(runningInstances), runningInstances)
		
		rebootRequest := &ecs.RebootInstancesRequest{
			InstanceId: tea.StringSlice(runningInstances),
			RegionId:   tea.String(regionId),
		}

		rebootResponse, err := EcsClient.RebootInstances(rebootRequest)
		if err != nil {
			log.Printf("❌ [NodeChaos rebootInstances] 重启实例失败: %v", err)
			return fmt.Errorf("调用 RebootInstances 失败: %v", err)
		}

		if rebootResponse.Body.RequestId != nil {
			log.Printf("✅ [NodeChaos rebootInstances] 重启请求成功")
			log.Printf("📝 [NodeChaos rebootInstances] 重启请求ID: %s", *rebootResponse.Body.RequestId)
		}
	}

	totalProcessed := len(stoppedInstances) + len(runningInstances)
	log.Printf("🎉 [NodeChaos rebootInstances] 批量重启操作完成")
	log.Printf("📊 [NodeChaos rebootInstances] 处理统计: 总计 %d 个实例 (启动: %d, 重启: %d)", 
		totalProcessed, len(stoppedInstances), len(runningInstances))
	log.Printf("⏳ [NodeChaos rebootInstances] 预计实例恢复时间: 1-2分钟")

	return nil
}
 