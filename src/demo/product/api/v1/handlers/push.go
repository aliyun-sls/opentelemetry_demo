package handlers

import (
	"bytes"
	"context"
	"fmt"
	"log"
	"sls-mall-go/common/model"

	"github.com/aliyun/alibabacloud-oss-go-sdk-v2/oss"
	"github.com/aliyun/alibabacloud-oss-go-sdk-v2/oss/credentials"
	openapicred "github.com/aliyun/credentials-go/credentials"
)

func PushOSS(products model.Product) {
	// 请根据实际要求设置region，以实例华东1（杭州）为例，regionID为cn-hangzhou
	region := "cn-heyuan"

	config := new(openapicred.Config).
		// 指定Credential类型，固定值为ecs_ram_role
		SetType("ecs_ram_role")
	// （可选项）指定角色名称。如果不指定，OSS会自动获取角色。强烈建议指定角色名称，以降低请求次数

	arnCredential, gerr := openapicred.NewCredential(config)
	provider := credentials.CredentialsProviderFunc(func(ctx context.Context) (credentials.Credentials, error) {
		if gerr != nil {
			return credentials.Credentials{}, gerr
		}
		cred, err := arnCredential.GetCredential()
		if err != nil {
			return credentials.Credentials{}, err
		}
		return credentials.Credentials{
			AccessKeyID:     *cred.AccessKeyId,
			AccessKeySecret: *cred.AccessKeySecret,
			SecurityToken:   *cred.SecurityToken,
		}, nil
	})

	// 加载默认配置并设置凭证提供者和region
	cfg := oss.LoadDefaultConfig().
		WithCredentialsProvider(provider).
		WithRegion(region)

	// 创建OSS客户端
	client := oss.NewClient(cfg)
	log.Printf("ossclient: %v", client)

	// 检查 products 结构体的内容
	log.Printf("products: %+v", products)
	if products.ProductsName == "" {
		log.Printf("products.ProductsName is empty")
	}
	if products.ProductBasicType.ProductsPic == nil {
		log.Printf("products.ProductBasicType.ProductsPic is nil")
	}

	// 获取图片数据
	imageData := products.ProductBasicType.ProductsPic
	fmt.Println("imageData: ", imageData, "products.ProductsName: ", products.ProductsName)
	// 创建上传对象的请求
	request := &oss.PutObjectRequest{
		Bucket: oss.Ptr("o11y-demo-cn-heyuan"),           // 存储空间名称
		Key:    oss.Ptr("test/" + products.ProductsName), // 对象名称，使用文件名作为对象名称
		Body:   bytes.NewReader(imageData),               // 要上传的图片数据
	}

	// 发送上传对象的请求
	result, err := client.PutObject(context.TODO(), request)
	if err != nil {
		log.Fatalf("failed to put object %v", err)
	}

	// 打印上传对象的结果
	log.Printf("put object result:%#v\n", result)
}
