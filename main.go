package main

import (
	"fmt"
	"log"
	"syscall"

	"github.com/gin-gonic/gin"

	models "general_ledger_golang/models"
	config "general_ledger_golang/pkg/config"
	"general_ledger_golang/pkg/gredis"
	"general_ledger_golang/pkg/logger"
	"general_ledger_golang/pkg/logging"
	"general_ledger_golang/pkg/util"
	routers "general_ledger_golang/routers"

	"github.com/fvbock/endless"
)

func init() {
	config.Setup()
	models.Setup()
	logging.Setup() // should deprecate
	logger.Setup()
	err := gredis.Setup()
	util.Setup()

	if err != nil {
		log.Fatalln(err)
	}
}

// @title Golang Gin API
// @version 1.0
// @description An example of gin
func main() {
	gin.SetMode(config.ServerSetting.RunMode)

	routersInit := routers.InitRouter()
	readTimeout := config.ServerSetting.ReadTimeout
	writeTimeout := config.ServerSetting.WriteTimeout
	endPoint := fmt.Sprintf("localhost:%d", config.ServerSetting.HttpPort) // endpoint should look like 'localhost:1234'
	maxHeaderBytes := 1 << 20

	// If you want Graceful Restart, you need a Unix system and download github.com/fvbock/endless
	endless.DefaultReadTimeOut = readTimeout
	endless.DefaultWriteTimeOut = writeTimeout
	endless.DefaultMaxHeaderBytes = maxHeaderBytes
	server := endless.NewServer(endPoint, routersInit)

	server.BeforeBegin = func(add string) {
		logger.Info("Actual pid is %d", syscall.Getpid())
		logger.Info("Http server listening on: %s", endPoint)
	}

	err := server.ListenAndServe()
	if err != nil {
		log.Printf("Server err: %v", err)
	}
}
