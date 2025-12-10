package main

import (
	"os"
	"syscall"

	"github.com/joho/godotenv"
	"github.com/thoas/go-funk"

	"general_ledger_golang/internal/config"
	"general_ledger_golang/internal/database"
	auto "general_ledger_golang/internal/database/migrations/auto"
	"general_ledger_golang/internal/logger"
	grpcserver "general_ledger_golang/internal/transport/grpc"
	"general_ledger_golang/internal/util"
	"general_ledger_golang/models"
)

func init() {
	// use dotenv if explicitly marked enabled, else use dotenv only in local.
	if os.Getenv("DOT_ENV") == "enable" || funk.ContainsString([]string{"local", "localhost"}, os.Getenv("APP_ENV")) {
		err := godotenv.Load()
		logger.Logger.Info(".env Loaded")
		if err != nil {
			logger.Logger.Fatalf("Couldn't load .env, error: %+v", err)
		}
	}
	config.Setup("./internal/config/")
	database.Setup()
	models.Setup()

	// Auto migration may not be the way to prod. It ideally should be disabled in prod.
	// Migration should only run from cli in production as that can lead to serious locking of rows in case of alters.
	if os.Getenv("AUTO_MIGRATE") == "enable" || funk.ContainsString([]string{"local", "localhost"}, os.Getenv("APP_ENV")) {
		auto.Migrate()
	}

	logger.Setup()
	util.Setup()
}

// In-case, only grpc server is needed, http is not required, start this.
func main() {
	conf := config.GetConfig()

	grpcserver.RegisterGrpcServer(conf.ServerSetting.GrpcPort)

	logger.Logger.Infof("Actual pid is %d", syscall.Getpid())
}
