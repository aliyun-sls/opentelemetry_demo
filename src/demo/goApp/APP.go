package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"sls-mall-go/common/model"
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
	urlShelve := "http://product.default.svc.cluster.local:8080/api/v1/products/shelve"
	urlGateway := "http://gateway:8080"

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
		//OSS.PushOSS()
		// GET 请求到 shelve 接口，添加查询参数
		path, _ := os.Getwd()
		files, err := os.ReadDir(path + "/tupian/")
		if err != nil {
			log.Fatalf("failed to read directory: %v", err)
		}

		var product model.Product
		var picPath []byte
		//var filePath = ""
		for _, file := range files {
			if !file.IsDir() {
				//filePath = file.Name()
				data, _ := product.ProductBasicType.ProductsPic.Value()
				picPath = data.([]byte)
				break // 只取第一张图片
			}
		}

		queryParams := fmt.Sprintf("?products_name=\"yichen\"&products_cate=1&products_name=electronics&unit_price=100&products_unit=1&brand_id=123&seller_id=456&products_pic=%s", picPath)
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

		respGateway, err := http.Get(urlGateway)
		if err != nil {
			log.Fatalf("Error calling login endpoint: %v", err)
		}
		defer respGateway.Body.Close()

		bodyGateway, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			log.Fatalf("Error reading login response body: %v", err)
		}
		log.Printf("Gateway Response: %s", bodyGateway)
	}
}
