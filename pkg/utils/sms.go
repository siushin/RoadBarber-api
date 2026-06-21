package utils

import (
	"crypto/rand"
	"math/big"
)

// GenerateCode 生成 6 位短信验证码（使用 crypto/rand，比 math/rand 更安全）
func GenerateCode() string {
	const digits = "0123456789"
	code := make([]byte, 6)
	max := big.NewInt(int64(len(digits)))
	for i := range code {
		n, err := rand.Int(rand.Reader, max)
		if err != nil {
			code[i] = '0'
			continue
		}
		code[i] = digits[n.Int64()]
	}
	return string(code)
}