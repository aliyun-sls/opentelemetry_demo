package main

import (
	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis/v8"
	"gorm.io/gorm"
	"os"
)

type User struct {
	gorm.Model
	Username string `json:"username"`
	Password string `json:"password"`
}

// 初始化 Redis 客户端
var redisClient *redis.Client

func main() {
	// 初始化数据库
	InitDB()

	// 初始化 Gin
	r := gin.Default()

	redisAddr := os.Getenv("REDIS_ADDR")
	redisClient = redis.NewClient(&redis.Options{
		Addr:     redisAddr, // Redis 服务器地址
		Password: "",        // 密码，如果没有设置密码则为空
		DB:       0,         // 默认数据库
	})

	// 注册路由
	r.POST("/register", register)
	r.POST("/login", login)

	// 启动服务
	r.Run(":8080")
}
