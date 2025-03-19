package service

import (
	"context"
	"sls-mall-go/common/model"
	"sls-mall-go/common/util"
)

// GetProduct 查询产品
func GetProduct(ctx context.Context, productsId uint) (*model.Product, error) {
	data := map[string]interface{}{}
	data["products_id"] = productsId
	product := &model.Product{}
	res := &util.Result{Data: product}
	err := serviceCallGet(ctx, "product:8085", "/api/v1/products/get_products", data, res)
	if err != nil {
		return product, err
	}
	return product, nil
}
