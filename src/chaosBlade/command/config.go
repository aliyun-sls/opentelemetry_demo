package command

import (
	cs "github.com/alibabacloud-go/cs-20151215/v5/client"
	ecs "github.com/alibabacloud-go/ecs-20140526/v4/client"
	"github.com/open-feature/go-sdk/openfeature"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/client-go/dynamic"
	"sync"
)

type RegionLoss struct {
	regionlossLastConfig int64
	flagdValue           int64
}

type NodeLoss struct {
	nodeLossLastConfig int64
	flagdValue         int64
}

type NodeCpu struct {
	nodeCpuLastConfig int64
	flagdValue        int64
}

type NodeMem struct {
	nodeMemLastConfig int64
	flagdValue        int64
}

type NodeDisk struct {
	nodeDiskLastConfig int64
	flagdValue         int64
}

type PodNetDelay struct {
	podNetDelayLastConfig string
	flagdValue            string
	labels                string
}

type PodCpu struct {
	podCpuLastConfig int64
	flagdValue       int64
}

type PodMem struct {
	podMemLastConfig int64
	flagdValue       int64
}

type Netloss struct {
	Name    string
	Scope   string
	Labels  string
	Port    string
	Percent string
	Timeout string
}

type Netdelay struct {
	Name    string
	Scope   string
	Labels  string
	Port    string
	Time    string
	Offset  string
	Timeout string
}

type CpuAndMem struct {
	Name    string
	Scope   string
	Labels  string
	Percent string
	Timeout string
}

var Dynamic *dynamic.DynamicClient
var Gvr schema.GroupVersionResource
var (
	FlagClient *openfeature.Client
	Once       sync.Once

	// 阿里云客户端 - 只初始化一次
	EcsClient  *ecs.Client
	CsClient   *cs.Client
	AliyunOnce sync.Once

	// 阿里云混沌工程实例 - 全局单例
	AliyunRegionChaosInstance *AliyunRegionChaos
	AliyunRegionChaosOnce     sync.Once

	// 节点混沌工程实例 - 全局单例
	AliyunNodeChaosInstance *AliyunNodeChaos
	AliyunNodeChaosOnce     sync.Once
)
