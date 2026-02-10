package main

import (
	"log"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"

	"life_journey/database"
	"life_journey/handler"
	"life_journey/middleware"
)

func main() {
	// 初始化数据库
	database.Init()

	// 创建 Gin 引擎
	r := gin.Default()

	// CORS 配置：允许前端 localhost:8080 跨域访问
	r.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"http://localhost:8080", "http://127.0.0.1:8080"},
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "PATCH", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Authorization", "Accept"},
		ExposeHeaders:    []string{"Content-Length", "Content-Disposition", "X-Original-Size", "X-Compressed-Size", "X-Compression-Ratio"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}))

	// ---------- 路由注册 ----------
	api := r.Group("/api")
	{
		// 健康检查（无需鉴权）
		api.GET("/health", handler.Health)

		// 认证接口（无需鉴权）
		auth := api.Group("/auth")
		{
			auth.POST("/login", handler.Login)
			auth.POST("/login-by-code", handler.LoginByCode)
			auth.POST("/register", handler.Register)
			auth.GET("/me", middleware.AuthRequired(), handler.GetMe)
			auth.POST("/logout", middleware.AuthRequired(), handler.Logout)
		}

		// 以下接口需要鉴权
		protected := api.Group("")
		protected.Use(middleware.AuthRequired())
		{
			// 笔记本 CRUD
			notebooks := protected.Group("/notebooks")
			{
				notebooks.GET("", handler.ListNotebooks)
				notebooks.POST("", handler.CreateNotebook)
				notebooks.PUT("/:id", handler.UpdateNotebook)
				notebooks.DELETE("/:id", handler.DeleteNotebook)
			}

			// 笔记 CRUD
			notes := protected.Group("/notes")
			{
				notes.GET("", handler.ListNotes)
				notes.GET("/:id", handler.GetNote)
				notes.POST("", handler.CreateNote)
				notes.PUT("/:id", handler.UpdateNote)
				notes.DELETE("/:id", handler.DeleteNote)
			}

			// 待办 CRUD
			todos := protected.Group("/todos")
			{
				todos.GET("", handler.ListTodos)
				todos.POST("", handler.CreateTodo)
				todos.PUT("/:id", handler.UpdateTodo)
				todos.DELETE("/:id", handler.DeleteTodo)
			}

			// 工具接口
			tools := protected.Group("/tools")
			{
				tools.POST("/image/compress", handler.CompressImage)
			}
		}
	}

	// 启动服务
	addr := "127.0.0.1:13245"
	log.Printf("Life Journey Go 本地后端启动: http://%s", addr)
	log.Printf("健康检查: http://%s/api/health", addr)
	if err := r.Run(addr); err != nil {
		log.Fatalf("服务启动失败: %v", err)
	}
}
