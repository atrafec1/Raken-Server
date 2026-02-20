package database

import (
	"fmt"
	"os"
	"path/filepath"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

var DB *gorm.DB

func InitDB(appName string) error {
	path, err := resolveDBPath(appName)
	if err != nil {
		return fmt.Errorf("failed to resolve database path: %w", err)
	}

	db, err := gorm.Open(sqlite.Open(path), &gorm.Config{})
	if err != nil {
		return fmt.Errorf("failed to open database: %w", err)
	}

	DB = db

	return migrateDB(db)
}

func resolveDBPath(appName string) (string, error) {
	dir, err := os.UserConfigDir()
	if err != nil {
		return "", fmt.Errorf("failed to get user config directory: %w", err)
	}
	appDir := filepath.Join(dir, appName)
	//0755 is a common permission for config directories, allowing the owner to read/write/execute and others to read/execute
	if err := os.MkdirAll(appDir, 0755); err != nil {
		return "", fmt.Errorf("failed to create app config directory: %w", err)
	}
	return filepath.Join(appDir, "app.db"), nil
}

func migrateDB(db *gorm.DB) error {
	return db.AutoMigrate(
		&Job{},
		&Material{},
		&MaterialLog{},
		&Employee{},
		&CostCode{},
	)
}

func InitTestDB() error {
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	if err != nil {
		return fmt.Errorf("failed to open in-memory database: %w", err)
	}
	DB = db
	return migrateDB(db)
}
