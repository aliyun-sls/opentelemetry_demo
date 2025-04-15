package command

import (
	"bytes"
	"chaosBlade/client"
	"os"
	"text/template"
)

// 通用函数，用于生成YAML文件
func generateYAML(templateFile string, data interface{}) string {
	// 解析模板文件
	tmpl, err := template.ParseFiles(templateFile)
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

// 配置node网络丢包压测
func NodeNetLossYaml(labels string) string {
	if labels != "" {
		// 定义模板数据
		data := &Netloss{
			Name:    "chaosblade-node-loss",
			Labels:  labels,
			Percent: "100",
			Timeout: "500",
		}

		path, _ := os.Getwd()
		return generateYAML(path+"/yaml/node_network_loss.yaml", data)
	} else {
		arr := client.ListCRD(Dynamic, Gvr)
		for _, s := range arr {
			if s == "chaosblade-node-loss" {
				DeleteCRD(Dynamic, Gvr, s)
			}
		}
	}
	return ""
}

// 配置RDS断连
func RDSlossYaml(labels string) string {
	if labels != "" {
		// 定义模板数据
		data := &Netloss{
			Name:    "chaosblade-rds-loss",
			Labels:  labels,
			Percent: "100",
			Port:    "3306",
			Timeout: "500",
		}

		path, _ := os.Getwd()
		return generateYAML(path+"/yaml/pod_network_loss.yaml", data)
	} else {
		arr := client.ListCRD(Dynamic, Gvr)
		for _, s := range arr {
			if s == "chaosblade-rds-loss" {
				DeleteCRD(Dynamic, Gvr, s)
			}
		}
	}
	return ""
}

// PodNetDelay 配置pod网络延时压测
func PodNetDelayYaml(labels string) string {
	if labels != "" {
		// 定义模板数据
		data := &Netdelay{
			Namespace: "default",
			Labels:    labels,
			Port:      "8080",
			Time:      "3000",
			Offset:    "1000",
			Timeout:   "500",
		}

		path, _ := os.Getwd()
		return generateYAML(path+"/yaml/pod_network_delay.yaml", data)
	} else {
		arr := client.ListCRD(Dynamic, Gvr)
		for _, s := range arr {
			if s == "chaosblade-rds-loss" {
				DeleteCRD(Dynamic, Gvr, s)
			}
		}
	}
	return ""
}

func PodCpuYaml(percent int64) string {
	if percent != 0 {
		// 定义模板数据
		data := &CpuAndMem{
			Labels:  "app.kubernetes.io/name=cart",
			Percent: percent,
		}

		path, _ := os.Getwd()
		return generateYAML(path+"/yaml/pod_cpu.yaml", data)
	} else {
		arr := client.ListCRD(Dynamic, Gvr)
		for _, s := range arr {
			if s == "cpu-load" {
				DeleteCRD(Dynamic, Gvr, s)
			}
		}
	}
	return ""
}

func PodMemYaml(percent int64) string {
	if percent != 0 {
		// 定义模板数据
		data := &CpuAndMem{
			Labels:  "app.kubernetes.io/name=cart",
			Percent: percent,
		}

		path, _ := os.Getwd()
		return generateYAML(path+"/yaml/pod_mem.yaml", data)
	} else {
		arr := client.ListCRD(Dynamic, Gvr)
		for _, s := range arr {
			if s == "mem-load" {
				DeleteCRD(Dynamic, Gvr, s)
			}
		}
	}
	return ""
}
