package model

import (
	"sls-mall-go/common/util"
)

// Card 付款方式
type Card struct {
	CardId uint   `json:"card_id" gorm:"primaryKey"`
	Name   string `json:"name" form:"name" gorm:"index"`
	Url    string `json:"url" form:"url"`
}

type ListCardRequest struct {
	Card
	util.Page
}

type CardIdType struct {
	PayType int `json:"pay_type" gorm:"not null;index:PayTypeIdx;comment:'支付方式'"`
}
