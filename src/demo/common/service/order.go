package service

import (
	"context"
	"sls-mall-go/common/model"
	"sls-mall-go/common/util"
)

func GetOrder(ctx context.Context, orderId string, userId uint) (*model.Order, error) {
	order := &model.Order{}
	res := &util.Result{
		Data: order,
	}
	m := map[string]interface{}{
		"order_id": orderId,
		"user_id":  userId,
	}
	err := serviceCallGet(ctx, "order:8087", "/order/details", m, res)
	if err != nil {
		return order, err
	}
	return order, nil
}

func PayOrder(ctx context.Context, orderId string, userId int) error {
	m := map[string]interface{}{
		"order_id": orderId,
		"user_id":  userId,
	}
	return serviceCallGet(ctx, "order:8087", "/order/pay", m, &util.Result{})
}

func ChangeOrderLogisticsStatus(ctx context.Context, orderId string, status uint) error {
	m := map[string]interface{}{
		"order_id":         orderId,
		"logistics_status": status,
	}
	return serviceCallGet(ctx, "order:8087", "/order/logistics_status", m, &util.Result{})
}

func SetCommentId(ctx context.Context, orderId string, productsId uint, commentId uint) error {
	m := map[string]interface{}{
		"order_id":    orderId,
		"products_id": productsId,
		"CommentId":   commentId,
	}
	return serviceCallGet(ctx, "order:8087", "/order/set_comment_id", m, &util.Result{})
}
