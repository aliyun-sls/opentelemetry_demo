package handlers

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	otelhooks "github.com/open-feature/go-sdk-contrib/hooks/open-telemetry/pkg"
	flagd "github.com/open-feature/go-sdk-contrib/providers/flagd/pkg"
	"github.com/open-feature/go-sdk/openfeature"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/trace"
	"log"
	"net/http"
	"sls-mall-go/common/config"
	"sls-mall-go/common/model"
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

func processProduct(ctx context.Context, client *openfeature.Client) error {
	// 获取当前配置的 variant
	delayVariant, err := client.StringValue(
		ctx,
		"productDelay",
		"off", // 默认值
		openfeature.EvaluationContext{},
	)
	if err != nil {
		return fmt.Errorf("failed to get feature flag value: %w", err)
	}
	fmt.Println(delayVariant)
	// 根据 variant 确定延迟时间
	var delay time.Duration
	switch delayVariant {
	case "10sec":
		delay = 10 * time.Second
	case "5sec":
		delay = 5 * time.Second
	case "off":
		delay = 0
	default:
		delay = 0
	}

	// 应用延迟
	if delay > 0 {
		fmt.Printf("Applying product delay: %v\n", delay)
		time.Sleep(delay)
	} else {
		fmt.Println("Product delay is disabled")
	}

	// 这里放置你的产品处理逻辑
	fmt.Println("Processing product...")
	// 模拟产品处理
	time.Sleep(1 * time.Second)
	fmt.Println("Product processing complete")

	return nil
}

func PutProducts(c *gin.Context) {
	openfeature.AddHooks(otelhooks.NewTracesHook())
	err := openfeature.SetProvider(flagd.NewProvider())
	if err != nil {
		log.Fatal(err)
	}

	// 获取 productDelay flag 的值
	client := openfeature.NewClient("products")
	processProduct(context.Background(), client)
	image, err := c.FormFile("products_pic")
	if image != nil {
		err = PushOSS(image, image.Filename)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to upload image file"})
			return
		}
	}

	util.Status200(c, "Put Products complete.")
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
