package model

import "sls-mall-go/common/util"

type Collect struct {
	Model
	UserIdType
	ProductsIdType
	ProductBasicType
}

type ListCollectRequest struct {
	Collect
	util.Page
}

type CollectIdsRequest struct {
	Ids []int `json:"ids"`
	UserIdType
}

type CollectResponse struct {
	Collect
	Like bool `json:"like"`
}
