package command

import (
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
	Labels  string
	Port    string
	Percent string
	Timeout string
}

type Netdelay struct {
	Namespace string
	Labels    string
	Port      string
	Time      string
	Offset    string
	Timeout   string
}

type CpuAndMem struct {
	Labels  string
	Percent int64
}

var Dynamic *dynamic.DynamicClient
var Gvr schema.GroupVersionResource
var (
	FlagClient *openfeature.Client
	Once       sync.Once
)
