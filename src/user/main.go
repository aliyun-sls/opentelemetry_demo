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
	Role     int    `json:"role"` // 添加角色字段
}

const (
	ROLE_ADMIN = iota
)

// 初始化 Redis 客户端
var redisClient *redis.Client

func main() {
	// 初始化数据库
	InitDB()

	// 初始化 Gin
	r := gin.Default()

	redisAddr := os.Getenv("REDIS_ADDR")
	redisPassword := os.Getenv("REDIS_PASSWORD")
	redisClient = redis.NewClient(&redis.Options{
		Addr:     redisAddr,
		Password: redisPassword,
		DB:       0,
	})

	// 注册路由
	//r.POST("/register", register)
	r.POST("/login", login)
	r.POST("/logout", logout)

	// 启动服务
	r.Run(":8080")
}
