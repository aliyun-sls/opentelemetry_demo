package main

import (
	"github.com/gin-gonic/gin"
	"log"
	"runtime"
	"time"
)

func ExceptionLoggerMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		defer func() {
			if err := recover(); err != nil {
				// 获取当前时间戳
				timestamp := time.Now().Format("2006-01-02 15:04:05")
				// 获取调用栈信息
				stack := make([]byte, 4<<10) // 4KB buffer
				n := runtime.Stack(stack, false)
				log.Printf("[PANIC] %s\nError: %v\nStack: %s", timestamp, err, stack[:n])
			}
		}()
		c.Next()
	}
}
