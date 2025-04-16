package main

import (
	"github.com/gin-gonic/gin"
)

func RegisterRoutes(r *gin.Engine) {
	// 商品相关路由
	productGroup := r.Group("/products")
	{
		productGroup.POST("/", createProduct)
		productGroup.PUT("/:id", updateProduct)
		productGroup.POST("/upload", uploadImage)
		productGroup.GET("/", listProducts)  // 新增：查询多个商品
		productGroup.GET("/:id", getProduct) // 新增：查询单个商品
	}
}
