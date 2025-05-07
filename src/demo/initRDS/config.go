package main

import (
	"gorm.io/gorm"
	"os"
)

var (
	MysqlUser = os.Getenv("MYSQL_USER")
	MysqlPass = os.Getenv("MYSQL_PASS")
	MysqlHost = os.Getenv("MYSQL_HOST")
)

var MDB *gorm.DB
