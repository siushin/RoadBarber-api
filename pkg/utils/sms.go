package utils

import (
	"fmt"
	"math/rand"
	"time"
)

// GenerateCode 生成 6 位短信验证码
func GenerateCode() string {
	rand.Seed(time.Now().UnixNano())
	return fmt.Sprintf("%06d", rand.Intn(1000000))
}
