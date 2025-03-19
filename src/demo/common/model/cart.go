package model

import (
	"sls-mall-go/common/util"
)

// Cart 购物车
type Cart struct {
	Model
	UserIdType
	ProductsIdType
	ProductBasicType
	ProductsPrice float64 `json:"products_price" gorm:"decimal(18,2);not null;comment:'商品价格'"`
	ProductsNum   int     `json:"products_num" gorm:"not null;;comment:'数量'"`
}

type ListCartRequest struct {
	Cart
	util.Page
}

type DeleteCartRequest struct {
	Ids []int `json:"ids"`
	UserIdType
}

// ListCartsResponse 浏览响应
type ListCartsResponse struct {
	Cart
	Like      bool  `json:"like"`
	CollectId *uint `json:"collect_id,omitempty"`
}
