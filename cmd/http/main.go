package main

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"syscall"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"github.com/thoas/go-funk"

	"general_ledger_golang/pkg/config"
	"general_ledger_golang/pkg/database"
	"general_ledger_golang/pkg/logger"
	"general_ledger_golang/pkg/util"
	"general_ledger_golang/routers"
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
	config.Setup("./pkg/config/")
	database.Setup()
	logger.Setup()
	util.Setup()
}

func main() {
	conf := config.GetConfig()
	gin.SetMode(conf.ServerSetting.RunMode)

	router := routers.InitRouter()
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
		if err := srv.ListenAndServe(); err != http.ErrServerClosed {
			logger.Logger.Errorf("Server Error: %v", err)
			panic(err)
		}
	}()

	logger.Logger.Infof("Actual pid is %d", syscall.Getpid())
	logger.Logger.Infof("Http server listening on: %s", endPoint)
	_, cancel := context.WithCancel(context.Background())

	// gracefully stopping logic...
	util.GracefulShutDown(cancel, srv)
}
