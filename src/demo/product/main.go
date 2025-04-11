package main

import (
	"github.com/gin-gonic/gin"
	"sls-mall-go/common/model"
	"sls-mall-go/common/util"
	apiv1 "sls-mall-go/product/api/v1"
)

func main() {
	util.InitInTimeZone()
	util.InitDB()
	// 修改为使用指针传递模型，并确保模型字段类型正确
	err := util.MDB.AutoMigrate(&model.Product{})
	util.Chk(err)
	//err := util.MDB.AutoMigrate(&model.Product{})
	//util.Chk(err)
	//err = util.MDB.AutoMigrate(&model.Collect{})
	//util.Chk(err)
	//util.InitTrace()

	//util.InitES()
	//util.InitPyroscope(config.ServiceName)
	/*r := util.InitGin()*/
	r := gin.Default()
	apiv1.Routers(r)
	err = r.Run(":8080")
	util.Chk(err)

}
