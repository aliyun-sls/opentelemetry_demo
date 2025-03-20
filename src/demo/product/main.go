package main

import (
	"sls-mall-go/common/util"
	apiv1 "sls-mall-go/product/api/v1"
)

func main() {
	util.InitInTimeZone()
	util.InitDB()
	//err := util.MDB.AutoMigrate(&model.Product{})
	//util.Chk(err)
	//err = util.MDB.AutoMigrate(&model.Collect{})
	//util.Chk(err)
	//util.InitTrace()

	//util.InitES()
	//util.InitPyroscope(config.ServiceName)
	r := util.InitGin()
	apiv1.Routers(r)
	//err := r.Run(":" + config.ServicePort)
	err := r.Run(":8080")
	util.Chk(err)

}
