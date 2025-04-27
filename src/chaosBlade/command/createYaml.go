package command

import (
	"bytes"
	"chaosBlade/client"
	"os"
	"strconv"
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
func (node *NodeLoss) NodeNetLossYaml() string {
	nodeid := os.Getenv("NODEID")
	labels := "alibabacloud.com/ecs-instance-id=" + nodeid
	name := "node-loss"
	if node.flagdValue != 0 {
		// 定义模板数据
		data := &Netloss{
			Name:    name,
			Labels:  labels,
			Percent: "100",
			Timeout: strconv.Itoa(int(node.flagdValue)),
		}

		path, _ := os.Getwd()
		return generateYAML(path+"/yaml/node_network_loss.yaml", data)
	} else {
		arr := client.ListCRD(Dynamic, Gvr)
		for _, s := range arr {
			if s == "node-loss" {
				DeleteCRD(Dynamic, Gvr, s)
			}
		}
	}
	return ""
}

func (region *RegionLoss) RegionNetLossYaml() string {
	zone := os.Getenv("ZONE")
	labels := "topology.kubernetes.io/zone=" + zone
	name := "region-loss"
	if region.flagdValue != 0 {
		// 定义模板数据
		data := &Netloss{
			Name:    name,
			Labels:  labels,
			Percent: "100",
			Timeout: strconv.Itoa(int(region.flagdValue)),
		}

		path, _ := os.Getwd()
		return generateYAML(path+"/yaml/node_network_loss.yaml", data)
	} else {
		arr := client.ListCRD(Dynamic, Gvr)
		for _, s := range arr {
			if s == "region-loss" {
				DeleteCRD(Dynamic, Gvr, s)
			}
		}
	}
	return ""
}

// PodNetDelay 配置pod网络延时压测
func (podnetdelay *PodNetDelay) PodNetDelayYaml() string {
	if podnetdelay.labels != "" {
		// 定义模板数据
		data := &Netdelay{
			Labels: podnetdelay.labels,
			Port:   "8080",
			Time:   "3000",
			Offset: "1000",
		}

		path, _ := os.Getwd()
		return generateYAML(path+"/yaml/pod_network_delay.yaml", data)
	} else {
		arr := client.ListCRD(Dynamic, Gvr)
		for _, s := range arr {
			if s == "pod-network-delay" {
				DeleteCRD(Dynamic, Gvr, s)
			}
		}
	}
	return ""
}

func (podcpu *PodCpu) PodCpuYaml() string {
	if podcpu.flagdValue != 0 {
		// 定义模板数据
		data := &CpuAndMem{
			Labels:  "app.kubernetes.io/name=cart",
			Percent: podcpu.flagdValue,
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

func (podmem *PodMem) PodMemYaml() string {
	if podmem.flagdValue != 0 {
		// 定义模板数据
		data := &CpuAndMem{
			Labels:  "app.kubernetes.io/name=cart",
			Percent: podmem.flagdValue,
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

func (nodecpu *NodeCpu) NodeCpuYaml() string {
	nodeid := os.Getenv("NODEID")
	labels := "alibabacloud.com/ecs-instance-id=" + nodeid
	if nodecpu.flagdValue != 0 {
		// 定义模板数据
		data := &CpuAndMem{
			Labels:  labels,
			Percent: nodecpu.flagdValue,
		}

		path, _ := os.Getwd()
		return generateYAML(path+"/yaml/node_cpu.yaml", data)
	} else {
		arr := client.ListCRD(Dynamic, Gvr)
		for _, s := range arr {
			if s == "node-cpu" {
				DeleteCRD(Dynamic, Gvr, s)
			}
		}
	}
	return ""
}

func (nodemem *NodeMem) NodeMemYaml() string {
	nodeid := os.Getenv("NODEID")
	labels := "alibabacloud.com/ecs-instance-id=" + nodeid
	if nodemem.flagdValue != 0 {
		// 定义模板数据
		data := &CpuAndMem{
			Labels:  labels,
			Percent: nodemem.flagdValue,
		}

		path, _ := os.Getwd()
		return generateYAML(path+"/yaml/node_mem.yaml", data)
	} else {
		arr := client.ListCRD(Dynamic, Gvr)
		for _, s := range arr {
			if s == "node-mem" {
				DeleteCRD(Dynamic, Gvr, s)
			}
		}
	}
	return ""
}
