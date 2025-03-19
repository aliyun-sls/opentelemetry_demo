package service

import (
	"context"
	"sls-mall-go/common/model"
	"sls-mall-go/common/util"
)

func Refund(ctx context.Context, pay model.Pay) error {
	return serviceCallPost(ctx, "pay:8088", "/pay/refund", pay, &util.Result{})
}
