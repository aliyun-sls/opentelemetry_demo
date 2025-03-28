package handlers

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
)

// GetApp01 调用 /getapp01 接口
func GetApp01() (string, error) {
	// 定义请求URL
	url := "http://goapp." + os.Getenv("NAMESPACE") + ".svc.cluster.local:8080/getapp01"

	// 发送GET请求
	resp, err := http.Get(url)
	if err != nil {
		return "", fmt.Errorf("请求失败: %v", err)
	}
	defer resp.Body.Close()

	// 读取响应体
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("读取响应体失败: %v", err)
	}

	// 返回响应内容
	return string(body), nil
}

func ToApp() {
	// 调用 GetApp01 函数
	response, err := GetApp01()
	if err != nil {
		fmt.Println("调用 /getapp01 接口失败:", err)
		return
	}
	fmt.Println("调用 /getapp01 接口成功，响应内容:", response)
}
