package config

import (
	"os"
)

var MysqlHost string
var MysqlUser string
var MysqlPass string
var MysqlDB string

var RedisAddr string
var RedisPass string

var AccessKeyID string
var AccessKeySecret string
var OSSEndpoint string
var OSSBucketName string

var EsAddr string
var EsUser string
var EsPass string
var EsIndex string

var TraceProject string
var TraceEndpoint string
var TraceInstance string

var ServiceName string
var ServicePort string

func init() {
	AccessKeyID = os.Getenv("SECRET_AK")
	AccessKeySecret = os.Getenv("SECRET_SK")

	if AccessKeyID == "" {
		AccessKeyID = os.Getenv("AK")
	}
	if AccessKeySecret == "" {
		AccessKeySecret = os.Getenv("SK")
	}

	MysqlHost = os.Getenv("MYSQL_HOST")
	MysqlUser = os.Getenv("MYSQL_USER")
	MysqlPass = os.Getenv("MYSQL_PASS")
	MysqlDB = os.Getenv("MYSQL_DB")

	RedisAddr = os.Getenv("REDIS_ADDR")
	RedisPass = os.Getenv("REDIS_PASS")

	OSSEndpoint = os.Getenv("OSS_ENDPOINT")
	OSSBucketName = os.Getenv("OSS_BUCKET_NAME")

	// ES地址 逗号分隔
	EsAddr = os.Getenv("ES_ADDR")
	EsUser = os.Getenv("ES_USER")
	EsPass = os.Getenv("ES_PASS")
	EsIndex = os.Getenv("ES_INDEX")

	TraceProject = os.Getenv("TRACE_PROJECT")
	TraceEndpoint = os.Getenv("TRACE_ENDPOINT")
	TraceInstance = os.Getenv("TRACE_INSTANCE")

	ServiceName = os.Getenv("SERVICE_NAME")
	ServicePort = os.Getenv("SERVICE_PORT")

	/*---------------------------------------------------*/
	//MysqlHost = "mall-mysql"
	//MysqlUser = "root"
	//MysqlPass = "root"
	//MysqlDB = "my_victore"
	//
	OSSEndpoint = "oss-cn-hangzhou.aliyuncs.com"
	OSSBucketName = "sls-mall-oss"
	//
	//RedisAddr = "mall-redis"
	//
	//TraceProject = "k8s-log-cc47e4a01e07d43128b7733651a954428"
	//TraceEndpoint = "k8s-log-cc47e4a01e07d43128b7733651a954428.cn-hangzhou-intranet.log.aliyuncs.com:10010"
	//TraceInstance = "sls-mall"
	//
}
