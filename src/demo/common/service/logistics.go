package service

import (
	"context"
	"sls-mall-go/common/model"
	"sls-mall-go/common/util"
)

func CreateLogistics(ctx context.Context, logistics model.Logistics) error {
	return serviceCallPost(ctx, "logistics:8183", "/api/v1/logistics/add", logistics, &util.Result{})
}
