package config

import (
	"encoding/json"
	"fmt"
	"os"
	"regexp"
	"strings"

	"github.com/joho/godotenv"
	"github.com/spf13/viper"

	"general_ledger_golang/pkg/gotypes"
	"general_ledger_golang/pkg/logger"
)

var conf *Config

func getMapStringByValue(valueStr string) map[string]string {
	mapOfString := map[string]string{}

	confVal := ReplaceEnvInConfig([]byte(valueStr))

	e := json.Unmarshal(confVal, &mapOfString)

	if e != nil {
		fmt.Printf("\nError while unmarshaling value of %+v, err: %+v\n", valueStr, e)
		panic(e)
	}

	return mapOfString
}

// getMapBoolByValue is useful when your config is a map of bool
//
// Example:
//	 `{"BINANCE":true,"WAZIRX":true,"COINDCX":true}`
func getMapBoolByValue(valueStr string) map[string]bool {
	mapOfBool := map[string]bool{}

	confVal := ReplaceEnvInConfig([]byte(valueStr))

	e := json.Unmarshal(confVal, &mapOfBool)

	if e != nil {
		fmt.Printf("\nError while unmarshaling value of %+v, err: %+v\n", valueStr, e)
		panic(e)
	}

	return mapOfBool
}

// getMapOfMapStringByValue is useful when your config is a nested map of strings
//
// Example:
//	 `{"BROKER_UUID":{"USDT":"1"}}`
func getMapOfMapStringByValue(valueStr string) map[string]map[string]string {
	mapOfMapString := map[string]map[string]string{}
	confVal := ReplaceEnvInConfig([]byte(valueStr))

	e := json.Unmarshal(confVal, &mapOfMapString)

	if e != nil {
		fmt.Printf("\nError while unmarshaling value of %+v, err: %+v\n", valueStr, e)
		panic(e)
	}

	return mapOfMapString
}

// getMapOfMapBoolByValue is useful when your config is a nested map of booleans
//
// Example:
//	 `{"WAZIRX":{"USDT/INR":true},"BINANCE":{"BTC/USDT":true,"ADA/USDT":true}}`
func getMapOfMapBoolByValue(valueStr string) map[string]map[string]bool {
	mapOfMapBool := map[string]map[string]bool{}

	confVal := ReplaceEnvInConfig([]byte(valueStr))

	e := json.Unmarshal(confVal, &mapOfMapBool)

	if e != nil {
		fmt.Printf("\nError while unmarshaling value of %+v, err: %+v\n", valueStr, e)
		panic(e)
	}

	return mapOfMapBool
}

// ReplaceEnvInConfig will replace the ${xx} value of yaml after fetching it from .env
func ReplaceEnvInConfig(body []byte) []byte {
	search := regexp.MustCompile(`\${([^{}]+)}`)
	replacer := func(b []byte) []byte {
		group1 := search.ReplaceAllString(string(b), `$1`)
		envValue := os.Getenv(group1)
		if len(envValue) > 0 {
			return []byte(envValue)
		}
		return []byte("")
	}
	replacedBody := search.ReplaceAllFunc(body, replacer)

	return replacedBody
}

func GetProjectRoot() string {
	rootPath, _ := os.Getwd()
	return rootPath
}

// Setup is an exported method that takes the environment which starts the viper
// lib and populates the conf struct.
// configPath defaults to `./config/` if you provide "" as input
// default config file name is `config.yaml`, for us, it's not.
// APP_ENV is what determines config file name (yaml file).
func Setup(configPath string) {
	if configPath == "" {
		configPath = "./config/"
	}
	root := GetProjectRoot()
	// If providing ./xyz then no need to append root as that will break
	if !strings.Contains(configPath, "./") {
		configPath = root + configPath
	}

	var config *viper.Viper
	var err error
	config = viper.New()
	config.SetConfigType("yaml")
	config.SetConfigName(os.Getenv("APP_ENV")) // config file is named dynamically based on APP_ENV
	config.AddConfigPath(root)
	config.AddConfigPath(configPath)
	//config.AutomaticEnv()
	err = config.ReadInConfig()
	if err != nil {
		logger.Logger.Fatalf("error parsing configuration file, err: %+v", err)
	}
	for _, key := range config.AllKeys() {
		valueStr := config.GetString(key)

		// regular string config
		envOrRaw := string(ReplaceEnvInConfig([]byte(valueStr)))
		// for map of bool
		if gotypes.IsMapBool(envOrRaw) {
			mapOfBool := getMapBoolByValue(envOrRaw)
			config.Set(key, mapOfBool)
			continue
		}

		// for map of map of bool
		if gotypes.IsMapMapBool(envOrRaw) {
			mapOfMapString := getMapOfMapBoolByValue(envOrRaw)
			config.Set(key, mapOfMapString)
			continue
		}

		// for map of map of string
		if gotypes.IsMapMapString(envOrRaw) {
			mapOfMapString := getMapOfMapStringByValue(envOrRaw)
			config.Set(key, mapOfMapString)
			continue
		}

		// for map of string
		if gotypes.IsMapString(envOrRaw) {
			mapOfStr := getMapStringByValue(valueStr)
			config.Set(key, mapOfStr)
			continue
		}

		config.Set(key, envOrRaw)
	}
	if err = config.Unmarshal(&conf); err != nil {
		logger.Logger.Printf("Could not parse config, Error: %+v", err)
		panic(err)
	}
	logger.Logger.Infof("serv conf: %+v\n", *conf.ServerSetting)
}

func SetupStub(projectRootRelativePath string, ConfigDirPathFromRoot string) {
	// Load the local .env file by switching to project root
	os.Chdir(projectRootRelativePath)

	os.Setenv("APP_ENV", "local")
	godotenv.Load()
	// uncomment below logs to debug
	// myEnv, err := godotenv.Read()
	// fmt.Printf("error is: %+v", err)
	// fmt.Printf("MyENV: %+v", myEnv)
	Setup(ConfigDirPathFromRoot)
}

// GetConfig returns the config struct, where all the configs are stored as a struct.
func GetConfig() *Config {
	return conf
}
