package command

import (
	"context"
	"fmt"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/apimachinery/pkg/runtime/serializer/yaml"
	"k8s.io/client-go/dynamic"
	"time"
)

func createCRD(dynamic *dynamic.DynamicClient, gvr schema.GroupVersionResource, file string) string {
	// 创建资源
	decoder := yaml.NewDecodingSerializer(unstructured.UnstructuredJSONScheme)

	obj := &unstructured.Unstructured{}

	_, _, err := decoder.Decode([]byte(file), nil, obj)
	if err != nil {
		panic(err.Error())
	}

	res, err := dynamic.Resource(gvr).Apply(context.TODO(), obj.GetName(), obj, metav1.ApplyOptions{FieldManager: "chaosBlade"})
	if err != nil {
		panic(err.Error())
	}
	time.Sleep(2 * time.Second)
	return GetStatus(dynamic, gvr, res.GetName())
}

func GetStatus(dynamic *dynamic.DynamicClient, gvr schema.GroupVersionResource, name string) string {
	// 获取资源
	obj, err := dynamic.Resource(gvr).Get(context.TODO(), name, metav1.GetOptions{})
	if err != nil {
		fmt.Println("get status error")
		panic(err.Error())
	}

	// 将资源对象转换为 JSON 格式并打印
	jsonData, err := obj.MarshalJSON()
	if err != nil {
		panic(err.Error())
	}
	return string(jsonData)
}
