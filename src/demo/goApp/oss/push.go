package OSS

import (
	"context"
	"github.com/gin-gonic/gin"
	"log"
	"os"
	"path/filepath"
	"time"

	"github.com/aliyun/alibabacloud-oss-go-sdk-v2/oss"
	"github.com/aliyun/alibabacloud-oss-go-sdk-v2/oss/credentials"
	openapicred "github.com/aliyun/credentials-go/credentials"
)

func PushOSS() {
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
	push(client)
	log.Printf("ossclient: %v", client)
}

func push(client *oss.Client) {
	// 获取指定目录下的所有文件
	path, _ := os.Getwd()
	files, err := os.ReadDir(path + "/tupian/")
	if err != nil {
		log.Fatalf("failed to read directory: %v", err)
	}

	// 遍历目录中的文件
	for _, file := range files {
		if !file.IsDir() {
			filePath := path + "/tupian/" + file.Name()
			file, err := os.Open(filePath)
			if err != nil {
				log.Fatalf("failed to open file: %v", err)
			}
			defer file.Close()

			// 创建上传对象的请求
			request := &oss.PutObjectRequest{
				Bucket: oss.Ptr("o11y-demo-cn-heyuan"), // 存储空间名称
				Key:    oss.Ptr("test/" + file.Name()), // 对象名称，使用文件名作为对象名称
				Body:   file,                           // 要上传的文件内容
			}

			// 发送上传对象的请求
			result, err := client.PutObject(context.TODO(), request)
			if err != nil {
				log.Fatalf("failed to put object %v", err)
			}

			// 打印上传对象的结果
			log.Printf("put object result:%#v\n", result)
		}
	}
}

// handler/upload.go
func UploadImage(c *gin.Context) {
	file, err := c.FormFile("image")
	if err != nil {
		c.JSON(400, gin.H{"error": "获取文件失败"})
		return
	}

	// 获取OSS Bucket
	bucket, err := config.OSSClient.Bucket("your-bucket-name")
	if err != nil {
		c.JSON(500, gin.H{"error": "OSS连接失败"})
		return
	}

	// 生成唯一文件名
	ext := filepath.Ext(file.Filename)
	filename := "products/" + time.Now().Format("20060102") + "/" + generateUUID() + ext

	// 打开文件
	src, err := file.Open()
	if err != nil {
		c.JSON(500, gin.H{"error": "文件读取失败"})
		return
	}
	defer src.Close()

	// 上传到OSS
	err = bucket.PutObject(filename, src)
	if err != nil {
		c.JSON(500, gin.H{"error": "上传OSS失败"})
		return
	}

	// 返回访问URL
	url := "https://your-bucket-name.oss-cn-hangzhou.aliyuncs.com/" + filename
	c.JSON(200, gin.H{"url": url})
}

func generateUUID() string {
	// 实现你自己的UUID生成逻辑
	return "unique-id"
}
