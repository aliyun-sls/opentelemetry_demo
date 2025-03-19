package util

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/go-redis/redis/extra/redisotel/v8"
	"github.com/go-redis/redis/v8"
	"go.opentelemetry.io/otel/trace"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"gorm.io/plugin/opentelemetry/tracing"
	"net/url"
	"sls-mall-go/common/config"
	"time"
)

var MDB *gorm.DB
var RDB *redis.Client

type MyDBLog struct{}

func (l *MyDBLog) LogMode(level logger.LogLevel) logger.Interface {
	newlogger := *l
	return &newlogger
}

func (l *MyDBLog) Info(ctx context.Context, msg string, data ...interface{}) {
	spanContex := trace.SpanContextFromContext(ctx)
	dataJson, err := json.Marshal(data)
	Chk(err)
	fmt.Printf(`{"level":"%s", "msg":"%s", "data": %s", "span_id":"%s", "trace_id":"%s"}
`, "Info", msg, dataJson, spanContex.SpanID(), spanContex.TraceID())
}
func (l *MyDBLog) Warn(ctx context.Context, msg string, data ...interface{}) {
	spanContex := trace.SpanContextFromContext(ctx)
	dataJson, err := json.Marshal(data)
	Chk(err)
	fmt.Printf(`{"level":"%s", "msg":"%s", "data": %s", "span_id":"%s", "trace_id":"%s"}
`, "Warn", msg, dataJson, spanContex.SpanID(), spanContex.TraceID())
}
func (l *MyDBLog) Error(ctx context.Context, msg string, data ...interface{}) {
	spanContex := trace.SpanContextFromContext(ctx)
	dataJson, err := json.Marshal(data)
	Chk(err)
	fmt.Printf(`{"level":"%s", "msg":"%s", "data": %s", "span_id":"%s", "trace_id":"%s"}
`, "Error", msg, dataJson, spanContex.SpanID(), spanContex.TraceID())
}
func (l *MyDBLog) Trace(ctx context.Context, begin time.Time, fc func() (sql string, rowsAffected int64), err error) {
	spanContex := trace.SpanContextFromContext(ctx)
	errJson, err := json.Marshal(err)
	Chk(err)
	sql, rowsAffected := fc()
	fmt.Printf(`{"level":"%s", "error":%s, "sql": "%s", "rowsAffected": %d, "span_id":"%s", "trace_id":"%s"}
`, "Trace", errJson, sql, rowsAffected, spanContex.SpanID(), spanContex.TraceID())
}

func InitMDB() {

	loc, _ := time.LoadLocation("Asia/Shanghai")
	//+url.QueryEscape(loc.String())
	//fmt.Println(url.QueryEscape(loc.String()))
	// "root:12345@tcp(127.0.0.1:3306)/go_db?charset=utf8mb4&parseTime=True"
	dsn := fmt.Sprintf("%s:%s@tcp(%s:3306)/%s?charset=utf8mb4&parseTime=True&loc=",
		config.MysqlUser, config.MysqlPass, config.MysqlHost, config.MysqlDB)
	//dsn := "root:root@tcp(127.0.0.1:3306)/my_victore?charset=utf8mb4&parseTime=True&loc=Asia%2FShanghai"
	//fmt.Println(dsn)
	dsn = dsn + url.QueryEscape(loc.String())
	var newLogger = &MyDBLog{}

	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{Logger: newLogger})
	Chk(err)
	err = db.Use(tracing.NewPlugin())
	Chk(err)

	//MDB = db.Debug()
	MDB = db
}

func InitRDB() {
	RDB = redis.NewClient(&redis.Options{
		Addr:     config.RedisAddr,
		Password: config.RedisPass,
		DB:       1,
	})
	RDB.AddHook(redisotel.NewTracingHook())
}

func InitDB() {
	InitRDB()
	InitMDB()
}
