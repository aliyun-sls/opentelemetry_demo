package main

import (
	"logistic-center/routers"
	"logistic-center/service"
	"sls-mall-go/common/util"
)

func main() {
	util.InitDB()
	r := util.InitGin()
	routers.Init(r)
	go service.Consumerkafka()
	err := r.Run(":8080")
	util.Chk(err)
}
