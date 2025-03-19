package main

import (
	"github.com/aliyun/aliyun-oss-go-sdk/oss"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"os"
)

var DB *gorm.DB
var OSSClient *oss.Client
var Bucket *oss.Bucket

func InitDB() error {
	user := os.Getenv("MYSQL_USER")
	password := os.Getenv("MYSQL_PASSWORD")
	host := os.Getenv("MYSQL_ENDPOINT")
	var err error
	dsn := user + ":" + password + "@tcp(" + host + ")/demo?charset=utf8mb4&parseTime=True&loc=Local"
	DB, err = gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		return err
	}
	// 自动迁移表结构
	DB.AutoMigrate(&Product{})

	return nil
}

var OSSEndpoint string
var OSSBucketName = "o11y-demo"

func InitOSS() error {
	var err error
	OSSEndpoint = os.Getenv("OSS_Endpoint")
	ossAk := os.Getenv("OSS_AccessKeyId")
	ossAs := os.Getenv("OSS_AccessKeySecret")
	OSSClient, err = oss.New(OSSEndpoint, ossAk, ossAs)
	if err != nil {
		return err
	}

	Bucket, err = OSSClient.Bucket(OSSBucketName)
	if err != nil {
		return err
	}

	return nil
}
