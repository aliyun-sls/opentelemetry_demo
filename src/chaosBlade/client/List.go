package client

import (
	"context"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/client-go/dynamic"
)

func ListCRD(dynamicClient *dynamic.DynamicClient, gvr schema.GroupVersionResource) []string {
	// 获取资源
	unStructData, err := dynamicClient.Resource(gvr).List(context.TODO(), metav1.ListOptions{})
	if err != nil {
		panic(err.Error())
	}

	// 打印结果
	var res corev1.PodList
	var arr []string
	runtime.DefaultUnstructuredConverter.FromUnstructured(unStructData.UnstructuredContent(), &res)
	for _, item := range res.Items {
		arr = append(arr, item.Name)
		//print(item.Name)
	}
	return arr
}

func (g *Gvrnode) ListNode(dynamicClient *dynamic.DynamicClient) []string {
	unStructData, err := dynamicClient.Resource(g.Gvr).List(context.TODO(), metav1.ListOptions{})
	if err != nil {
		panic(err.Error())
	}

	// 打印结果
	var res corev1.PodList
	var arr []string
	err = runtime.DefaultUnstructuredConverter.FromUnstructured(unStructData.UnstructuredContent(), &res)
	if err != nil {
		return nil
	}
	for _, item := range res.Items {
		arr = append(arr, item.Name)
		//print(item.Name)
	}
	return arr
}

func (g *Gvrpod) ListPod(dynamicClient *dynamic.DynamicClient) []string {
	// 获取资源
	//gvrPod := schema.GroupVersionResource{
	//	Group:    "",
	//	Version:  "v1",
	//	Resource: "pods",
	//}

	unStructData, err := dynamicClient.Resource(g.Gvr).Namespace("test").List(context.TODO(), metav1.ListOptions{})
	if err != nil {
		panic(err.Error())
	}

	// 打印结果
	var res corev1.PodList
	var arr []string
	err = runtime.DefaultUnstructuredConverter.FromUnstructured(unStructData.UnstructuredContent(), &res)
	if err != nil {
		return nil
	}
	for _, item := range res.Items {
		arr = append(arr, item.Name)
		//print(item.Name)
	}
	return arr
}
