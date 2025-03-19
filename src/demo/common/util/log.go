package util

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	"go.opentelemetry.io/otel/trace"
	"io/ioutil"
	"os"
	"runtime"
	"strings"
	"time"
)

func LoggerMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()
		traceParent := c.Request.Header.Get("Traceparent")
		traceArr := strings.Split(traceParent, "-")
		traceId := ""
		spanId := ""
		if len(traceArr) > 1 {
			traceId = traceArr[1]
			spanId = traceArr[2]
		}
		method := c.Request.Method
		url := c.Request.URL.String()
		requestHeaders := c.Request.Header
		reqBody, _ := ioutil.ReadAll(c.Request.Body)
		c.Request.Body = ioutil.NopCloser(bytes.NewBuffer(reqBody)) // 用来恢复原始数据
		// 创建一个缓冲，并将其设置为响应体的输出设备
		responseBodyBuffer := bytes.NewBuffer([]byte{})
		writer := responseBodyWrapper{
			ResponseWriter: c.Writer,
			responseBody:   responseBodyBuffer,
		}
		c.Writer = writer
		// 在请求结束时，获取响应体内容
		c.Next()
		responseHeaders := c.Writer.Header()
		responseBodyBytes := responseBodyBuffer.Bytes()
		requestHeadersJson, err := json.Marshal(requestHeaders)
		Chk(err)
		responseHeadersJson, err := json.Marshal(responseHeaders)
		Chk(err)
		reqs := string(reqBody)
		if reqs == "" {
			reqs = "{}"
		}
		resps := string(responseBodyBytes)
		if resps == "" {
			resps = "{}"
		}
		// 计算请求时长
		duration := time.Since(start).Milliseconds()
		_, file1, line1, _ := runtime.Caller(1)
		location := fmt.Sprintf(`%s:%d`, file1, line1)
		fmt.Printf(`{"Method":"%s","URL":"%s","Request_Headers":%s,"Request_Body":%s,"Response_Headers":%s,"Response_Body":%s,"trace_id":"%s","span_id":"%s","time_Consuming":"%d","location":"%s"}%s`, method, url, requestHeadersJson, reqs, responseHeadersJson, resps, traceId, spanId, duration, location, "\n")
	}
}

type responseBodyWrapper struct {
	gin.ResponseWriter
	responseBody *bytes.Buffer
}

func (w responseBodyWrapper) Write(b []byte) (int, error) {
	// 将响应数据写入缓冲区和原始响应体
	w.responseBody.Write(b)
	return w.ResponseWriter.Write(b)
}

func (w responseBodyWrapper) WriteString(s string) (int, error) {
	// 将响应数据写入缓冲区和原始响应体
	w.responseBody.WriteString(s)
	return w.ResponseWriter.WriteString(s)
}

func LogErr(ctx context.Context, c *gin.Context, err error, msg string) {
	spanContext := trace.SpanContextFromContext(ctx)
	_, file1, line1, _ := runtime.Caller(1)
	location := fmt.Sprintf(`%s:%d`, file1, line1)
	fmt.Fprintf(os.Stderr, `{"Method":"%s","URL":"%s","Message":"%s","Response_Body":"%s","trace_id":"%s","span_id":"%s","location":"%s"}%s`,
		c.Request.Method, c.Request.URL.String(), msg, err,
		spanContext.TraceID().String(),
		spanContext.SpanID().String(), location,
		"\n")
}
