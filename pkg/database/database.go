package database

import (
	"database/sql"
	"fmt"
	"log"
	"os"
	"time"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	gLogger "gorm.io/gorm/logger"

	"general_ledger_golang/pkg/config"
	"general_ledger_golang/pkg/logger"
	"general_ledger_golang/pkg/util"
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
	cfg := *config.GetConfig()

	dsn := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=%s TimeZone=Asia/Shanghai",
		cfg.DatabaseSetting.Host,
		cfg.DatabaseSetting.Port,
		cfg.DatabaseSetting.User,
		cfg.DatabaseSetting.Password,
		cfg.DatabaseSetting.Name,
		cfg.DatabaseSetting.SSLMode,
	)
	conf := &gorm.Config{}

	// Log the queries if environment is not prod
	if !util.Includes(os.Getenv("APP_ENV"), []interface{}{"prod", "production", "release"}) {
		newLogger := gLogger.New(
			log.New(os.Stdout, "\r\n", log.LstdFlags), // io writer
			gLogger.Config{
				SlowThreshold:             time.Millisecond * time.Duration(100), // Slow SQL threshold
				LogLevel:                  gLogger.Warn,                          // Log level
				IgnoreRecordNotFoundError: false,                                 // Ignore ErrRecordNotFound error for logger
				Colorful:                  false,                                 // Disable color
			},
		)
		conf.Logger = newLogger
	}

	db, err = gorm.Open(postgres.Open(dsn), conf)
	if err != nil {
		logger.Logger.Fatalf("Gorm connection open err: %v", err)
	}

	sqlDB, err = db.DB()
	if err != nil {
		logger.Logger.Fatalf("Gorm sqlDB err: %v", err)
	}

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

//// GetCode Given a gorm response, it returns the sql error code.
//func GetCode(value *gorm.DB) string {
//	if value.Error != nil {
//		switch value.Dialector.Name() {
//		case "sqlite":
//			if err, ok := value.Error.(sqlite3.Error); ok {
//				return err.ExtendedCode.Error()
//			}
//		case "postgres":
//			if err, ok := value.Error.(*pgconn.PgError); ok {
//				return err.Code
//			}
//		}
//	}
//	return value.Error.Error()
//}
