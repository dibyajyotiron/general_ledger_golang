package main

import (
	"context"
	"fmt"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"github.com/thoas/go-funk"

	"general_ledger_golang/internal/config"
	"general_ledger_golang/internal/database"
	auto "general_ledger_golang/internal/database/migrations/auto"
	"general_ledger_golang/internal/logger"
	grpcserver "general_ledger_golang/internal/transport/grpc"
	httpserver "general_ledger_golang/internal/transport/http"
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

// In case, http and grpc both are required, start this.
func main() {
	conf := config.GetConfig()
	gin.SetMode(conf.ServerSetting.RunMode)

	go grpcserver.RegisterGrpcServer(conf.ServerSetting.GrpcPort)

	router := httpserver.InitRouter()
	readTimeout := conf.ServerSetting.ReadTimeout
	writeTimeout := conf.ServerSetting.WriteTimeout
	endPoint := fmt.Sprintf("localhost:%d", conf.ServerSetting.HttpPort) // endpoint should look like 'localhost:1234'
	maxHeaderBytes := 1 << 24                                            // this is around 16 mb.

	srv := &http.Server{
		Addr:           endPoint,
		Handler:        router,
		ReadTimeout:    readTimeout,
		WriteTimeout:   writeTimeout,
		MaxHeaderBytes: maxHeaderBytes,
	}

	go func() {
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			logger.Logger.Errorf("Server Error: %v", err)
			panic(err)
		}
	}()

	logger.Logger.Infof("Http server listening on: %s", endPoint)
	_, cancel := context.WithCancel(context.Background())

	// gracefully stopping logic...
	util.GracefulShutDown(cancel, srv)
}
