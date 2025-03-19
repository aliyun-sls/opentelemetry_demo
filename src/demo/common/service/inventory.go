package service

import (
	"context"
	"sls-mall-go/common/model"
	"sls-mall-go/common/util"
	"strconv"
)

// GetInventory 加入购物车获取库存
func GetInventory(ctx context.Context, productsId int) (*model.Inventory, error) {
	data := map[string]interface{}{}
	data["products_id"] = strconv.Itoa(productsId)
	inventory := &model.Inventory{}
	res := &util.Result{Data: inventory}
	err := serviceCallGet(ctx, "inventory:8181", "/api/v1/inventory/get_inventory", data, res)
	if err != nil {
		return inventory, err
	}
	return inventory, nil
}

// CreateInventory 初始化库存
func CreateInventory(ctx context.Context, inventory model.Inventory) error {
	return serviceCallPost(ctx, "inventory:8181", "/api/v1/inventory/create_inventory", inventory, &util.Result{})
}

// Predecrement 预减库存
func Predecrement(ctx context.Context, productsId uint, productsNum int) error {
	m := map[string]interface{}{
		"products_id":   productsId,
		"inventory_num": productsNum,
	}
	return serviceCallGet(ctx, "inventory:8181", "/api/v1/inventory/predecrement", m, &util.Result{})
}

// Decrement 支付减库存
func Decrement(ctx context.Context, productsId uint, productsNum int) error {
	m := map[string]interface{}{
		"products_id":   productsId,
		"inventory_num": productsNum,
	}
	return serviceCallGet(ctx, "inventory:8181", "/api/v1/inventory/decrement", m, &util.Result{})
}
