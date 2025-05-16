package main

import (
	"pay-center/routers"
	"pay-center/service"
	"sls-mall-go/common/util"
)

func main() {
	util.InitDB()
	r := util.InitGin()
	service.InitFeatureFlag()
	routers.Init(r)
	err := r.Run(":8080")
	util.Chk(err)
}
