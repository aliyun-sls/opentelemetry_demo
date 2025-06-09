package command

import (
	"os"
	"log"
	openapi "github.com/alibabacloud-go/darabonba-openapi/v2/client"
	ecs "github.com/alibabacloud-go/ecs-20140526/v4/client"
	cs "github.com/alibabacloud-go/cs-20151215/v5/client"
	"github.com/alibabacloud-go/tea/tea"
)

// InitAliyunClients 初始化阿里云客户端（只执行一次）
func InitAliyunClients() {
	AliyunOnce.Do(func() {
		accessKeyId := os.Getenv("ALIBABA_CLOUD_ACCESS_KEY_ID")
		accessKeySecret := os.Getenv("ALIBABA_CLOUD_ACCESS_KEY_SECRET")
		regionId := os.Getenv("ALIBABA_CLOUD_REGION_ID")
		
		if accessKeyId == "" || accessKeySecret == "" || regionId == "" {
			log.Fatal("阿里云认证信息缺失，请设置环境变量: ALIBABA_CLOUD_ACCESS_KEY_ID, ALIBABA_CLOUD_ACCESS_KEY_SECRET, ALIBABA_CLOUD_REGION_ID")
		}

		// 创建 ECS 客户端配置
		ecsConfig := &openapi.Config{
			AccessKeyId:     tea.String(accessKeyId),
			AccessKeySecret: tea.String(accessKeySecret),
			RegionId:        tea.String(regionId),
			Endpoint:        tea.String("ecs." + regionId + ".aliyuncs.com"),
		}

		// 创建 CS 客户端配置
		csConfig := &openapi.Config{
			AccessKeyId:     tea.String(accessKeyId),
			AccessKeySecret: tea.String(accessKeySecret),
			RegionId:        tea.String(regionId),
			Endpoint:        tea.String("cs." + regionId + ".aliyuncs.com"),
		}

		var err error
		// 初始化 ECS 客户端
		EcsClient, err = ecs.NewClient(ecsConfig)
		if err != nil {
			log.Fatalf("初始化 ECS 客户端失败: %v", err)
		}

		// 初始化 CS 客户端
		CsClient, err = cs.NewClient(csConfig)
		if err != nil {
			log.Fatalf("初始化 CS 客户端失败: %v", err)
		}

		log.Println("阿里云客户端初始化成功")
	})
} 