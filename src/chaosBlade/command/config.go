package command

import (
	"github.com/open-feature/go-sdk/openfeature"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/client-go/dynamic"
	"sync"
)

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
	flagClient *openfeature.Client
	once       sync.Once
)
