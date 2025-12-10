package config

import (
	"time"
)

// App settings Section
type App struct {
	JwtSecret   string
	PrefixUrl   string
	LogSavePath string
	LogSaveName string
	LogFileExt  string
	TimeFormat  string
}

// Server settings Section
type Server struct {
	RunMode      string
	HttpPort     int
	GrpcPort     int
	ReadTimeout  time.Duration
	WriteTimeout time.Duration
	// ServiceTokenWhitelist
	//
	// Example:
	//		{"service_name":{"read":"abc","write":"cde"}}
	ServiceTokenWhitelist map[string]map[string]string
}

// Database DB settings Section
type Database struct {
	Type        string
	User        string
	Password    string
	Host        string
	Port        string
	Name        string
	TablePrefix string
	SSLMode     string
}

// Redis settings Section
type Redis struct {
	Host        string
	MaxIdle     int
	MaxActive   int
	IdleTimeout time.Duration
}

type Config struct {
	AppSetting      *App      `mapstructure:"app"`
	ServerSetting   *Server   `mapstructure:"server"`
	DatabaseSetting *Database `mapstructure:"database"`
	RedisSetting    *Redis    `mapstructure:"redis"`
}
