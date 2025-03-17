package oteldemo

import (
	"io/ioutil"
	"log"
	"net/http"
)

func ToApp01() {
	//service := os.Getenv("SERVICE_NAME")
	url := "http://goapp01:8080/getapp01"

	// 发送HTTP GET请求
	resp, err := http.Get(url)
	if err != nil {
		log.Fatalf("Error calling B Pod endpoint: %v", err)
	}
	defer resp.Body.Close()

	// 读取响应体
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatalf("Error reading response body: %v", err)
	}

	// 打印响应体
	log.Printf("Response from B Pod: %s", body)
}
