package handler

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"

	"life_journey/middleware"
	"life_journey/response"
)

// ---------- 请求/响应结构体 ----------

// LoginRequest 邮箱+密码登录请求
type LoginRequest struct {
	Email    string `json:"email" binding:"required"`
	Password string `json:"password" binding:"required"`
}

// LoginByCodeRequest 邮箱+验证码登录请求
type LoginByCodeRequest struct {
	Email string `json:"email" binding:"required"`
	Code  string `json:"code" binding:"required"`
}

// RegisterRequest 注册请求
type RegisterRequest struct {
	Email    string `json:"email" binding:"required"`
	Password string `json:"password" binding:"required"`
	Nickname string `json:"nickname"`
}

// UserInfo 用户信息
type UserInfo struct {
	ID       uint   `json:"id"`
	Email    string `json:"email"`
	Nickname string `json:"nickname"`
}

// ---------- 工具函数 ----------

func generateToken(userID uint, email string) (string, error) {
	claims := jwt.MapClaims{
		"user_id": userID,
		"email":   email,
		"exp":     time.Now().Add(7 * 24 * time.Hour).Unix(),
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(middleware.JWTSecret)
}

// ---------- 处理函数 ----------

// Login 邮箱+密码登录（Mock：任意邮箱密码均可登录）
func Login(c *gin.Context) {
	var req LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Fail(c, http.StatusBadRequest, "参数错误：需要 email 和 password")
		return
	}

	token, err := generateToken(1, req.Email)
	if err != nil {
		response.Fail(c, http.StatusInternalServerError, "生成 token 失败")
		return
	}

	response.Success(c, gin.H{
		"token": token,
		"user": UserInfo{
			ID:       1,
			Email:    req.Email,
			Nickname: "本地用户",
		},
	})
}

// LoginByCode 邮箱+验证码登录（Mock）
func LoginByCode(c *gin.Context) {
	var req LoginByCodeRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Fail(c, http.StatusBadRequest, "参数错误：需要 email 和 code")
		return
	}

	token, err := generateToken(1, req.Email)
	if err != nil {
		response.Fail(c, http.StatusInternalServerError, "生成 token 失败")
		return
	}

	response.Success(c, gin.H{
		"token": token,
		"user": UserInfo{
			ID:       1,
			Email:    req.Email,
			Nickname: "本地用户",
		},
	})
}

// Register 注册（Mock）
func Register(c *gin.Context) {
	var req RegisterRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Fail(c, http.StatusBadRequest, "参数错误：需要 email 和 password")
		return
	}

	token, err := generateToken(1, req.Email)
	if err != nil {
		response.Fail(c, http.StatusInternalServerError, "生成 token 失败")
		return
	}

	nickname := req.Nickname
	if nickname == "" {
		nickname = "本地用户"
	}

	response.Success(c, gin.H{
		"token": token,
		"user": UserInfo{
			ID:       1,
			Email:    req.Email,
			Nickname: nickname,
		},
	})
}

// GetMe 获取当前用户信息
func GetMe(c *gin.Context) {
	email, _ := c.Get("email")
	userID, _ := c.Get("user_id")

	// JWT MapClaims 中数字默认为 float64
	var uid uint
	switch v := userID.(type) {
	case float64:
		uid = uint(v)
	default:
		uid = 1
	}

	emailStr, _ := email.(string)

	response.Success(c, UserInfo{
		ID:       uid,
		Email:    emailStr,
		Nickname: "本地用户",
	})
}

// Logout 登出
func Logout(c *gin.Context) {
	response.Success(c, nil)
}
