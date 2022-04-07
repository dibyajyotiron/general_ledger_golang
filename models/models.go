package models

import (
	"database/sql"
	"fmt"
	"log"
	"os"
	"time"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"

	"general_ledger_golang/pkg/config"
)

var db *gorm.DB
var sqlDB *sql.DB

type Model struct {
	Id        uint64    `gorm:"primaryKey;autoIncrement;" json:"id"`
	CreatedAt time.Time `gorm:"autoCreateTime;column:createdAt" json:"createdAt"`
	UpdatedAt time.Time `gorm:"autoUpdateTime;column:updatedAt" json:"updatedAt"`
}

// Setup initializes the database instance
func Setup() {
	var err error
	dsn := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=%s TimeZone=Asia/Shanghai",
		config.DatabaseSetting.Host,
		config.DatabaseSetting.Port,
		config.DatabaseSetting.User,
		config.DatabaseSetting.Password,
		config.DatabaseSetting.Name,
		config.DatabaseSetting.SSLMode,
	)
	conf := &gorm.Config{}

	if os.Getenv("APP_ENV") != "prod" {
		conf.Logger = logger.Default.LogMode(logger.Info)
	}

	db, err = gorm.Open(postgres.Open(dsn), conf)
	if err != nil {
		log.Fatalf("models.Setup err: %v", err)
	}

	if err = db.AutoMigrate(&Book{}, &Operation{}, &Posting{}, &BookBalance{}); err != nil {
		log.Fatalf("Automigration failed, error: %+v", err) // fataF is printf followed by panic
	}

	hasConstraint := db.Migrator().HasConstraint(&BookBalance{}, "non_negative_balance")

	if hasConstraint == false {
		// create this constraint only if it doesn't exist, if any modification needed,
		// drop and create the constraint as postgres doesn't have update constraint provision.
		err = db.Migrator().CreateConstraint(&BookBalance{}, "non_negative_balance")
		if err != nil {
			log.Fatalf("non_negative_balance check constraint add failed, error: %+v", err)
		}
	}

	sqlDB, err = db.DB()

	// SetMaxIdleConns sets the maximum number of connections in the idle connection pool.
	sqlDB.SetMaxIdleConns(10)

	sqlDB.SetMaxOpenConns(100)
}

// GetDB returns the database connections
func GetDB() (*gorm.DB, *sql.DB) {
	return db, sqlDB
}

// CloseDB closes database connection (unnecessary)
func CloseDB() {
	defer func(sqlDB *sql.DB) {
		err := sqlDB.Close()
		if err != nil {
			panic(err)
		}
	}(sqlDB)
}
