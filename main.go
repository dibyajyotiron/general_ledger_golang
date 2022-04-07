package main

import (
	"fmt"
	"log"
	"net/http"
	"syscall"

	"github.com/gin-gonic/gin"

	"general_ledger_golang/models"
	"general_ledger_golang/pkg/config"
	"general_ledger_golang/pkg/gredis"
	"general_ledger_golang/pkg/logger"
	"general_ledger_golang/pkg/logging"
	"general_ledger_golang/pkg/util"
	"general_ledger_golang/routers"
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

func main() {
	gin.SetMode(config.ServerSetting.RunMode)

	router := routers.InitRouter()
	readTimeout := config.ServerSetting.ReadTimeout
	writeTimeout := config.ServerSetting.WriteTimeout
	endPoint := fmt.Sprintf("localhost:%d", config.ServerSetting.HttpPort) // endpoint should look like 'localhost:1234'
	maxHeaderBytes := 1 << 24                                              // this is around 16 mb.

	srv := &http.Server{
		Addr:           endPoint,
		Handler:        router,
		ReadTimeout:    readTimeout,
		WriteTimeout:   writeTimeout,
		MaxHeaderBytes: maxHeaderBytes,
	}

	done := make(chan bool)
	go func() {
		err := srv.ListenAndServe()
		if err != nil {
			logger.Logger.Errorf("Server Error: %v", err)
			panic(err)
		}
	}()

	logger.Logger.Infof("Actual pid is %d", syscall.Getpid())
	logger.Logger.Infof("Http server listening on: %s", endPoint)

	<-done
}
