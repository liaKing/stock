package util

import (
	"encoding/base64"
)

// Encode 加密
func Encode(id string) string {
	// 编码为Base64
	encoded := base64.StdEncoding.EncodeToString([]byte(id))
	return encoded
}

// Decode 解密
func Decode(id string) string {
	// 解码Base64
	decodedBytes, _ := base64.StdEncoding.DecodeString(id)

	decodedStr := string(decodedBytes)
	return decodedStr
}
