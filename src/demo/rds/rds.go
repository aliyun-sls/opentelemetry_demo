// This file is auto-generated, don't edit it. Thanks.
package main

import (
	openapi "github.com/alibabacloud-go/darabonba-openapi/v2/client"
	openapiutil "github.com/alibabacloud-go/openapi-util/service"
	util "github.com/alibabacloud-go/tea-utils/v2/service"
	"github.com/alibabacloud-go/tea/tea"
	"os"
	"time"
)

// Description:
//
// 使用AK&SK初始化账号Client
//
// @return Client
//
// @throws Exception
func CreateClient() (_result *openapi.Client, _err error) {
	// 工程代码泄露可能会导致 AccessKey 泄露，并威胁账号下所有资源的安全性。以下代码示例仅供参考。
	// 建议使用更安全的 STS 方式，更多鉴权访问方式请参见：https://help.aliyun.com/document_detail/378661.html。
	config := &openapi.Config{
		// 必填，请确保代码运行环境设置了环境变量 ALIBABA_CLOUD_ACCESS_KEY_ID。
		AccessKeyId: tea.String(os.Getenv("ALIBABA_CLOUD_ACCESS_KEY_ID")),
		// 必填，请确保代码运行环境设置了环境变量 ALIBABA_CLOUD_ACCESS_KEY_SECRET。
		AccessKeySecret: tea.String(os.Getenv("ALIBABA_CLOUD_ACCESS_KEY_SECRET")),
	}
	// Endpoint 请参考 https://api.aliyun.com/product/Rds
	config.Endpoint = tea.String("rds.aliyuncs.com")
	_result = &openapi.Client{}
	_result, _err = openapi.NewClient(config)
	return _result, _err
}

// Description:
//
// # API 相关
//
// @param path - string Path parameters
//
// @return OpenApi.Params
func CreateApiInfo() (_result *openapi.Params) {
	params := &openapi.Params{
		// 接口名称
		Action: tea.String("DescribeDBInstanceHAConfig"),
		// 接口版本
		Version: tea.String("2014-08-15"),
		// 接口协议
		Protocol: tea.String("HTTPS"),
		// 接口 HTTP 方法
		Method:   tea.String("POST"),
		AuthType: tea.String("AK"),
		Style:    tea.String("RPC"),
		// 接口 PATH
		Pathname: tea.String("/"),
		// 接口请求体内容格式
		ReqBodyType: tea.String("json"),
		// 接口响应体内容格式
		BodyType: tea.String("json"),
	}
	_result = params
	return _result
}

func _main(args []*string) (_err error) {
	client, _err := CreateClient()
	if _err != nil {
		return _err
	}

	params := CreateApiInfo()
	// query params
	queries := map[string]interface{}{}
	queries["DBInstanceId"] = tea.String("rm-f8za7coa50f6j2l82")
	// runtime options
	runtime := &util.RuntimeOptions{}
	request := &openapi.OpenApiRequest{
		Query: openapiutil.Query(queries),
	}
	// 复制代码运行请自行打印 API 的返回值
	// 返回值实际为 Map 类型，可从 Map 中获得三类数据：响应体 body、响应头 headers、HTTP 返回的状态码 statusCode。
	_, _err = client.CallApi(params, request, runtime)
	if _err != nil {
		return _err
	}
	return _err
}

func main() {
	for {
		err := _main(tea.StringSlice(os.Args[1:]))
		if err != nil {
			panic(err)
		}
		time.Sleep(30 * time.Second)
	}

}
