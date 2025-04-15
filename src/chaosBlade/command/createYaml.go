package command

import (
	"bytes"
	"chaosBlade/client"
	"context"
	"fmt"
	"github.com/open-feature/go-sdk/openfeature"
	"log"
	"os"
	"text/template"
)

// 配置node网络丢包压测
func NodeNetLossYaml() string {

	// 获取feature flag值
	labels := flagClient.String(
		context.Background(),
		"nodeLoss",
		"",
		openfeature.EvaluationContext{},
	)
	log.Printf("获取feature : %v", labels)
	if labels != "" {
		// 定义模板数据
		data := &Netloss{
			Name:    "chaosblade-node-loss",
			Labels:  labels,
			Percent: "100",
			Port:    "8080",
			Timeout: "500",
		}

		// 解析模板文件
		path, _ := os.Getwd()
		tmpl, err := template.ParseFiles(path + "/yaml/node_network_loss.yaml")
		if err != nil {
			panic(err)
		}

		// 渲染模板并保存到文件
		var buf bytes.Buffer
		err = tmpl.Execute(&buf, data)
		if err != nil {
			panic(err)
		}
		return buf.String()
	} else {
		DeleteCRD(Dynamic, Gvr, "chaosblade-node-loss")
	}
	return ""
}

// 配置RDS断连
func RDSlossYaml(labels string) string {
	if labels != "" {
		fmt.Println(labels)
		// 定义模板数据
		data := &Netloss{
			Name:    "chaosblade-rds-loss",
			Labels:  labels,
			Percent: "100",
			Port:    "3306",
			Timeout: "500",
		}

		// 解析模板文件
		path, _ := os.Getwd()
		tmpl, err := template.ParseFiles(path + "/yaml/node_network_loss.yaml")
		if err != nil {
			fmt.Println("模版error")
			panic(err)
		}

		// 渲染模板并保存到文件
		var buf bytes.Buffer
		err = tmpl.Execute(&buf, data)
		if err != nil {
			panic(err)
		}
		return buf.String()
	} else {
		fmt.Println(labels)
		arr := client.ListCRD(Dynamic, Gvr)
		if len(arr) != 0 {
			DeleteCRD(Dynamic, Gvr, "chaosblade-rds-loss")
		}

	}
	return ""
}

// PodNetDelay 配置pod网络延时压测
func PodNetDelay() string {

	namespace := os.Getenv("NAMESPACE")
	labels := os.Getenv("LABELS")
	port := os.Getenv("PORT")
	time := os.Getenv("TIME")
	offset := os.Getenv("OFFSET")
	timeout := os.Getenv("TIMEOUT")

	// 定义模板数据
	data := &Netdelay{
		Namespace: namespace,
		Labels:    labels,
		Port:      port,
		Time:      time,
		Offset:    offset,
		Timeout:   timeout,
	}

	// 解析模板文件
	tmpl, err := template.ParseFiles("/Users/chenxiaohui/opentelemetry_demo/src/chaosBlade/chaosblade/yaml/node_network_loss.yaml")
	if err != nil {
		panic(err)
	}

	// 渲染模板并保存到文件
	var buf bytes.Buffer
	err = tmpl.Execute(&buf, data)
	if err != nil {
		panic(err)
	}
	return buf.String()
}
