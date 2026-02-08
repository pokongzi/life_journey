package otp

import (
	"fmt"
	"testing"
	"time"
)

func TestGenerateOTP(t *testing.T) {
	secretKey := "17B588035513210A"
	// keyBytes := []byte(secretKey + "dimpt@@.com")
	const beijingLocation = "Asia/Shanghai" // 或者使用 "PRC"

	loc, err := time.LoadLocation(beijingLocation)
	if err != nil {
		fmt.Println("无法加载时区:", err)
		return
	}

	otp, remainingSeconds := generateOTPFixed(secretKey, time.Date(2025, 7, 25, 14, 0, 0, 0, loc))
	fmt.Println(otp, remainingSeconds)

	// otp, remainingSeconds = generateOTP1(secretKey)
	// fmt.Println(otp, remainingSeconds)
}

// func TestGenerateOTP1(t *testing.T) {
// 	now := time.Date(2025, 7, 25, 14, 0, 0, 0, time.UTC)
// 	otp, remainingSeconds := generateOTPFixed("MTdCNjg4MDM1NTEzMjEwQQ==", now)
// 	fmt.Println("OTP:", otp)
// 	fmt.Println("Remaining Seconds:", remainingSeconds)
// }
