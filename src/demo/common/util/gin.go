package util

import (
	"github.com/gin-gonic/gin"
	"go.opentelemetry.io/contrib/instrumentation/github.com/gin-gonic/gin/otelgin"
	"net/http"
)

func InitGin() *gin.Engine {
	r := gin.New()
	r.Use(gin.Recovery())
	r.Use(otelgin.Middleware("my-server"))
	r.Use(LoggerMiddleware())
	return r
}

func Status200(c *gin.Context, data interface{}) {
	c.JSON(http.StatusOK, Result{
		Code:    http.StatusOK,
		Message: "ok",
		Data:    data,
	})
}

func Status400(c *gin.Context, err error) {
	//log.Println(err)
	c.JSON(http.StatusBadRequest, Result{
		Code:    http.StatusBadRequest,
		Message: err.Error(),
		Data:    err,
	})
}

func Status500(c *gin.Context, err error) {
	//log.Println(err)
	c.JSON(http.StatusInternalServerError, Result{
		Code:    http.StatusInternalServerError,
		Message: err.Error(),
		Data:    err,
	})
}
