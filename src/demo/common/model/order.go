package model

import (
	"sls-mall-go/common/util"
)

type OrderStatus uint

// 下单-付款-发货(物流)-收货-完成-评价
//
//		  \
//	    取消
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

// Order 订单
type Order struct {
	//todo 待使用 购物车id 读写忽略该字段
	CartId uint `json:"cart_id" gorm:"-"`
	Model
	OrderIdType
	UserIdType
	CardIdType
	AddressId       int             `json:"address_id"`
	SumPrice        float64         `json:"sum_price"`
	Freight         float64         `json:"freight"`
	OrderStatus     OrderStatus     `json:"order_status" form:"order_status" gorm:"index"`
	OrderDetails    []OrderDetail   `json:"order_details" gorm:"foreignKey:OrderId;references:OrderId"`
	LogisticsStatus LogisticsStatus `json:"logistics_status" form:"logistics_status" gorm:"size:10;index;not null;comment:'物流状态'"`
	Logistics       []Logistics     `json:"logistics" gorm:"foreignKey:OrderId;references:OrderId"`
	Card            *Card           `json:"card" gorm:"foreignKey:CardId;references:PayType"`
}

// OrderDetail 订单明细
type OrderDetail struct {
	OrderDetailId uint `json:"order_detail_id" gorm:"primaryKey"`
	OrderIdType
	UserIdType
	ProductsIdType
	ProductBasicType
	ProductsNum   int      `json:"products_num"`
	ProductsPrice float64  `json:"products_price"`
	CommentId     *uint    `json:"comment_id" gorm:"index"`
	Comment       *Comment `json:"comment,omitempty"  gorm:"foreignKey:CommentId"`
}

type OrderIdType struct {
	OrderId string `json:"order_id" form:"order_id" gorm:"not null;index:OrderIdIdx;comment:'订单ID'"`
}

type ListOrderRequest struct {
	Order
	util.Page
}

type ListDetailsRequest struct {
	OrderDetail
	CommentStatus bool `json:"comment_status"`
	util.Page
}
