package routers

import (
	"github.com/gin-gonic/gin"
	"order/service"
)

func Init(r *gin.Engine) {
	order := r.Group("/order")
	{
		order.POST("/Create", service.CreateOrder)
		order.POST("/Pay", service.PayOrder)
		order.POST("/list", service.ListOrderAndDetails)
		order.GET("/get", service.GetOrder)
		order.POST("/LogisticStatusUpdate", service.LogisticStatusUpdate)
		order.GET("/shipping", service.CreateShipping)
	}
}
