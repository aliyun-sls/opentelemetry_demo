package service

import (
	"context"
	"github.com/gin-gonic/gin"
	"os"
	"sls-mall-go/common/model"
	"sls-mall-go/common/util"
	"sync"
)

type Action uint

const (
	Create Action = 1
)

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

type Logistic struct {
	model.Model
	OrderId          string          `json:"order_id" form:"order_id" gorm:"not null;index:OrderIdIdx;comment:'订单ID'"`
	UserId           uint            `json:"user_id" form:"user_id" gorm:"not null;index:UserIdIdx;comment:'用户ID'"`
	LogisticStatus   LogisticsStatus `json:"logistic_status" form:"logistic_status" gorm:"size:10;index;not null;comment:'物流状态'"`
	LogisticPosition string          `json:"logistic_position" form:"logistic_status" gorm:"size:1000;not null;comment:'物流位置'"`
}

func (Logistic) TableName() string {
	return "logistic" // 使用单数形式
}

type LogisticsMsg struct {
	UserId           uint            `json:"user_id"`
	OrderId          string          `json:"order_id"`
	Action           Action          `json:"action"`
	LogisticStatus   LogisticsStatus `json:"logistic_status" gorm:"size:10;index;not null;comment:'物流状态'"`
	LogisticPosition string          `json:"logistic_position" gorm:"size:1000;not null;comment:'物流位置'"`
}

type ListLogisticsRequest struct {
	Logistic
	util.Page
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
	OrderStatus        OrderStatus     `gorm:"column:order_status" json:"status"`
	LogisticStatus     LogisticsStatus `gorm:"column:logistic_status" form:"logistic_status" json:"logistic_status"`
	TotalPrice         float64         `gorm:"column:total_price" json:"total_price"`
	Logistics          []Logistic      `json:"logistics" gorm:"foreignKey:order_id;references:OrderId"`
}

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

// TableName 方法指定了 Order 结构体对应的表名为 `order`
func (Order) TableName() string {
	return "order" // 使用单数形式
}

func AddLogistic(c *gin.Context) {
	ctx := c.Request.Context()
	var logistic Logistic
	err := c.BindJSON(&logistic)
	if err != nil {
		util.Status400(c, err)
		return
	}

	//如果状态不是 Shipping 则调用订单服务更新订单状态
	if logistic.LogisticStatus != Shipping {
		req := map[string]interface{}{
			"order_id":        logistic.OrderId,
			"logistic_status": logistic.LogisticStatus,
		}
		err = ServiceCallPost(ctx, os.Getenv("OrderHost"), "/order/LogisticStatusUpdate", req, &util.Result{})
		if err != nil {
			util.Status500(c, err)
			return
		}
	}

	err = util.MDB.WithContext(ctx).Save(&logistic).Error
	if err != nil {
		util.Status500(c, err)
		return
	}
	util.Status200(c, true)
}

// ListLogistics 列物流状态
func ListLogistics(c *gin.Context) {
	ctx := c.Request.Context()
	var listLogisticsRequest ListLogisticsRequest
	err := c.BindJSON(&listLogisticsRequest)
	if err != nil {
		util.Status400(c, err)
		return
	}
	var orders []Order
	if listLogisticsRequest.Limit == 0 {
		listLogisticsRequest.Limit = 10
	}
	err = util.MDB.WithContext(ctx).Model(Order{}).Preload("Logistics").
		Where(&listLogisticsRequest.Logistic).
		Limit(listLogisticsRequest.Limit).Offset(listLogisticsRequest.Offset).Find(&orders).Error
	if err != nil {
		util.Status500(c, err)
	}
	util.Status200(c, orders)
}

func CreateMsg(msg *LogisticsMsg) {
	if msg.OrderId != "" {
		background := context.Background()
		//创建订单成功后 发送支付订单流程
		var wg sync.WaitGroup
		wg.Add(1)
		go func(msg *LogisticsMsg) {
			defer wg.Done()
			ServiceCallPost(background, os.Getenv("LogisticHost"), "/logistic/Create", msg, &util.Result{})
		}(msg)
	}
}
