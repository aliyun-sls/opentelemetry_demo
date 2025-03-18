package main

import (
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type User struct {
	gorm.Model
	Username string `json:"username"`
	Password string `json:"password"`
}

func main() {
	// 初始化数据库
	err := InitDB()
	if err != nil {
		return
	}

	// 初始化 Gin
	r := gin.Default()

	// 注册路由
	r.POST("/register", register)
	r.POST("/login", login)

	// 启动服务
	r.Run(":8080")
}
