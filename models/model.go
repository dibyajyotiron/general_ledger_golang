package models

import (
	"database/sql"
	"time"

	"gorm.io/gorm"

	"general_ledger_golang/pkg/database"
	"general_ledger_golang/pkg/logger"
)

var db *gorm.DB
var sqlDB *sql.DB

type Model struct {
	Id        uint64    `gorm:"primaryKey;autoIncrement;" json:"id"`
	CreatedAt time.Time `gorm:"autoCreateTime;column:createdAt" json:"createdAt"`
	UpdatedAt time.Time `gorm:"autoUpdateTime;column:updatedAt" json:"updatedAt"`
}

// Setup pulls in the created connection in to the models directory for future use
func Setup() {
	db, sqlDB = database.GetDB()
}

func GetDB() (*gorm.DB, *sql.DB) {
	logger.Logger.Infof("DB: %+v, SQLDB: %+v", db, sqlDB)
	return db, sqlDB
}
