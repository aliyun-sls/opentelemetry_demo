package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"math/rand"
	"mime/multipart"
	"net/http"
	"os"
	"path/filepath"
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
		user(urlLogin)

		shelve(urlShelve)

		gateway(urlGateway)
	}
}
func gateway(urlGateway string) {
	if rand.Int()%5 == 0 {
		time.Sleep(5 * time.Second)
	} else if rand.Int()%3 == 0 {
		time.Sleep(2 * time.Second)
	}
	respGateway, err := http.Get(urlGateway)
	if err != nil {
		log.Fatalf("Error calling login endpoint: %v", err)
	}
	defer respGateway.Body.Close()

	bodyGateway, err := ioutil.ReadAll(respGateway.Body)
	if err != nil {
		log.Fatalf("Error reading login response body: %v", err)
	}
	log.Printf("Gateway Response: %s", bodyGateway)
}
func user(urlLogin string) {
	if rand.Int()%5 == 0 {
		time.Sleep(5 * time.Second)
	} else if rand.Int()%3 == 0 {
		time.Sleep(2 * time.Second)
	}
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
}

func shelve(urlShelve string) {
	/*if rand.Int()%5 == 0 {
		time.Sleep(5 * time.Second)
	} else if rand.Int()%3 == 0 {
		time.Sleep(2 * time.Second)
	}*/
	// GET 请求到 shelve 接口，添加查询参数
	path, _ := os.Getwd()
	path = path + "/tupian/"
	files, err := os.ReadDir(path)
	if err != nil {
		log.Fatalf("failed to read directory: %v", err)
	}
	var data *os.File
	var product model.Product
	var picPath []uint8
	for _, file := range files {
		if !file.IsDir() {
			filePath := file.Name()
			data, err = os.Open(filePath)
			if err != nil {
				fmt.Println("打开文件失败:", err)
				return
			}
			break // 只取第一张图片
		}
	}
	//创建multipart writer
	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)

	//创建文件表单字段
	part, err := writer.CreateFormFile("image", filepath.Base(path))
	if err != nil {
		fmt.Println("创建表单字段失败:", err)
		return
	}

	//将文件内容拷贝到表单
	_, err = io.Copy(part, data)
	if err != nil {
		fmt.Println("拷贝文件内容失败:", err)
		return
	}

	product = model.Product{
		ProductBasicType: model.ProductBasicType{
			ProductsName: "yichen",
			ProductsPic:  picPath,
			UnitPrice:    100,
			ProductsUnit: 1,
		},
		BrandId:        123,
		SellerId:       456,
		ProductsCate:   1,
		ProductsDesc:   "electronics",
		ProductsStatus: model.Shelve,
	}

	writer.WriteField("brand_id", "123")
	writer.WriteField("product_name", "测试产品")

	//关闭writer完成表单构建
	err = writer.Close()
	if err != nil {
		fmt.Println("关闭writer失败:", err)
		return
	}

	respShelve, err := http.Post(urlShelve, "Content-Type", body)
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
