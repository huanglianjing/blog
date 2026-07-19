package model

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/glebarez/sqlite"
	"gorm.io/gorm"
)

// DB 是全局的 GORM 数据库句柄，由 InitDB 初始化。
var DB *gorm.DB

// InitDB 打开（不存在则创建）dsn 指向的 SQLite 数据库，
// 通过 AutoMigrate 建好各表结构，赋值给全局 DB。
func InitDB(dsn string) error {
	if dir := filepath.Dir(dsn); dir != "." {
		if err := os.MkdirAll(dir, 0755); err != nil {
			return fmt.Errorf("create db dir %q: %w", dir, err)
		}
	}

	db, err := gorm.Open(sqlite.Open(dsn), &gorm.Config{})
	if err != nil {
		return fmt.Errorf("open sqlite %q: %w", dsn, err)
	}

	DB = db
	if err := migrate(); err != nil {
		return fmt.Errorf("migrate: %w", err)
	}
	return nil
}

// migrate 根据各模型的结构体标签自动建表 / 补列 / 建索引。
func migrate() error {
	return DB.AutoMigrate(
		&Article{},
		&Category{},
		&Tag{},
		&ArticleTag{},
	)
}
