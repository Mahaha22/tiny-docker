package utils

import (
	"crypto/rand"
	"crypto/sha256"
	"encoding/hex"
)

func GetUniqueId() string {
	// 生成一个256位的随机数
	randomBytes := make([]byte, 256/8)
	if _, err := rand.Read(randomBytes); err != nil {
		panic(err)
	}
	// 计算SHA-256哈希值
	hashBytes := sha256.Sum256(randomBytes)
	// 取前12个字节作为容器ID
	containerId := hex.EncodeToString(hashBytes[:12])
	return containerId[:8]
}
