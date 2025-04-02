package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
)

type CPU struct {
}

// ConsumeCPU 通过遍历一个大数组来消耗 CPU 资源
func (a *CPU) ConsumeCPU() {
	// 创建一个包含几百万个元素的数组
	const size = 100000000
	arr := make([]int, size)

	// 填充数组
	for i := 0; i < size; i++ {
		arr[i] = i
	}

	// 遍历数组并进行一些简单的计算来消耗 CPU
	sum := 0
	for _, value := range arr {
		sum += value
	}
	fmt.Println("Sum:", sum)
}

func (a *CPU) cpu(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	// 返回 JSON 格式的响应
	response := map[string]string{
		"message": "Hello from /cpu",
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(response)

	a.ConsumeCPU()
}

func main() {
	app := &CPU{}
	fmt.Println("开启监听。。。。。")
	http.HandleFunc("/cpu", app.cpu)

	// 将监听服务放在一个独立的goroutine中运行
	if err := http.ListenAndServe(":8080", nil); err != nil {
		log.Fatalf("Error starting server: %v", err)
	}
}
