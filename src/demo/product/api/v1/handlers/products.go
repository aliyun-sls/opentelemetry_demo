package handlers

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/trace"
	"math/rand"
	"net/http"
	"os"
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
	delay := c.Query("delay")
	if delay != "" {
		duration, err := time.ParseDuration(delay)
		if err != nil {
			c.JSON(400, gin.H{"error": "invalid delay format"})
			return
		}
		time.Sleep(duration)
		c.JSON(200, gin.H{"message": "delay finish"})
		return
	}
	//util.Status200(c, true)
}
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
	if os.Getenv("ISTROUBLE") == "true" {
		if rand.Int()%5 == 0 {
			time.Sleep(5 * time.Second)
		} else if rand.Int()%3 == 0 {
			time.Sleep(2 * time.Second)
		}
	}
	ctx := c.Request.Context()
	var product model.Product
	image, err := c.FormFile("products_pic")
	if image != nil {
		err = PushOSS(image, image.Filename)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to upload image file"})
			return
		}
	}
	product = model.Product{
		ID:             0,
		ProductsName:   c.PostForm("products_name"),
		ProductsPrice:  c.PostForm("price"),
		ProductsCate:   c.PostForm("products_cate"),
		ProductsDesc:   c.PostForm("products_desc"),
		BrandId:        c.PostForm("brand_id"),
		SellerId:       c.PostForm("seller_id"),
		ProductsStatus: 0,
		Inventory: &model.Inventory{
			InventoryId:      0,
			ProductsIdType:   model.ProductsIdType{},
			InventoryName:    c.PostForm("inventory_name"),
			InventoryNum:     c.PostForm("inventory_num"),
			InventoryUnit:    0,
			InventoryAddress: c.PostForm("inventory_address"),
		},
		ProductCategory: nil,
	}
	product.ProductsStatus = model.Shelve

	// 显式指定表名
	err = util.MDB.WithContext(ctx).Table("product").Model(&product).Create(&product).Error
	if err != nil {
		fmt.Println(err)
		util.Status500(c, err)
		return
	}

	err = esIndex(ctx, product)
	if err != nil {
		fmt.Println(err)
		util.Status500(c, err)
		return
	}

	// InventoryNum
	inventory := product.Inventory
	if inventory == nil {
		inventory = &model.Inventory{
			InventoryName: c.PostForm("inventory_name"),
			InventoryNum:  c.PostForm("inventory_num"),
		}
	}
	err = service.CreateInventory(ctx, *inventory)
	if err != nil {
		util.Status500(c, err)
		return
	}

	util.Status200(c, product)
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
		err = util.MDB.WithContext(ctx).Model(&model.Collect{}).Where(collect).Find(collect).Error
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

// esIndex 索引到Elasticsearch
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

// getProduct 获取产品详情
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
