package client

import (
	flagd "github.com/open-feature/go-sdk-contrib/providers/flagd/pkg"
	"github.com/open-feature/go-sdk/openfeature"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/client-go/dynamic"
	"k8s.io/client-go/rest"
	"log"
	"sync"
)

var flagClient *openfeature.Client
var once sync.Once

func initFeatureFlag() {
	once.Do(func() {
		provider := flagd.NewProvider(
			flagd.WithHost("flagd"),
			flagd.WithPort(8013),
		)
		if err := openfeature.SetProvider(provider); err != nil {
			log.Printf("设置provider失败: %v", err)
			return
		}
		flagClient = openfeature.NewClient("product")
	})
}

func Client() (*dynamic.DynamicClient, schema.GroupVersionResource) {
	// 初始化flag客户端(只执行一次)
	initFeatureFlag()

	config, err := rest.InClusterConfig()
	if err != nil {
		panic(err.Error())
	}

	dynamicClient, err := dynamic.NewForConfig(config)
	if err != nil {
		panic(err.Error())
	}

	gvr := schema.GroupVersionResource{
		Group:    "chaosblade.io",
		Version:  "v1alpha1",
		Resource: "chaosblades",
	}

	return dynamicClient, gvr
}
