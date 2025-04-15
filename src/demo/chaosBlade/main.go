package main

import (
	"context"
	"fmt"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
)

func main() {
	// 尝试创建集群内配置
	config, err := rest.InClusterConfig()
	/*if err != nil {
		// 如果不是在集群内运行，尝试使用kubeconfig
		home, _ := os.UserHomeDir()
		kubeconfig := filepath.Join(home, ".kube", "config")
		config, err = clientcmd.BuildConfigFromFlags("", kubeconfig)
		if err != nil {
			panic(err.Error())
		}
	}*/

	// 创建clientset
	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		panic(err.Error())
	}

	// 示例：列出所有pod
	pods, err := clientset.CoreV1().Pods("").List(context.TODO(), metav1.ListOptions{})
	if err != nil {
		panic(err.Error())
	}

	fmt.Println("Pods in the cluster:")
	for _, pod := range pods.Items {
		fmt.Printf("- %s in namespace %s\n", pod.Name, pod.Namespace)
	}
}
