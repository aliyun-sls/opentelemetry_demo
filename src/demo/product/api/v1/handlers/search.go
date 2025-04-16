package handlers

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/trace"
	"math/rand"
	"sls-mall-go/common/config"
	"sls-mall-go/common/model"
	"sls-mall-go/common/util"
)

// SearchSimple 查询
func SearchSimple(c *gin.Context) {

	ctx := c.Request.Context()
	//  主动注入错误
	if rand.Int()%5 == 0 {
		err := errors.New("查询商品超时，请稍后再试")
		util.LogErr(ctx, c, err, "product SearchSimple fail")
		util.Status500(c, err)
		return
	}
	var products model.Product
	_ = c.BindJSON(products)
	var ps []model.Product
	util.MDB.WithContext(ctx).
		Where("products_cate like ?", products.ProductsCate).
		Where("products_name like ?", products.ProductsName).
		Where("products_desc like ?", products.ProductsDesc).
		Find(&ps)
	util.Status200(c, ps)
}

func SearchByContent(c *gin.Context) {
	ctx := c.Request.Context()
	//  主动注入错误
	if rand.Int()%5 == 0 {
		err := errors.New("获取商品超时，请稍后再试")
		util.LogErr(ctx, c, err, "product SearchByContent fail")
		util.Status500(c, err)
		return
	}

	var searchRequest model.SearchRequest
	err := c.BindQuery(&searchRequest)
	if err != nil {
		util.Status400(c, err)
		return
	}
	var buf bytes.Buffer
	//query := map[string]interface{}{
	//	"query": map[string]interface{}{
	//		"multi_match": map[string]interface{}{
	//			"query": info,
	//			//"fields": []string{"products_name", "products_desc"},
	//			"type": "best_fields",
	//		},
	//	},
	//}
	conditions := []map[string]interface{}{
		{"match": map[string]interface{}{"products_status": model.Shelve}},
	}
	if searchRequest.Info != "" {
		conditions = append(conditions, map[string]interface{}{
			"match": map[string]interface{}{"products_name": searchRequest.Info},
		})
	}
	if searchRequest.CategoryId != 0 {
		conditions = append(conditions, map[string]interface{}{
			"match": map[string]interface{}{"products_cate": searchRequest.CategoryId},
		})
	}
	query := map[string]interface{}{
		"query": map[string]interface{}{
			"bool": map[string]interface{}{
				"must": conditions,
			},
		},
	}
	if err := json.NewEncoder(&buf).Encode(query); err != nil {
		fmt.Printf("Error encoding query: %s \n", err)
		return
	}

	tracer := otel.Tracer("elasticsearch-t")
	ctx1, finish1 := tracer.Start(ctx, "Search", trace.WithSpanKind(trace.SpanKindClient))
	span := trace.SpanFromContext(ctx1)
	span.SetAttributes(
		attribute.KeyValue{
			Key:   "db.system",
			Value: attribute.StringValue("elasticsearch")},
		attribute.KeyValue{
			Key:   "db.index",
			Value: attribute.StringValue(config.EsIndex)},
		attribute.KeyValue{
			Key:   "db.query",
			Value: attribute.StringValue(searchRequest.Info)},
	)
	if searchRequest.Limit == 0 {
		searchRequest.Limit = 10
	}
	res, err := util.ESClient.Search(
		util.ESClient.Search.WithContext(ctx),
		util.ESClient.Search.WithIndex(config.EsIndex),
		util.ESClient.Search.WithSize(searchRequest.Limit),
		util.ESClient.Search.WithFrom(searchRequest.Offset),
		util.ESClient.Search.WithBody(&buf),
		util.ESClient.Search.WithTrackTotalHits(true),
		util.ESClient.Search.WithPretty(),
	)
	finish1.End()
	if err != nil {
		fmt.Printf("Error getting response: %s \n", err)
		return
	}
	defer res.Body.Close()

	if res.IsError() {
		var e map[string]interface{}
		if err = json.NewDecoder(res.Body).Decode(&e); err != nil {
			fmt.Printf("Error parsing the response body: %s \n", err)
			return
		} else {
			// Print the response status and error information.
			fmt.Printf("[%s] %s: %s \n",
				res.Status(),
				e["error"].(map[string]interface{})["type"],
				e["error"].(map[string]interface{})["reason"],
			)
		}
	}
	var r model.EsResult
	if err = json.NewDecoder(res.Body).Decode(&r); err != nil {
		fmt.Printf("Error parsing the response body: %s \n", err)
	}
	// Print the response status, number of results, and request duration.
	//log.Printf(
	//	"[%s] %d hits; took: %dms",
	//	res.Status(),
	//	int(r["hits"].(map[string]interface{})["total"].(map[string]interface{})["value"].(float64)),
	//	int(r["took"].(float64)),
	//)
	//bts, _ := json.Marshal(&r)
	//c.String(http.StatusOK, string(bts))
	var productsList []model.Product
	for _, hit := range r.Hits.Hits {
		productsList = append(productsList, hit.Source)
	}
	util.Status200(c, productsList)
}
