package util

import (
	"crypto/rand"
	"encoding/json"
	"fmt"
	othttptrace "go.opentelemetry.io/contrib/instrumentation/net/http/otelhttp"
	"log"
	"net/http"
	"time"
)

func GetRandomString(n int) string {
	randBytes := make([]byte, n/2)
	_, err := rand.Read(randBytes)
	Chk(err)
	return fmt.Sprintf("%x", randBytes)
}

func Chk(err error) {
	if err != nil {
		log.Fatal(err)
	}
}

func ChkLogin(bytes []byte) bool {
	res := &Result{}
	err := json.Unmarshal(bytes, res)
	if err != nil {
		log.Fatalln(err)
	}
	if res.Code != 200 {
		fmt.Println(res)
	}
	if res.Code == 401 {
		return false
	}
	return true
}

func ResChk(bytes []byte) {
	res := &Result{}
	err := json.Unmarshal(bytes, res)
	if err != nil {
		log.Fatalln(err)
	}
	if res.Code != 200 {
		fmt.Println(res)
	}
}

func NewOTELHTTPTransport(transport http.RoundTripper) http.RoundTripper {
	if transport == nil {
		transport = http.DefaultTransport
	}
	return othttptrace.NewTransport(
		transport,
	)
}

func SleepRandomS(max int) {
	s := GenNum(max)
	if s < 3 {
		s = 3
	}
	time.Sleep(time.Duration(s) * time.Second)
}

// 设置全局时区 东八区
func InitInTimeZone() {
	var cstZone = time.FixedZone("CST", 8*3600) // 东八
	time.Local = cstZone
}
