package handlers

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/trace"
	"sls-mall-go/common/config"
	"sls-mall-go/common/model"
	"sls-mall-go/common/service"
	"sls-mall-go/common/util"
	"strconv"
	"strings"
	"time"
)

// Shelve 上架
func Shelve(c *gin.Context) {
	ctx := c.Request.Context()
	var products model.Product
	err := c.BindQuery(&products)
	if err != nil {
		util.Status400(c, err)
		return
	}
	err = util.MDB.WithContext(ctx).Model(&products).Update("products_status", model.Shelve).Error
	if err != nil {
		util.Status500(c, err)
		return
	}
	err = esIndex(ctx, products)
	if err != nil {
		util.Status500(c, err)
		return
	}
	util.Status200(c, true)
}

// Unshelve 下架
func Unshelve(c *gin.Context) {
	ctx := c.Request.Context()
	var products model.Product
	err := c.BindQuery(&products)
	if err != nil {
		util.Status400(c, err)
		return
	}
	err = util.MDB.WithContext(ctx).Model(&products).Update("products_status", model.Unshelve).Error
	if err != nil {
		util.Status500(c, err)
		return
	}
	err = esIndex(ctx, products)
	if err != nil {
		util.Status500(c, err)
		return
	}
	util.Status200(c, true)
}

// ModifyProducts 修改
func ModifyProducts(c *gin.Context) {
	ctx := c.Request.Context()
	var products model.Product
	err := c.BindJSON(&products)
	if err != nil {
		util.Status400(c, err)
		return
	}
	err = util.MDB.WithContext(ctx).Model(&products).Updates(&products).Error
	if err != nil {
		util.Status500(c, err)
		return
	}
	err = esIndex(ctx, products)
	if err != nil {
		util.Status500(c, err)
		return
	}
	productKey := fmt.Sprintf("mall-products-%d", products.ID)
	util.RDB.Del(ctx, productKey)
	util.Status200(c, true)
}

// GetProducts 查询
//func GetProducts(c *gin.Context) {
//	ctx := c.Request.Context()
//	var products model.Product
//	err := c.BindQuery(&products)
//	if err != nil {
//		c.JSON(http.StatusBadRequest, util.Result{
//			Code:    http.StatusBadRequest,
//			Message: err.Error(),
//			Data:    nil,
//		})
//		return
//	}
//	var list []model.Product
//	err = util.MDB.WithContext(ctx).Where(&products).Find(&list).Error
//	if err != nil {
//		c.JSON(http.StatusInternalServerError, util.Result{
//			Code:    http.StatusInternalServerError,
//			Message: err.Error(),
//			Data:    nil,
//		})
//		return
//	}
//	c.JSON(http.StatusOK, util.Result{
//		Code:    http.StatusOK,
//		Message: "ok",
//		Data:    list,
//	})
//}

// PutProducts 保存
func PutProducts(c *gin.Context) {
	ctx := c.Request.Context()
	var products model.Product
	err := c.BindJSON(&products)
	if err != nil {
		fmt.Println(err)
		util.Status400(c, err)
		return
	}
	products.ProductsStatus = model.Shelve
	err = util.MDB.WithContext(ctx).Create(&products).Error
	if err != nil {
		fmt.Println(err)
		util.Status500(c, err)
		return
	}

	err = esIndex(ctx, products)
	if err != nil {
		fmt.Println(err)
		util.Status500(c, err)
		return
	}

	// InventoryNum
	inventory := products.Inventory
	if inventory == nil {
		inventory = &model.Inventory{}
		inventory.ProductsId = products.ID
		inventory.InventoryName = products.ProductsName
		inventory.InventoryNum = 100000
	}
	err = service.CreateInventory(ctx, *inventory)
	if err != nil {
		util.Status500(c, err)
		return
	}

	//err = util.RDB.HSet(ctx, "mall-products", products.ProductsId, string(bts)).Err()
	//if err != nil {
	//	fmt.Println(err)
	//	c.JSON(http.StatusInternalServerError, util.Result{
	//		Code:    http.StatusInternalServerError,
	//		Message: err.Error(),
	//		Data:    nil,
	//	})
	//	return
	//}

	util.Status200(c, products)

}

// GetProductsDetail 产品详细信息
func GetProductsDetail(c *gin.Context) {
	ctx := c.Request.Context()
	var getProductsDetailRequest model.GetProductsDetailRequest
	err := c.BindQuery(&getProductsDetailRequest)
	if err != nil || getProductsDetailRequest.ProductsId == 0 {
		util.Status400(c, err)
		return
	}
	var product model.GetProductsDetailResponse
	product.Product, err = getProduct(ctx, getProductsDetailRequest.ProductsId)
	collect := &model.Collect{
		UserIdType:     getProductsDetailRequest.UserIdType,
		ProductsIdType: getProductsDetailRequest.ProductsIdType,
	}
	if collect.UserId != 0 {
		err = util.MDB.WithContext(ctx).Where(collect).Find(collect).Error
		if err != nil {
			util.Status500(c, err)
			return
		}
		if collect.ID != 0 {
			product.Like = true
			product.CollectId = &collect.ID
		}
	}
	util.Status200(c, product)
}

func esIndex(ctx context.Context, products model.Product) error {
	var body strings.Builder
	docID := strconv.FormatInt(int64(products.ID), 10)

	body.Reset()
	bts, err := json.Marshal(products)
	if err != nil {
		return err
	}
	body.WriteString(string(bts))

	tracer := otel.Tracer("elasticsearch-t")
	ctx1, finish1 := tracer.Start(ctx, "Index", trace.WithSpanKind(trace.SpanKindClient))
	defer finish1.End()
	span := trace.SpanFromContext(ctx1)
	span.SetAttributes(
		attribute.KeyValue{
			Key:   "db.system",
			Value: attribute.StringValue("elasticsearch")},
		attribute.KeyValue{
			Key:   "db.index",
			Value: attribute.StringValue(config.EsIndex)},
		attribute.KeyValue{
			Key:   "db.docID",
			Value: attribute.StringValue(docID)},
		attribute.KeyValue{
			Key:   "db.body",
			Value: attribute.StringValue(body.String())},
	)
	span.AddEvent("event", trace.WithAttributes(
		attribute.String("docID", docID),
		attribute.String("body", body.String())),
	)
	_, err = util.ESClient.Index(
		config.EsIndex,
		strings.NewReader(body.String()),
		util.ESClient.Index.WithDocumentID(docID),
		util.ESClient.Index.WithRefresh("true"),
		util.ESClient.Index.WithPretty(),
		util.ESClient.Index.WithTimeout(100),
		util.ESClient.Index.WithContext(ctx),
	)
	if err != nil {
		return err
	}
	return nil
}

func getProduct(ctx context.Context, productsId uint) (product model.Product, err error) {
	productKey := fmt.Sprintf("mall-products-%d", productsId)
	productsJs, err := util.RDB.Get(ctx, productKey).Result()
	if err == nil {
		err = json.Unmarshal([]byte(productsJs), &product)
		if err != nil {
			return
		}
		err = util.RDB.Expire(ctx, productKey, time.Hour).Err()
		if err != nil {
			return
		}
		return
	}
	err = util.MDB.WithContext(ctx).Find(&product, productsId).Error
	if err != nil {
		return
	}
	productsJ, err := json.Marshal(product)
	if err != nil {
		return
	}
	err = util.RDB.SetEX(ctx, productKey, string(productsJ), time.Hour).Err()
	return
}
