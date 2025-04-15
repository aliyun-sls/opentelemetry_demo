package main

import (
	"context"
	"fmt"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
)

func main() {
	// 创建集群内配置，并指定使用 system:controller:pod-garbage-collector 的 Token
	config, err := rest.InClusterConfig()
	if err != nil {
		panic(err.Error())
	}

	// 设置使用 system:controller:pod-garbage-collector 的 Token
	config.BearerTokenFile = "/var/run/secrets/kubernetes.io/serviceaccount/token"

	// 创建clientset
	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		panic(err.Error())
	}

	// 示例：列出所有pod
	pods, err := clientset.CoreV1().Pods("default").List(context.TODO(), metav1.ListOptions{})
	if err != nil {
		panic(err.Error())
	}

	fmt.Println("Pods in the cluster:")
	for _, pod := range pods.Items {
		fmt.Printf("- %s in namespace %s\n", pod.Name, pod.Namespace)
	}
}
