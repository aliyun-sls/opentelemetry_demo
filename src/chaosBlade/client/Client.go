package client

import (
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/client-go/dynamic"
	"k8s.io/client-go/rest"
)

func Client() (*dynamic.DynamicClient, schema.GroupVersionResource) {
	// 初始化flag客户端(只执行一次)

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
