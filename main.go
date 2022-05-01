package main

import (
	"fmt"
	"net/http"
	"os"
	"syscall"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"github.com/thoas/go-funk"

	"general_ledger_golang/models"
	"general_ledger_golang/pkg/config"
	"general_ledger_golang/pkg/logger"
	"general_ledger_golang/pkg/util"
	"general_ledger_golang/routers"
)

func init() {
	if funk.ContainsString([]string{"local", "localhost"}, os.Getenv("APP_ENV")) {
		err := godotenv.Load()
		logger.Logger.Info(".env Loaded")
		if err != nil {
			logger.Logger.Fatalf("Couldn't load .env for local development, error: %+v", err)
		}
	}
	config.Setup("./pkg/config/")
	models.Setup()
	logger.Setup()
	//err := gredis.Setup()
	util.Setup()

	//if err != nil {
	//	log.Fatalln(err)
	//}
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

	serverCreationDone := make(chan bool)
	go func() {
		err := srv.ListenAndServe()
		if err != nil {
			logger.Logger.Errorf("Server Error: %v", err)
			panic(err)
		}
	}()

	logger.Logger.Infof("Actual pid is %d", syscall.Getpid())
	logger.Logger.Infof("Http server listening on: %s", endPoint)

	<-serverCreationDone
}
