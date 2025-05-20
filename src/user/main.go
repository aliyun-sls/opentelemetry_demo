package main

import (
	"os"

	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis/v8"
	"gorm.io/gorm"
)

type User struct {
	gorm.Model
	Username string `json:"username"`
	Password string `json:"password"`
	Role     int    `json:"role"` // 添加角色字段
}

const (
	ROLE_ADMIN = iota
	ROLE_USRR
)

// 初始化 Redis 客户端
var redisClient *redis.Client

// 初始化内置管理员用户
func initAdminUser() {
	var adminUser User
	result := db.Where("username = ?", "cms").First(&adminUser)
	if result.Error != nil {
		// 如果管理员用户不存在，则创建
		adminUser = User{
			Username: "cms",
			Password: "ali88",
			Role:     ROLE_ADMIN,
		}
		db.Create(&adminUser)
	}
}

func main() {
	// 初始化数据库
	InitDB()

	// 初始化内置管理员用户
	initAdminUser()

	// 初始化 Gin
	r := gin.Default()

	redisAddr := os.Getenv("REDIS_ENDPOINT")
	redisPassword := os.Getenv("REDIS_PASSWORD")
	redisClient = redis.NewClient(&redis.Options{
		Addr:     redisAddr,
		Password: redisPassword,
		DB:       0,
	})

	// 注册路由
	r.POST("/register", register)
	r.POST("/login", login)

	// 启动服务
	r.Run(":8080")
}
