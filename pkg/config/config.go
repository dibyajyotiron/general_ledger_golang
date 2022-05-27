package config

import (
	"log"
	"os"
	"regexp"

	"github.com/spf13/viper"

	"general_ledger_golang/pkg/logger"
)

var conf *Config

// ReplaceEnvInConfig will replace the ${xx} value of yaml after fetching it from .env
func ReplaceEnvInConfig(body []byte) []byte {
	search := regexp.MustCompile(`\${([^{}]+)}`)
	replacedBody := search.ReplaceAllFunc(body, func(b []byte) []byte {
		group1 := search.ReplaceAllString(string(b), `$1`)
		envValue := os.Getenv(group1)
		if len(envValue) > 0 {
			return []byte(envValue)
		}
		return []byte("")
	})

	return replacedBody
}

func GetProjectRoot() string {
	rootPath, _ := os.Getwd()
	return rootPath
}

// Setup is an exported method that takes the environment which starts the viper
// lib and populates the conf struct.
// configPath defaults to ./config/ if you provide "" as input
func Setup(configPath string) {
	if configPath == "" {
		configPath = "./config/"
	}
	root := GetProjectRoot()
	var config *viper.Viper
	var err error
	config = viper.New()
	config.SetConfigType("yaml")
	config.SetConfigName(os.Getenv("APP_ENV")) // config file is named dynamically based on APP_ENV
	config.AddConfigPath(root)
	config.AddConfigPath(configPath)
	err = config.ReadInConfig()
	if err != nil {
		log.Fatalf("error parsing configuration file, err: %+v", err)
	}
	for _, key := range config.AllKeys() {
		value := config.GetString(key)
		envOrRaw := ReplaceEnvInConfig([]byte(value))
		config.Set(key, string(envOrRaw))
	}
	if err = config.Unmarshal(&conf); err != nil {
		logger.Logger.Printf("Could not parse config, Error: %+v", err)
	}
}

// GetConfig returns the config struct, where all the configs are stored as a struct.
func GetConfig() *Config {
	return conf
}
