package model

import (
	"sls-mall-go/common/util"
)

// Logistics 物流
type Logistics struct {
	Model
	OrderIdType
	UserIdType
	LogisticsStatus   LogisticsStatus `json:"logistics_status" gorm:"size:10;index;not null;comment:'物流状态'"`
	LogisticsPosition string          `json:"logistics_position" gorm:"size:1000;not null;comment:'物流位置'"`
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

type ListLogisticsRequest struct {
	Logistics
	util.Page
}
