package config

import (
	"bytes"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/joho/godotenv"
	"github.com/spf13/viper"

	"general_ledger_golang/internal/logger"
)

var conf *Config

// readConfigFile expands environment variables and feeds the YAML into viper.
func readConfigFile(v *viper.Viper, absoluteConfigDir string, env string) error {
	configFile := filepath.Join(absoluteConfigDir, fmt.Sprintf("%s.yaml", env))
	fileBytes, err := os.ReadFile(configFile)
	if err != nil {
		return fmt.Errorf("read config file: %w", err)
	}

	expanded := os.ExpandEnv(string(fileBytes))
	if err := v.ReadConfig(bytes.NewBufferString(expanded)); err != nil {
		return fmt.Errorf("parse config file: %w", err)
	}

	return nil
}

// Setup wires viper, honors env overrides, and unmarshals into Config in a step-by-step, readable manner.
func Setup(configPath string) {
	if configPath == "" {
		configPath = "./internal/config/"
	}

	projectRoot, _ := os.Getwd()
	absoluteConfigDir := configPath
	if !filepath.IsAbs(configPath) {
		absoluteConfigDir = filepath.Clean(filepath.Join(projectRoot, configPath))
	}

	activeEnv := os.Getenv("APP_ENV")
	if activeEnv == "" {
		activeEnv = "local"
	}

	v := viper.New()
	v.SetConfigType("yaml")
	v.SetConfigName(activeEnv)
	v.AddConfigPath(absoluteConfigDir)
	v.AddConfigPath(projectRoot)

	// Allow environment variables to override config keys. Dots become underscores (server.httpport => APP_SERVER_HTTPPORT).
	v.SetEnvPrefix("APP")
	v.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))
	v.AutomaticEnv()

	if err := readConfigFile(v, absoluteConfigDir, activeEnv); err != nil {
		logger.Logger.Fatalf("Config load failed: %+v", err)
	}

	if err := v.Unmarshal(&conf); err != nil {
		logger.Logger.Fatalf("Config decode failed: %+v", err)
	}

	logger.Logger.Infof("Server config loaded: %+v", conf.ServerSetting)
}

func SetupStub(projectRootRelativePath string, configDirPathFromRoot string) {
	// Load the local .env file by switching to project root
	_ = os.Chdir(projectRootRelativePath)

	os.Setenv("APP_ENV", "local")
	_ = godotenv.Load()
	Setup(configDirPathFromRoot)
}

// GetConfig returns the populated Config struct.
func GetConfig() *Config {
	return conf
}
