package service

import (
	"context"
	"sls-mall-go/common/model"
	"sls-mall-go/common/util"
)

func DeleteCart(ctx context.Context, userId uint, cartId uint) error {
	m := map[string]interface{}{
		"ids":     []uint{cartId},
		"user_id": userId,
	}
	return serviceCallPost(ctx, "carts:8182", "/api/v1/carts/delete_products", m, &util.Result{})
}

func GetCart(ctx context.Context, id uint) (cart model.Cart, err error) {
	m := map[string]interface{}{
		"id": id,
	}
	err = serviceCallGet(ctx, "carts:8182", "/api/v1/carts/get_products", m,
		&util.Result{Data: &cart})
	return
}
