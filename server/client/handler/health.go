package handler

import (
	"github.com/gin-gonic/gin"

	"life_journey/response"
	"life_journey/vault"
)

// Health 健康检查接口
func Health(c *gin.Context) {
	response.Success(c, gin.H{
		"status":     "ok",
		"vault_path": vault.Dir,
	})
}
