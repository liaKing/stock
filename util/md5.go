package util

import (
	"crypto/md5"
	"encoding/hex"
)

// HashPassword 加密
func HashPassword(password string) string {
	hasher := md5.New()
	hasher.Write([]byte(password))
	hashedPassword := hex.EncodeToString(hasher.Sum(nil))
	return hashedPassword
}

// ComparePassword 解密
func ComparePassword(hashedPassword, password string) bool {
	return hashedPassword == HashPassword(password)
}
