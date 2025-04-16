package model

import (
	"sls-mall-go/common/util"
)

// PostSale 售后
type PostSale struct {
	Model
	OrderIdType
	UserIdType
	CardIdType
	Reason       string        `json:"reason"`
	Status       Status        `json:"status"`
	Mode         Mode          `json:"mode"`
	RefundFee    float64       `json:"refund_fee"`
	Order        *Order        `json:"order" gorm:"foreignKey:OrderId;references:OrderId"`
	OrderDetails []OrderDetail `json:"order_details" gorm:"foreignKey:OrderId;references:OrderId"`
}

type ListPostSaleRequest struct {
	PostSale
	util.Page
}

type Status uint

const (
	Processing Status = 1
	Agree      Status = 2
	Refuse     Status = 3
)

type Mode uint

const (
	// Exchange 换货
	Exchange Mode = 1
	// Refund 退款
	Refund Mode = 2
	// Return 退货退款
	Return Mode = 3
)
