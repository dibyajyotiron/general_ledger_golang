package config

import (
	"log"
	"time"

	"github.com/gookit/ini/v2"
	"github.com/gookit/ini/v2/dotnev"
)

// mapTo map section
// This will map the section part of .ini with the interfaces like (RedisSetting, DatabaseSetting etc..).
func mapTo(section string, conf interface{}) {
	err := ini.MapStruct(section, conf)

	if err != nil {
		log.Fatalf("Cfg.MapTo %s err: %v", section, err)
	}
}

// Setup initialize the configuration instance
func Setup() {
	err := dotnev.Load("./", ".env")
	err2 := ini.LoadFiles("conf/app.ini")

	if err != nil || err2 != nil {
		log.Fatalf("config.Setup, fail to parse 'conf/app.ini': %v", err)
	}

	mapTo("app", AppSetting)
	mapTo("server", ServerSetting)
	mapTo("database", DatabaseSetting)
	mapTo("redis", RedisSetting)

	AppSetting.ImageMaxSize = AppSetting.ImageMaxSize * 1024 * 1024
	ServerSetting.ReadTimeout = ServerSetting.ReadTimeout * time.Second
	ServerSetting.WriteTimeout = ServerSetting.WriteTimeout * time.Second
	RedisSetting.IdleTimeout = RedisSetting.IdleTimeout * time.Second
}
