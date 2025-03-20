package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"time"
)

type APP struct {
}

func (a *APP) GetApp01(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Hello from /getapp01"))
}

type User struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

func main() {
	app := &APP{}
	fmt.Println("开启监听。。。。。")
	http.HandleFunc("/getapp01", app.GetApp01)

	// 将监听服务放在一个独立的goroutine中运行
	go func() {
		if err := http.ListenAndServe(":8080", nil); err != nil {
			log.Fatalf("Error starting server: %v", err)
		}
	}()

	urlLogin := "http://user.default.svc.cluster.local:8080/login"
	urlShelve := "http://product.default.svc.cluster.local:8080/shelve"

	// 创建一个每10秒触发一次的定时器
	ticker := time.NewTicker(10 * time.Second)
	defer ticker.Stop()

	for range ticker.C {
		// 创建请求体
		user := User{
			Username: "admin",
			Password: "admin",
		}
		jsonData, err := json.Marshal(user)
		if err != nil {
			log.Fatalf("Error marshalling user data: %v", err)
		}

		// POST 请求到登录接口
		resp, err := http.Post(urlLogin, "application/json", bytes.NewBuffer(jsonData))
		if err != nil {
			log.Fatalf("Error calling login endpoint: %v", err)
		}
		defer resp.Body.Close()

		body, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			log.Fatalf("Error reading login response body: %v", err)
		}
		log.Printf("Login Response: %s", body)

		// GET 请求到 shelve 接口，添加查询参数
		queryParams := "?products_cate=1&products_name=electronics&unit_price=100&products_unit=1&brand_id=123&seller_id=456"
		respShelve, err := http.Get(urlShelve + queryParams)
		if err != nil {
			log.Fatalf("Error calling shelve endpoint: %v", err)
		}
		defer respShelve.Body.Close()

		bodyShelve, err := ioutil.ReadAll(respShelve.Body)
		if err != nil {
			log.Fatalf("Error reading shelve response body: %v", err)
		}
		log.Printf("Shelve Response: %s", bodyShelve)
	}
}
