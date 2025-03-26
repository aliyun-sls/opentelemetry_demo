package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
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
	urlShelve := "http://product.default.svc.cluster.local:8080/api/v1/products/shelve/?"
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
	// GET 请求到 shelve 接口，添加查询参数
	path, _ := os.Getwd()
	fmt.Println(path)
	files, err := os.ReadDir(path + "/tupian/")
	if err != nil {
		log.Fatalf("failed to read directory: %v", err)
	}

	var product model.Product
	var picPath []uint8
	//var filePath = ""
	for _, file := range files {
		if !file.IsDir() {
			//filePath = file.Name()
			data, _ := product.ProductBasicType.ProductsPic.Value()
			picPath = data.([]uint8)
			fmt.Println(file.Name(), picPath)
			break // 只取第一张图片
		}
	}
	product = model.Product{
		ProductBasicType: model.ProductBasicType{
			ProductsName: "yichen",
			ProductsPic:  model.PicPath(picPath),
			UnitPrice:    100,
			ProductsUnit: 1,
		},
		BrandId:        123,
		SellerId:       456,
		ProductsCate:   1,
		ProductsDesc:   "electronics",
		ProductsStatus: model.Shelve,
	}
	fmt.Println(product.ProductsPic, string(product.ProductsPic), fmt.Sprintf("%s", product.ProductsPic))
	// 将product结构体转换为url.Values
	res := url.Values{}
	res.Add("products_name", product.ProductBasicType.ProductsName)
	res.Add("products_pic", string(product.ProductBasicType.ProductsPic))
	res.Add("unit_price", fmt.Sprintf("%d", product.ProductBasicType.UnitPrice))
	res.Add("products_unit", fmt.Sprintf("%d", product.ProductBasicType.ProductsUnit))
	res.Add("brand_id", fmt.Sprintf("%d", product.BrandId))
	res.Add("seller_id", fmt.Sprintf("%d", product.SellerId))
	res.Add("products_cate", fmt.Sprintf("%d", product.ProductsCate))
	res.Add("products_desc", product.ProductsDesc)
	res.Add("products_status", string(product.ProductsStatus))

	queryParams := res.Encode()
	fmt.Println(urlShelve + queryParams)
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
