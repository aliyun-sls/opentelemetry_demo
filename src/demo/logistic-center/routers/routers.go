package routers

import (
	"github.com/gin-gonic/gin"
	service "logistic-center/service"
)

func Init(r *gin.Engine) {
	logistic := r.Group("/logistic")
	{
		logistic.POST("/Create", service.AddLogistic)
		logistic.POST("/list", service.ListLogistics)

	}
}
