package main

import (
	"sls-mall-go/common/config"
	"sls-mall-go/common/model"
	"sls-mall-go/common/util"
	apiv1 "sls-mall-go/product/api/v1"
)

func main() {
	util.InitInTimeZone()
	util.InitDB()
	err := util.MDB.AutoMigrate(&model.Product{})
	util.Chk(err)
	err = util.MDB.AutoMigrate(&model.Collect{})
	util.Chk(err)
	util.InitTrace()

	util.InitES()
	util.InitPyroscope(config.ServiceName)
	r := util.InitGin()
	apiv1.Routers(r)
	err = r.Run(":" + config.ServicePort)
	util.Chk(err)

}
