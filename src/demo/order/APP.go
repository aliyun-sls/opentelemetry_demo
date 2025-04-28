package main

import (
	"bytes"
	"context"
	"fmt"
	flagd "github.com/open-feature/go-sdk-contrib/providers/flagd/pkg"
	"github.com/open-feature/go-sdk/openfeature"
	"io"
	"io/ioutil"
	"log"
	"mime/multipart"
	"net/http"
	"os"
	"path/filepath"
)

var flagClient *openfeature.Client

type APP struct {
}

func (a *APP) Order(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("order"))
	urlShelve := "http://product:8080/api/v1/products/put_products"

	PutProducts(urlShelve)
	isSwitch := flagClient.String(
		context.Background(),
		"ServiceAbnormal",
		"",
		openfeature.EvaluationContext{},
	)
	log.Printf("获取feature ServiceAbnormal: %v", isSwitch)

	if isSwitch == "on" {
		service()
	}
}

func initFeatureFlag() {
	provider := flagd.NewProvider(
		flagd.WithHost("flagd"),
		flagd.WithPort(8013),
	)
	if err := openfeature.SetProvider(provider); err != nil {
		log.Printf("设置provider失败: %v", err)
		return
	}
	flagClient = openfeature.NewClient("order")
}

func main() {
	app := &APP{}

	initFeatureFlag()

	http.HandleFunc("/order", app.Order)
	log.Println("监听端口: 8080")
	if err := http.ListenAndServe(":8080", nil); err != nil {
		log.Fatalf("无法启动服务器: %v", err)
	}
}

func service() {
	response, err := http.Get("http://abnormal:8080/order")
	if err != nil {
		log.Printf("Error calling shelve endpoint: %v", err)
	}
	body, err := ioutil.ReadAll(response.Body)
	if err != nil {
		log.Printf("Error reading shelve response body: %v", err)
	}
	log.Printf("service Response: %s", body)
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
	writer.WriteField("inventory_num", "10")
	writer.WriteField("brand_id", "apple")
	writer.WriteField("seller_id", "1")
	writer.WriteField("products_status", "1")
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
