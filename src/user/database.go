package main

import (
	"database/sql"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"os"
	"sync"
	"time"
)

var (
	db   *gorm.DB
	once sync.Once
)

func InitDB() *gorm.DB {
	once.Do(func() {
		user := os.Getenv("MYSQL_USER")
		password := os.Getenv("MYSQL_PASSWORD")
		host := os.Getenv("MYSQL_ENDPOINT")
		dsn := user + ":" + password + "@tcp(" + host + ")/dbname?charset=utf8mb4&parseTime=True&loc=Local"
		sqlDB, err := sql.Open("mysql", dsn)
		if err != nil {
			panic("failed to connect database")
		}

		// 配置连接池参数
		sqlDB.SetMaxOpenConns(100)                 // 最大打开连接数
		sqlDB.SetMaxIdleConns(10)                  // 最大空闲连接数
		sqlDB.SetConnMaxLifetime(time.Minute * 10) // 连接的最大生命周期

		db, err = gorm.Open(mysql.New(mysql.Config{
			Conn: sqlDB,
		}), &gorm.Config{})
		if err != nil {
			panic("failed to connect database")
		}

		err = db.AutoMigrate(&User{})
		if err != nil {
			return
		}
	})
	return db
}
