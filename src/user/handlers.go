package main

import (
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"log"
	"net/http"
	"time"
)

func login(c *gin.Context) {
	var user User
	if err := c.ShouldBindJSON(&user); err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"error":   "Invalid request payload",
			"details": err.Error(),
		})
		return
	}

	var foundUser User
	if err := db.Where("username = ? AND password = ?", user.Username, user.Password).First(&foundUser).Error; err != nil {
		c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
			"error":   "Invalid credentials",
			"details": err.Error(),
		})
		return
	}

	uid := uuid.New().String()
	if err := redisClient.Set(c, uid, foundUser.ID, 24*time.Hour).Err(); err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
			"error":   "Failed to create session",
			"details": err.Error(),
		})
		return
	}

	c.SetCookie("sid", uid, 3600*24, "/", "", false, true)

	c.JSON(http.StatusOK, gin.H{"message": "Login successful", "sid": uid, "role": foundUser.Role}) // 返回角色信息
}

func logout(c *gin.Context) {
	log.Println("[INFO] Logout 接口被调用") // 👈 新增日志：记录接口被调用

	// 获取 cookie
	sessionID, err := c.Cookie("sid")
	if err != nil {
		log.Printf("[ERROR] 无法获取 session cookie: %v\n", err) // 👈 新增日志：记录 cookie 获取失败
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"error":   "No active session",
			"details": err.Error(),
		})
		return
	}

	// 删除 Redis 中的会话
	if err := redisClient.Del(c, sessionID).Err(); err != nil {
		log.Printf("[ERROR] 删除 Redis 会话失败: %v\n", err) // 👈 新增日志：记录 Redis 删除失败
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
			"error":   "Failed to delete session",
			"details": err.Error(),
		})
		return
	}

	// 清除浏览器 cookie
	c.SetCookie("sid", "", -1, "/", "", false, true)
	log.Println("[INFO] 用户已成功登出，session 已清除") // 👈 新增日志：记录登出完成

	c.JSON(http.StatusOK, gin.H{"message": "Logout successful"})
}
