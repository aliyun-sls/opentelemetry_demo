package service

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	flagd "github.com/open-feature/go-sdk-contrib/providers/flagd/pkg"
	"github.com/open-feature/go-sdk/openfeature"
	"github.com/rs/xid"
	"io"
	"log"
	"mime/multipart"
	"net/http"
	"os"
	"path/filepath"
	"sls-mall-go/common/model"
	"sls-mall-go/common/util"
	"sync"
)

type OrderStatus uint

const (
	// WaitForPayment 待付款
	WaitForPayment OrderStatus = 1
	// WaitForSending 待发货
	WaitForSending OrderStatus = 2
	// WaitForReceiving 待收货
	WaitForReceiving OrderStatus = 3
	// Complete 完成
	Complete OrderStatus = 4
	// Cancel 取消
	Cancel OrderStatus = 5
	//已评论
	Commented OrderStatus = 6
)

func GenId() string {
	return xid.New().String()
}

type OrderResult struct {
	OrderId            string       `json:"order_id,omitempty"`
	ShippingTrackingId string       `json:"shipping_tracking_id,omitempty"`
	ShippingCost       *Money       `json:"shipping_cost,omitempty"`
	ShippingAddress    *Address     `json:"shipping_address,omitempty"`
	Items              []*OrderItem `json:"items,omitempty"`
	TotalPrice         float64      `json:"total_price,omitempty"`
	UserId             int64        `json:"user_id,omitempty"`
}

type Money struct {
	// The 3-letter currency code defined in ISO 4217.
	CurrencyCode string `json:"currency_code,omitempty"`
	// The whole units of the amount.
	// For example if `currencyCode` is `"USD"`, then 1 unit is one US dollar.
	Units int64 `json:"units,omitempty"`
	// Number of nano (10^-9) units of the amount.
	// The value must be between -999,999,999 and +999,999,999 inclusive.
	// If `units` is positive, `nanos` must be positive or zero.
	// If `units` is zero, `nanos` can be positive, zero, or negative.
	// If `units` is negative, `nanos` must be negative or zero.
	// For example $-1.75 is represented as `units`=-1 and `nanos`=-750,000,000.
	Nanos int32 `json:"nanos,omitempty"`
}

type Address struct {
	StreetAddress string `json:"street_address,omitempty"`
	City          string `json:"city,omitempty"`
	State         string `json:"state,omitempty"`
	Country       string `json:"country,omitempty"`
	ZipCode       string `json:"zip_code,omitempty"`
}

type OrderItem struct {
	ProductId string   `json:"product_id"`
	Quantity  int64    `json:"quantity"`
	Cost      *Money   `protobuf:"bytes,2,opt,name=cost,proto3" json:"cost,omitempty"`
	Product   *Product `json:"product"`
}

type Product struct {
	Name        string `json:"name"`
	Picture     string `json:"picture"`
	Description string `json:"description"`
	PriceUsd    *Money `json:"price_usd"`
}

type LogisticsStatus uint

const (
	// Shipping 发货 src
	Shipping LogisticsStatus = 1
	// Collecting 揽收 src
	Collecting LogisticsStatus = 2
	// Transportation 运输 src
	Transportation LogisticsStatus = 3
	// Delivery 派送 dst
	Delivery LogisticsStatus = 4
	// Signing 签收 dst
	Signing LogisticsStatus = 5
)

type Action uint

const (
	Create Action = 1
)

type LogisticsMsg struct {
	UserId           int64           `gorm:"column:user_id" json:"user_id"`
	OrderId          string          `gorm:"column:order_id;primaryKey" json:"order_id"`
	Action           Action          `gorm:"column:action" json:"action"`
	LogisticStatus   LogisticsStatus `json:"logistic_status" gorm:"size:10;index;not null;comment:'物流状态'"`
	LogisticPosition string          `json:"logistic_position" gorm:"size:1000;not null;comment:'物流位置'"`
}

type Order struct {
	model.Model
	UserId             int64           `gorm:"column:user_id" json:"user_id" form:"user_id"`
	OrderId            string          `gorm:"column:order_id;primaryKey" json:"order_id" form:"order_id"`
	ShippingTrackingId string          `gorm:"column:shipping_tracking_id" json:"shipping_tracking_id"`
	CurrencyCode       string          `gorm:"column:currency_code" json:"currency_code"`
	Units              int64           `gorm:"column:units" json:"units"`
	Nanos              int32           `gorm:"column:nanos" json:"nanos"`
	StreetAddress      string          `gorm:"column:street_address" json:"street_address"`
	City               string          `gorm:"column:city" json:"city"`
	State              string          `gorm:"column:state" json:"state"`
	Country            string          `gorm:"column:country" json:"country"`
	ZipCode            string          `gorm:"column:zip_code" json:"zip_code"`
	OrderStatus        OrderStatus     `gorm:"column:order_status" json:"order_status"`
	LogisticStatus     LogisticsStatus `gorm:"column:logistic_status" form:"logistic_status" json:"logistic_status"`
	OrderDetails       []OrderDetail   `json:"order_details" gorm:"foreignKey:OrderId;references:OrderId"`
	TotalPrice         float64         `gorm:"column:total_price" json:"total_price"`
}

// TableName 方法指定了 Order 结构体对应的表名为 `order`
func (Order) TableName() string {
	return "order" // 使用单数形式
}

type OrderDetail struct {
	OrderDetailId  uint   `json:"order_detail_id" gorm:"primaryKey"`
	ProductId      string `json:"product_id"`
	OrderId        string `json:"order_id"`
	Quantity       int64  `json:"quantity"`
	ProductName    string `json:"product_name"`
	ProductPicture string `json:"product_picture"`
	Units          int64  `gorm:"column:units" json:"units"`
	Nanos          int32  `gorm:"column:nanos" json:"nanos"`
	UserId         int64  `gorm:"column:user_id" json:"user_id"`
	Description    string `json:"description"`
}

func (OrderDetail) TableName() string {
	return "order_detail" // 使用单数形式
}

type ListOrderRequest struct {
	Order
	util.Page
}

var flagClient *openfeature.Client

func CreateOrder(c *gin.Context) {
	// 处理请求体数据
	ctx := c.Request.Context()
	var orderRequest OrderResult
	err := c.BindJSON(&orderRequest)
	if err != nil {
		util.Status400(c, err)
		return
	}
	marshal, err := json.Marshal(orderRequest)
	log.Println(string(marshal))

	/*split := strings.Split(orderRequest.OrderId, "$")
	var orderId string
	var userId int64
	if len(split) == 2 {
		orderId = split[0]
		i, err := strconv.ParseInt(split[1], 10, 0)
		if err != nil {
			log.Fatal(err)
		}
		userId = i
	}*/
	// 将 OrderResult 转换为 Order
	var orderDetails []OrderDetail
	for _, item := range orderRequest.Items {
		orderDetail := OrderDetail{
			ProductId:      item.ProductId,
			Quantity:       item.Quantity,
			UserId:         orderRequest.UserId,
			ProductName:    item.Product.Name,
			ProductPicture: item.Product.Picture,
			Description:    item.Product.Description,
			Units:          item.Cost.Units,
			Nanos:          item.Cost.Nanos,
			OrderId:        orderRequest.OrderId,
		}
		orderDetails = append(orderDetails, orderDetail)
	}

	order := Order{
		OrderId:            orderRequest.OrderId,
		ShippingTrackingId: orderRequest.ShippingTrackingId,
		CurrencyCode:       orderRequest.ShippingCost.CurrencyCode,
		Units:              orderRequest.ShippingCost.Units,
		Nanos:              orderRequest.ShippingCost.Nanos,
		StreetAddress:      orderRequest.ShippingAddress.StreetAddress,
		City:               orderRequest.ShippingAddress.City,
		State:              orderRequest.ShippingAddress.State,
		Country:            orderRequest.ShippingAddress.Country,
		ZipCode:            orderRequest.ShippingAddress.ZipCode,
		OrderStatus:        WaitForPayment, // 初始状态为待付款
		OrderDetails:       orderDetails,
		UserId:             orderRequest.UserId,
		TotalPrice:         orderRequest.TotalPrice,
	}

	err = util.MDB.WithContext(ctx).Create(&order).Error
	if err != nil {
		util.Status500(c, err)
		return
	}
	util.Status200(c, order)
	//创建订单成功后 发送支付订单流程
	var wg sync.WaitGroup
	background := context.Background()
	wg.Add(1)
	go func(o *Order) {
		defer wg.Done()
		if err := ServiceCallPost(background, os.Getenv("PayHost"), "/pay/Create", o, &util.Result{}); err != nil {
			log.Printf("Failed to call payment service: %v", err)
		}
	}(&order)
	urlShelve := "http://product:8080/api/v1/products/put_products"

	if err := PutProducts(urlShelve); err != nil {
		log.Printf("Failed to put products: %v", err)
	}
	isSwitch := flagClient.String(
		context.Background(),
		"ServiceAbnormal",
		"",
		openfeature.EvaluationContext{},
	)
	log.Printf("获取feature ServiceAbnormal: %v", isSwitch)

	if isSwitch == "on" {
		Service()
	}
}

func PayOrder(c *gin.Context) {
	ctx := c.Request.Context()
	var order Order
	err := c.BindJSON(&order)
	if err != nil {
		util.Status400(c, err)
		return
	}

	tx := util.MDB.WithContext(ctx).Model(&order).Where("user_id = ? and order_id = ?", order.UserId, order.OrderId).
		Update("order_status", WaitForSending)
	err = tx.Error
	if err != nil {
		util.Status500(c, err)
		return
	}
	util.Status200(c, true)
}

// Shipping 发货
func CreateShipping(c *gin.Context) {
	ctx := c.Request.Context()
	var order Order
	orderId := c.Query("order_id")
	err := util.MDB.WithContext(ctx).Model(&order).
		Where("order_id = ?", orderId).
		Updates(Order{
			OrderStatus:    WaitForReceiving,
			LogisticStatus: Shipping,
		}).Error
	if err != nil {
		util.Status500(c, err)
		return
	}

	var wg sync.WaitGroup
	wg.Add(1)
	go func(order *Order) {
		// 创建物流 通过kafka通知 创建物流信息
		err = SendKafka(LogisticsMsg{
			OrderId:          order.OrderId,
			UserId:           order.UserId,
			Action:           Create,
			LogisticStatus:   Shipping,
			LogisticPosition: "",
		})
		if err != nil {
			log.Printf("SendKafka is fail: %v", err)
		}
		defer wg.Done()
	}(&order)

	util.Status200(c, true)
}

func ListOrderAndDetails(c *gin.Context) {
	ctx := c.Request.Context()

	var listOrderRequest ListOrderRequest
	err := c.BindJSON(&listOrderRequest)
	if err != nil {
		util.Status400(c, err)
		return
	}
	var orders []Order
	if listOrderRequest.Limit == 0 {
		listOrderRequest.Limit = 10
	}
	tx := util.MDB.WithContext(ctx).Model(&listOrderRequest.Order).Preload("OrderDetails")
	err = tx.Limit(listOrderRequest.Limit).Offset(listOrderRequest.Offset).Where(&listOrderRequest.Order).Order("id desc").Find(&orders).Error
	if err != nil {
		util.Status500(c, err)
		return
	}
	util.Status200(c, orders)
}

func GetOrder(c *gin.Context) {
	ctx := c.Request.Context()

	var order Order
	err := c.BindJSON(&order)
	if err != nil {
		util.Status500(c, err)
		return
	}
	err = util.MDB.WithContext(ctx).Model(&order).Preload("OrderDetails").Where("user_id = ? and order_id = ?", order.UserId, order.OrderId).Find(&order).Error
	if err != nil {
		util.Status500(c, err)
		return
	}
	util.Status200(c, order)
}

// LogisticsStatus 物流
func LogisticStatusUpdate(c *gin.Context) {
	ctx := c.Request.Context()
	var order Order
	err := c.BindJSON(&order)
	if err != nil {
		util.Status400(c, err)
		return
	}
	if order.LogisticStatus == Signing {
		err = util.MDB.WithContext(ctx).Model(&order).Where("order_id =?", order.OrderId).
			Updates(Order{LogisticStatus: order.LogisticStatus, OrderStatus: Complete}).Error
	}
	err = util.MDB.WithContext(ctx).Model(&order).Where("order_id =?", order.OrderId).
		Updates(Order{LogisticStatus: order.LogisticStatus}).Error
	if err != nil {
		util.Status500(c, err)
		return
	}
	util.Status200(c, true)
}

func InitFeatureFlag() {
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

func Service() {
	response, err := http.Get("http://abnormal:8080/order")
	if err != nil {
		log.Printf("Error calling shelve endpoint: %v", err)
	}
	body, err := io.ReadAll(response.Body)
	if err != nil {
		log.Printf("Error reading shelve response body: %v", err)
	}
	log.Printf("service Response: %s", body)
}

func PutProducts(urlShelve string) error {
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

	if data == nil {
		log.Printf("未找到文件")
		return fmt.Errorf("未找到文件")
	}
	defer func() {
		if err := data.Close(); err != nil {
			log.Printf("关闭文件时发生错误: %v", err)
		}
	}()
	//创建multipart writer
	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)

	// 创建文件表单字段
	if err := writer.WriteField("inventory_num", "10"); err != nil {
		log.Printf("写入inventory_num字段失败: %v", err)
		return err
	}
	if err := writer.WriteField("brand_id", "apple"); err != nil {
		log.Printf("写入brand_id字段失败: %v", err)
		return err
	}
	if err := writer.WriteField("seller_id", "1"); err != nil {
		log.Printf("写入seller_id字段失败: %v", err)
		return err
	}
	if err := writer.WriteField("products_status", "1"); err != nil {
		log.Printf("写入products_status字段失败: %v", err)
		return err
	}
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
		return err // 返回错误信息以便调用者处理
	}
	// 检查 respShelve 是否为 nil
	if respShelve == nil {
		log.Printf("响应为空，无法读取响应体")
		return fmt.Errorf("响应为空")
	} else {
		defer func() {
			if err := respShelve.Body.Close(); err != nil {
				log.Printf("Error closing response body: %v", err)
			}
		}()
	}
	bodyShelve, err := io.ReadAll(respShelve.Body)
	if err != nil {
		log.Printf("Error reading shelve response body: %v", err)
	}
	log.Printf("Shelve Response: %s", bodyShelve)

	return nil
}
