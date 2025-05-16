package routers

import (
	"github.com/gin-gonic/gin"
	"pay-center/service"
)

func Init(r *gin.Engine) {
	order := r.Group("/pay")
	{
		order.POST("/Create", service.PayOrder)
	}
}
