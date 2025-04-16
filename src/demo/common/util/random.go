package util

import (
	"github.com/rs/xid"
	"math/rand"
	"time"
)

func GenId() string {
	return xid.New().String()
}

func GenNum(max int) int {
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	g := r.Intn(max)
	if g == 0 {
		return 1
	}
	return g
}

// 获取0到60的随机数 * 60  0分-->60分钟
func GetRandomInt64() int64 {
	randNum := 60 * (int64)(rand.Intn(121))
	return randNum
}

// 获取1到5的随机数
func GetRandomInt() int {
	randNum := rand.Intn(5) + 1
	return randNum
}

// 根据参入获取随机数
func RandomIntByNum(num int) int {
	randNum := rand.Intn(num) + 1
	return randNum
}
