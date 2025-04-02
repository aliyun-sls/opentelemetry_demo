package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
)

type CPU struct {
}

// ConsumeCPU 通过循环计算来消耗 CPU 资源
func ConsumeCPU() {
	// 通过循环计算来消耗 CPU
	const iterations = 1000000000
	sum := 0
	for i := 0; i < iterations; i++ {
		sum += i
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

	ConsumeCPU()
}

func run() {
	ConsumeCPU()
}

func main() {
	app := &CPU{}
	fmt.Println("开启监听。。。。。")
	http.HandleFunc("/cpu", app.cpu)
	go func() {
		for {
			run()
		}

	}()
	// 将监听服务放在一个独立的goroutine中运行
	if err := http.ListenAndServe(":8080", nil); err != nil {
		log.Fatalf("Error starting server: %v", err)
	}

}
