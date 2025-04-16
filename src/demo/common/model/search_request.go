package model

import "sls-mall-go/common/util"

type SearchRequest struct {
	Info       string `json:"info" form:"info"`
	CategoryId int    `json:"category_id" form:"category_id"`
	util.Page
}
