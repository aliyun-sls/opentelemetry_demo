package command

import (
	"context"
	"fmt"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/apimachinery/pkg/types"
	"k8s.io/client-go/dynamic"
	"time"
)

func DeleteCRD(dynamicClient *dynamic.DynamicClient, gvr schema.GroupVersionResource, name string) {
	Patch(dynamicClient, gvr, name)
	time.Sleep(2 * time.Second)
	Delete(dynamicClient, gvr, name)
	time.Sleep(2 * time.Second)
}

func Patch(dynamicClient *dynamic.DynamicClient, gvr schema.GroupVersionResource, name string) {
	patch := `{"metadata":{"finalizers":[]}}`
	_, err := dynamicClient.Resource(gvr).Patch(context.TODO(), name, types.MergePatchType, []byte(patch), metav1.PatchOptions{})
	if err != nil {
		fmt.Println("patch error")
		panic(err)
	}
}

func Delete(dynamicClient *dynamic.DynamicClient, gvr schema.GroupVersionResource, name string) {
	err := dynamicClient.Resource(gvr).Delete(context.TODO(), name, metav1.DeleteOptions{})
	if err != nil {
		fmt.Println("delete error")
		panic(err)
	}
}
