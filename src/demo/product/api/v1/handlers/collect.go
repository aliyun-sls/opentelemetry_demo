package handlers

import (
	"github.com/gin-gonic/gin"
	"sls-mall-go/common/model"
	"sls-mall-go/common/util"
)

func Collect(c *gin.Context) {
	ctx := c.Request.Context()
	var collectIdsRequest model.CollectIdsRequest
	err := c.BindJSON(&collectIdsRequest)
	if err != nil {
		util.Status400(c, err)
		return
	}
	var collects []model.Collect
	if len(collectIdsRequest.Ids) == 0 {
		util.Status200(c, collects)
		return
	}
	for _, productsId := range collectIdsRequest.Ids {
		product, err := getProduct(ctx, uint(productsId))
		if err != nil {
			util.Status500(c, err)
			return
		}
		collect := model.Collect{
			UserIdType:       model.UserIdType{UserId: collectIdsRequest.UserId},
			ProductsIdType:   model.ProductsIdType{ProductsId: uint(productsId)},
			ProductBasicType: product.ProductBasicType,
		}
		collects = append(collects, collect)
	}
	if len(collects) > 0 {
		err = util.MDB.WithContext(ctx).Create(&collects).Error
		if err != nil {
			util.Status500(c, err)
			return
		}
	}
	util.Status200(c, collects)
}

func CancelCollect(c *gin.Context) {
	ctx := c.Request.Context()
	var deleteCollectRequest model.CollectIdsRequest
	err := c.BindJSON(&deleteCollectRequest)
	if err != nil {
		util.Status400(c, err)
		return
	}
	err = util.MDB.WithContext(ctx).Where("user_id = ?", deleteCollectRequest.UserId).Delete(&model.Collect{}, deleteCollectRequest.Ids).Error
	if err != nil {
		util.Status500(c, err)
		return
	}
	util.Status200(c, true)
}

func ListCollect(c *gin.Context) {
	ctx := c.Request.Context()
	var listCollectRequest model.ListCollectRequest
	err := c.BindQuery(&listCollectRequest)
	if err != nil {
		util.Status400(c, err)
		return
	}
	if listCollectRequest.Limit == 0 {
		listCollectRequest.Limit = 10
	}
	var collects []model.CollectResponse
	err = util.MDB.WithContext(ctx).Model(&model.Collect{}).Select("*, true as `like` ").Where("user_id = ?", listCollectRequest.UserId).
		Limit(listCollectRequest.Limit).Offset(listCollectRequest.Offset).Scan(&collects).Error
	if err != nil {
		util.Status500(c, err)
		return
	}
	util.Status200(c, collects)
}
