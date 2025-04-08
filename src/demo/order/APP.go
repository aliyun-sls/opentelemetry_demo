package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"mime/multipart"
	"net/http"
	"os"
	"path/filepath"
)

type APP struct {
}

func (a *APP) GetApp01(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("order"))
	//urlCpu := "http://gocpu." + os.Getenv("NAMESPACE") + ".svc.cluster.local:8080/cpu"
	urlShelve := "http://product." + os.Getenv("NAMESPACE") + ".svc.cluster.local:8080/api/v1/products/put_products"

	//cpu(urlCpu)
	PutProducts(urlShelve)
}

type User struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

func main() {
	app := &APP{}
	fmt.Println("开启监听。。。。。")
	http.HandleFunc("/order", app.GetApp01)

	// 将监听服务放在一个独立的goroutine中运行
	http.ListenAndServe(":8080", nil)

}

func cpu(urlCpu string) {
	body, err := http.Get(urlCpu)
	if err != nil {
		log.Printf("Error calling cpu endpoint: %v", err)
		return
	}
	log.Printf("gocpu Response: %s", body)
}

func user(urlLogin string) {
	user := User{
		Username: "admin",
		Password: "admin",
	}
	jsonData, err := json.Marshal(user)
	if err != nil {
		log.Printf("Error marshalling user data: %v", err)
		return
	}

	// POST 请求到登录接口
	resp, err := http.Post(urlLogin, "application/json", bytes.NewBuffer(jsonData))
	if err != nil {
		log.Printf("Error calling login endpoint: %v", err)
		return
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Printf("Error reading login response body: %v", err)
		return
	}
	log.Printf("Login Response: %s", body)
}

func PutProducts(urlShelve string) {
	// GET 请求到 shelve 接口，添加查询参数
	path, _ := os.Getwd()
	path = path + "/tupian/"
	files, err := os.ReadDir(path)
	if err != nil {
		log.Printf("failed to read directory: %v", err)
	}
	var data *os.File

	for _, file := range files {
		if !file.IsDir() {
			filePath := file.Name()
			data, err = os.Open(path + filePath)
			if err != nil {
				log.Printf("打开文件失败: %v", err)
			}
			break // 只取第一张图片
		}
	}

	defer data.Close()
	//创建multipart writer
	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)

	//创建文件表单字段
	writer.WriteField("products_name", "Test Product")
	writer.WriteField("products_price", "100")
	writer.WriteField("products_unit", "1")
	writer.WriteField("products_desc", "This is a test product")
	writer.WriteField("products_category", "1")
	writer.WriteField("products_status", "1")
	writer.WriteField("seller_id", "1")
	writer.WriteField("brand_id", "1")
	writer.WriteField("inventory_name", "Test Inventory")
	writer.WriteField("inventory_address", "Test Address")
	writer.WriteField("inventory_num", "10")
	part, err := writer.CreateFormFile("products_pic", filepath.Base(data.Name()))
	if err != nil {
		fmt.Println(err)
	}

	//将文件内容拷贝到表单
	_, err = io.Copy(part, data)
	if err != nil {
		log.Printf("拷贝文件内容失败: %v", err)
	}

	//关闭writer完成表单构建
	err = writer.Close()
	if err != nil {
		log.Printf("关闭writer失败: %v", err)
	}

	respShelve, err := http.Post(urlShelve, writer.FormDataContentType(), body)
	if err != nil {
		log.Printf("Error calling shelve endpoint: %v", err)
		return // 直接返回，避免后续操作
	}
	if respShelve != nil {
		defer respShelve.Body.Close()
	}

	bodyShelve, err := ioutil.ReadAll(respShelve.Body)
	if err != nil {
		log.Printf("Error reading shelve response body: %v", err)
	}
	log.Printf("Shelve Response: %s", bodyShelve)
}
