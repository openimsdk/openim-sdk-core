package db

import (
	"testing"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func Test_DataBase_migrateToVersion_returns_error_when_auto_migrate_fails(t *testing.T) {
	// Given
	conn, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	if err != nil {
		t.Fatalf("open database: %v", err)
	}
	sqlDB, err := conn.DB()
	if err != nil {
		t.Fatalf("get sql database: %v", err)
	}
	if err := sqlDB.Close(); err != nil {
		t.Fatalf("close sql database: %v", err)
	}
	database := &DataBase{conn: conn}

	// When
	err = database.migrateToVersion("3.8.0")

	// Then
	if err == nil {
		t.Fatal("expected auto-migration error")
	}
}
