package oteldemo

import (
	"database/sql"
	"fmt"
	"log"
	"os"
	"time"

	_ "github.com/go-sql-driver/mysql"
)

func Mysql() {
	// 从环境变量中读取数据库连接信息
	dbUser := os.Getenv("MYSQL_USER")
	dbPassword := os.Getenv("MYSQL_PASSWORD")
	dbHost := os.Getenv("MYSQL_ENDPOINT")
	dbName := "demo"

	// 构建数据库连接字符串
	dsn := fmt.Sprintf("%s:%s@tcp(%s)/%s", dbUser, dbPassword, dbHost, dbName)

	// 打开数据库连接
	db, err := sql.Open("mysql", dsn)
	if err != nil {
		log.Fatalf("无法打开数据库连接: %v", err)
	}
	defer db.Close()

	// 测试数据库连接
	err = db.Ping()
	if err != nil {
		log.Fatalf("无法连接到数据库: %v", err)
	}

	// 每5秒查询一次数据库
	ticker := time.NewTicker(5 * time.Second)
	defer ticker.Stop()

	for range ticker.C {
		var result string
		err := db.QueryRow("SELECT id").Scan(&result)
		if err != nil {
			log.Printf("查询失败: %v", err)
			continue
		}
		log.Printf("查询结果: %s", result)
	}
}
