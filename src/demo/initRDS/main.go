package main

import (
	"fmt"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/plugin/opentelemetry/tracing"
	"net/url"
	"sls-mall-go/common/model"
	"time"
)

func main() {
	InitMDB()
	err := MDB.AutoMigrate(&model.Product{}, &model.AdEntity{}, &model.MarketingEntity{}, &model.NotificationEntity{}, model.PromotionEntity{})
	if err != nil {
		panic(err)
	}
	initData()
}

func InitMDB() {
	loc, _ := time.LoadLocation("Asia/Shanghai")
	dsn := fmt.Sprintf("%s:%s@tcp(%s)/%s?charset=utf8mb4&parseTime=True&loc=",
		MysqlUser, MysqlPass, MysqlHost, "demo")
	dsn = dsn + url.QueryEscape(loc.String())

	db, err := gorm.Open(mysql.Open(dsn))
	if err != nil {
		panic(err)
	}
	err = db.Use(tracing.NewPlugin())
	if err != nil {
		panic(err)
	}

	MDB = db
}

func initData() {
	ads := model.AdEntity{
		Model:       gorm.Model{},
		ID:          1,
		RedirectUrl: "/product/2ZYFJ3GM2P",
		Text:        "/product/2ZYFJ3GM2P Roof Binoculars for sale. 50% off.",
	}
	err := MDB.Create(&ads).Error
	if err != nil {
		panic(err)
	}

	marketing := model.MarketingEntity{
		Model:       gorm.Model{},
		ID:          1,
		RedirectUrl: "/product/2ZYFJ3GM2P",
		Text:        "/product/2ZYFJ3GM2P Roof Binoculars for sale. 50% off.",
	}
	err = MDB.Create(&marketing).Error
	if err != nil {
		panic(err)
	}

	notification := model.NotificationEntity{
		Model:       gorm.Model{},
		ID:          1,
		RedirectUrl: "/product/2ZYFJ3GM2P",
		Text:        "/product/2ZYFJ3GM2P Roof Binoculars for sale. 50% off.",
	}
	err = MDB.Create(&notification).Error
	if err != nil {
		panic(err)
	}

	promotion := model.PromotionEntity{
		Model:       gorm.Model{},
		ID:          1,
		RedirectUrl: "/product/2ZYFJ3GM2P",
		Text:        "/product/2ZYFJ3GM2P Roof Binoculars for sale. 50% off.",
	}
	err = MDB.Create(&promotion).Error
	if err != nil {
		panic(err)
	}
}
