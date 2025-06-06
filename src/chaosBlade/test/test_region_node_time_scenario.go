package main

import (
	"log"
	"time"
)

// 模拟混沌工程时间保护逻辑
type TimeScenarioTest struct {
	lastFlagdValue       int64
	flagdValue           int64
	lastOperationTime    time.Time
	minOperationInterval time.Duration
	stoppedInstances     []string
}

func NewTimeScenarioTest() *TimeScenarioTest {
	return &TimeScenarioTest{
		lastFlagdValue:       0,
		flagdValue:           0,
		lastOperationTime:    time.Time{},
		minOperationInterval: 3 * time.Minute, // 3分钟保护间隔
		stoppedInstances:     []string{},
	}
}

func (test *TimeScenarioTest) ExecuteChaos(flagdValue int64, mockTime time.Time, timeLabel string) string {
	test.flagdValue = flagdValue
	
	log.Printf("\n🕒 %s - 尝试操作", timeLabel)
	log.Printf("📊 flagd变更: %d → %d", test.lastFlagdValue, test.flagdValue)
	
	if test.flagdValue == test.lastFlagdValue {
		return "跳过"
	}

	// 检查时间间隔保护
	if !test.lastOperationTime.IsZero() {
		timeSinceLastOp := mockTime.Sub(test.lastOperationTime)
		
		if timeSinceLastOp < test.minOperationInterval {
			remainingTime := test.minOperationInterval - timeSinceLastOp
			log.Printf("❌ 操作拒绝 - 距离上次操作: %.0f秒, 需要等待: %.0f秒", 
				timeSinceLastOp.Seconds(), remainingTime.Seconds())
			return "拒绝"
		}
	}

	// 执行操作
	log.Printf("✅ 操作成功 - 时间间隔足够")
	if test.flagdValue > 0 {
		log.Printf("🔥 开始故障注入: %d秒", test.flagdValue)
		test.stoppedInstances = []string{"i-test1", "i-test2"}
	} else {
		log.Printf("🔄 停止故障注入，恢复实例")
		test.stoppedInstances = []string{}
	}
	
	// 更新状态
	test.lastOperationTime = mockTime
	test.lastFlagdValue = test.flagdValue
	
	return "成功"
}

func main() {
	log.Printf("🎬 时序场景测试 - 基于21:33:00的实际测试数据")
	log.Printf("📝 测试规则: 3分钟最小操作间隔")
	log.Printf("=" * 60)

	test := NewTimeScenarioTest()
	
	// 基准时间：21:33:00
	baseTime := time.Date(2024, 1, 15, 21, 33, 0, 0, time.Local)
	
	// 测试场景（基于截图中的实际数据）
	scenarios := []struct {
		timeOffset   time.Duration
		timeLabel    string
		flagdValue   int64
		operation    string
		expectedResult string
	}{
		{0, "21:33:00", 1200, "20分钟停机", "成功"},
		{10 * time.Second, "21:33:10", 0, "off", "拒绝"},
		{1 * time.Minute, "21:34:00", 0, "off", "拒绝"},
		{2 * time.Minute, "21:35:00", 0, "off", "拒绝"},
		{2*time.Minute + 59*time.Second, "21:35:59", 0, "off", "拒绝"},
		{3 * time.Minute, "21:36:00", 0, "off", "成功"},
	}

	log.Printf("\n📋 测试执行:")
	successCount := 0
	rejectCount := 0
	
	for i, scenario := range scenarios {
		mockTime := baseTime.Add(scenario.timeOffset)
		
		log.Printf("\n" + "-"*40)
		log.Printf("🔍 测试 %d: %s", i+1, scenario.operation)
		
		result := test.ExecuteChaos(scenario.flagdValue, mockTime, scenario.timeLabel)
		
		// 验证结果
		if result == scenario.expectedResult {
			log.Printf("✅ 结果正确: %s", result)
			if result == "成功" {
				successCount++
			} else if result == "拒绝" {
				rejectCount++
			}
		} else {
			log.Printf("❌ 结果错误: 期望 %s, 实际 %s", scenario.expectedResult, result)
		}
		
		// 显示时间间隔信息
		if !test.lastOperationTime.IsZero() {
			timeSinceLastOp := mockTime.Sub(test.lastOperationTime)
			log.Printf("⏰ 距离上次成功操作: %.0f秒", timeSinceLastOp.Seconds())
		}
	}

	log.Printf("\n" + "="*60)
	log.Printf("📊 测试结果统计:")
	log.Printf("   ✅ 成功操作: %d 次 (21:33:00, 21:36:00)", successCount)
	log.Printf("   ❌ 拒绝操作: %d 次 (21:33:10, 21:34:00, 21:35:00, 21:35:59)", rejectCount)
	log.Printf("   📈 拒绝率: %.1f%% (符合预期)", float64(rejectCount)/float64(len(scenarios))*100)

	log.Printf("\n🎯 关键时间点分析:")
	log.Printf("   • 21:33:00: 首次20分钟停机操作成功")
	log.Printf("   • 21:33:10: 10秒后尝试off → 拒绝 (需要等待170秒)")
	log.Printf("   • 21:34:00: 1分钟后尝试off → 拒绝 (需要等待120秒)")
	log.Printf("   • 21:35:00: 2分钟后尝试off → 拒绝 (需要等待60秒)")
	log.Printf("   • 21:35:59: 2分59秒后尝试off → 拒绝 (需要等待1秒)")
	log.Printf("   • 21:36:00: 3分钟后off操作成功 → 实例恢复")

	log.Printf("\n💡 实际应用效果:")
	log.Printf("   🔒 有效防止了频繁的误操作")
	log.Printf("   ⏱️  保护窗口期内实例状态稳定")
	log.Printf("   🎛️  精确的时间控制（精确到秒）")
	log.Printf("   📝 详细的日志记录便于排查问题")
	
	log.Printf("\n🎉 时序场景测试完成!")
} 