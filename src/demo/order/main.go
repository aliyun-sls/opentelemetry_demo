package main

import (
	"order/routers"
	"order/service"
	"sls-mall-go/common/util"
)

func main() {

	util.InitDB()
	service.InitFeatureFlag()
	r := util.InitGin()
	routers.Init(r)
	err := r.Run(":8080")
	util.Chk(err)
}
