package main

import (
	"github.com/gin-gonic/gin"
	"log"
)

func main() {
	// 初始化数据库
	err := InitDB()
	if err != nil {
		log.Fatalf("failed to initialize database: %v", err)
	}

	// 初始化 OSS
	err = InitOSS()
	if err != nil {
		log.Fatalf("failed to initialize OSS: %v", err)
	}

	// 初始化 Gin
	r := gin.Default()

	// 注册路由
	RegisterRoutes(r)

	// 启动服务
	r.Run(":8080")
}
