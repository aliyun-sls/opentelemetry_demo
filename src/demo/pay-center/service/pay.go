package service

import (
	"context"
	"errors"
	"github.com/gin-gonic/gin"
	flagd "github.com/open-feature/go-sdk-contrib/providers/flagd/pkg"
	"github.com/open-feature/go-sdk/openfeature"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/codes"
	"go.opentelemetry.io/otel/trace"
	"log"
	"sls-mall-go/common/model"
	"sls-mall-go/common/util"
	"strconv"
)

var flagClient *openfeature.Client

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
	OrderDetails       []OrderDetail   `json:"order_details" gorm:"foreignKey:OrderId;references:OrderId"`
	TotalPrice         float64         `gorm:"column:total_price" json:"total_price"`
}

// TableName 方法指定了 Order 结构体对应的表名为 `order`
func (Order) TableName() string {
	return "order" // 使用单数形式
}

func PayOrder(c *gin.Context) {
	ctx := context.Background()
	var order Order
	err := c.BindJSON(&order)
	if err != nil {
		util.Status400(c, err)
		return
	}

	isFail := flagClient.Boolean(
		context.Background(),
		"PayServiceAbnormal",
		false,
		openfeature.EvaluationContext{},
	)
	if isFail {
		span := trace.SpanFromContext(ctx)
		span.AddEvent("pay.center -> fail")
		span.SetAttributes(
			attribute.String("app.user.id", strconv.FormatInt(order.UserId, 10)),
			attribute.String("app.order.id", order.OrderId),
		)

		err := errors.New("支付失败")
		// 标记Span为错误状态，并记录错误信息
		span.RecordError(err, trace.WithAttributes(
			attribute.String("app.user.id", strconv.FormatInt(order.UserId, 10)),
			attribute.String("app.order.id", order.OrderId),
			attribute.String("error.type", "pay.center -> create"),
			attribute.String("error.stack", "pay.center -> fail"),
		))
		span.SetStatus(codes.Error, err.Error())
		util.Status500(c, err)
		span.End()
		return
	}

	tx := util.MDB.WithContext(ctx).Model(&order).
		Where("user_id = ? AND order_id = ?", order.UserId, order.OrderId).Update("order_status", WaitForSending)

	if tx.Error != nil {
		util.Status500(c, tx.Error)
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
	flagClient = openfeature.NewClient("pay")
}
