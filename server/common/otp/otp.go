package otp

import (
	"crypto/hmac"
	"crypto/sha512"
	"fmt"
	"math/big"
	"time"
)

func generateOTP1(secretKey string) (string, int) {
	secretKey = secretKey + "dimpt@@.com"

	// 获取当前时间戳，通常以秒为单位
	const beijingLocation = "Asia/Shanghai" // 或者使用 "PRC"

	loc, _ := time.LoadLocation(beijingLocation)
	currentTime := time.Date(2025, 7, 25, 14, 0, 0, 0, loc).Unix()
	fmt.Println(currentTime)
	timeInterval := currentTime / 60

	// 计算下一分钟的时间戳
	nextMinute := (timeInterval + 1) * 60

	// 计算当前时间到下一分钟的剩余秒数
	remainingSeconds := nextMinute - currentTime

	// 将时间间隔和secretKey合并为一个字节数组
	data := make([]byte, 8)
	for i := 3; i >= 0; i-- {
		data[i] = byte(timeInterval & 0xff)
		timeInterval >>= 8
	}
	fmt.Println("data", data)

	keyBytes := []byte(secretKey)
	hmacHash := hmac.New(sha512.New, keyBytes)
	hmacHash.Write(data)
	hash := hmacHash.Sum(nil)

	// 获取哈希值的最后4个字节
	offset := hash[len(hash)-1] & 0x0f
	hashPart := hash[offset : offset+4]
	fmt.Println("hashPart", hashPart)

	// Convert the extracted bytes to a big integer
	bigInt := new(big.Int).SetBytes(hashPart)

	// Mask off the most significant bit and get the last 31 bits
	maxBigInt := new(big.Int).Lsh(big.NewInt(1), 31)
	bigInt.And(bigInt, maxBigInt)

	// 生成6位OTP（取后6位的结果）
	otp := fmt.Sprintf("%06d", bigInt.Mod(bigInt, big.NewInt(1000000)).Int64())
	return otp, int(remainingSeconds)
}

// generateOTPFixed 使用正确的密钥处理方式（UTF-8 字节）
func generateOTPFixed(secretKey string, now time.Time) (string, int) {
	fullSecretKey := secretKey + "dimpt@@.com"

	// 将拼接后的字符串转换为字节数组
	keyBytes := []byte(fullSecretKey)

	currentTime := now.Unix()
	fmt.Println("currentTime", currentTime)
	timeInterval := uint32(currentTime / 60)
	fmt.Println("timeInterval", timeInterval)

	nextMinute := (timeInterval + 1) * 60
	remainingSeconds := int(nextMinute - uint32(currentTime))

	// 构造数据: 4 字节时间 + 密钥字节
	data := make([]byte, 4+len(keyBytes))
	data[0] = byte(timeInterval >> 24)
	data[1] = byte(timeInterval >> 16)
	data[2] = byte(timeInterval >> 8)
	data[3] = byte(timeInterval)
	fmt.Println("data-1", data)
	wordSize := 4 // 每个 word 是 32 位（4 字节）
	for i := 0; i < len(keyBytes); i += wordSize {
		word := uint32(0)
		for j := 0; j < wordSize && i+j < len(keyBytes); j++ {
			word |= uint32(keyBytes[i+j]) << uint32((wordSize-j-1)*8)
		}
		offset := 4 + i/wordSize*wordSize
		for j := 0; j < wordSize && offset+j < len(data); j++ {
			data[offset+j] = byte(word >> uint32((wordSize-j-1)*8))
		}
	}
	// copy(data[4:], keyBytes)
	fmt.Println("data-2", data)

	// 使用 keyBytes 作为 HMAC 密钥
	mac := hmac.New(sha512.New, keyBytes)
	mac.Write(data)
	hash := mac.Sum(nil)

	// 提取 hash[8:12] -> 第 3 个 uint32（即 words[2]）
	value := uint32(hash[8])<<24 | uint32(hash[9])<<16 | uint32(hash[10])<<8 | uint32(hash[11])

	// 取低 31 位
	value &= 0x7FFFFFFF

	// 生成 6 位 OTP
	otp := fmt.Sprintf("%06d", value%1000000)

	return otp, remainingSeconds
}
