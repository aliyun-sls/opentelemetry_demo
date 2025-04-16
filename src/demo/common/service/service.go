package service

import (
	"context"
	"errors"
	"fmt"
	"sls-mall-go/common/util"
)

func serviceCallPost(ctx context.Context, host string, path string, data interface{}, res *util.Result) error {
	client := &util.HttpClient{
		Host: host,
	}
	err := client.Post(ctx, path, data, res)
	if err != nil {
		return err
	}
	if res.Code != 200 {
		s := fmt.Sprintf("%d %s %v", res.Code, res.Message, res.Data)
		return errors.New(s)
	}
	return nil
}

func serviceCallGet(ctx context.Context, host string, path string, data map[string]interface{}, res *util.Result) error {
	client := &util.HttpClient{
		Host: host,
	}
	err := client.Get(ctx, path, data, res)
	if err != nil {
		return err
	}
	if res.Code != 200 {
		s := fmt.Sprintf("%d %s %v", res.Code, res.Message, res.Data)
		return errors.New(s)
	}
	return nil
}
