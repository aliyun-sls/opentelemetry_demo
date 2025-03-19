package model

import (
	"sls-mall-go/common/util"
)

// ProductCategory 商品类目
type ProductCategory struct {
	ProductCategoryId uint   `json:"product_category_id" gorm:"primaryKey"`
	Name              string `json:"name" gorm:"index"`
	Description       string `json:"description"`
	Status            int    `json:"status"`
}

type ListProductCategoryRequest struct {
	ProductCategory
	util.Page
}
