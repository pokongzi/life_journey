package database

import (
	"log"
	"os"
	"path/filepath"

	"github.com/glebarez/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"

	"life_journey/model"
	"life_journey/vault"
)

// DB 全局数据库实例
var DB *gorm.DB

// Init 初始化 Vault 目录 + SQLite 数据库（元数据索引）
func Init() {
	// 1. 初始化 Vault（Markdown 文件库）
	vault.Init()

	// 2. 初始化 SQLite（元数据存储，与 Vault 同级目录）
	dataDir := getDataDir()
	if err := os.MkdirAll(dataDir, 0755); err != nil {
		log.Fatalf("无法创建数据目录: %v", err)
	}

	dbPath := filepath.Join(dataDir, "life_journey.db")

	var err error
	DB, err = gorm.Open(sqlite.Open(dbPath), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Warn),
	})
	if err != nil {
		log.Fatalf("无法连接数据库: %v", err)
	}

	// 自动迁移建表（新增 path / file_path 列会自动加上）
	err = DB.AutoMigrate(
		&model.Notebook{},
		&model.Note{},
		&model.Todo{},
	)
	if err != nil {
		log.Fatalf("数据库迁移失败: %v", err)
	}

	log.Printf("数据库初始化完成: %s", dbPath)
}

// getDataDir 获取数据存储目录（存放 SQLite 等配置文件）
func getDataDir() string {
	if dir, err := os.UserConfigDir(); err == nil {
		return filepath.Join(dir, "life_journey")
	}
	home, _ := os.UserHomeDir()
	return filepath.Join(home, ".life_journey")
}
