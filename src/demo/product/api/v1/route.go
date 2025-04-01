package v1

import (
	"github.com/gin-gonic/gin"
	"sls-mall-go/product/api/v1/handlers"
)

func Routers(r *gin.Engine) {

	api := r.Group("/api")
	{
		// v1
		v1 := api.Group("/v1")
		productsGroup := v1.Group("/products")
		{
			productsGroup.GET("/shelve", handlers.Shelve)
			productsGroup.GET("/unshelve", handlers.Unshelve)
			productsGroup.POST("/modify_products", handlers.ModifyProducts)
			productsGroup.GET("/get_products", handlers.GetProductsDetail)
			productsGroup.POST("/put_products", handlers.PutProducts)
			productsGroup.GET("/search_simple", handlers.SearchSimple)
			productsGroup.GET("/search_by_content", handlers.SearchByContent)

			productsGroup.POST("/collect", handlers.Collect)
			productsGroup.POST("/cancel_collect", handlers.CancelCollect)
			productsGroup.GET("/list_collect", handlers.ListCollect)
		}
	}
}
