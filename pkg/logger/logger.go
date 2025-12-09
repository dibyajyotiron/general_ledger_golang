package logger

import (
	"io"
	"os"
	"path"
	"runtime"
	"strconv"

	"github.com/sirupsen/logrus"
	"github.com/thoas/go-funk"
)

var Logger = logrus.New()

type Fields map[string]interface{}

func Setup() {
	Logger.SetOutput(getWriter())
	layout := "2006-Jan-02T15:04:05.000Z"

	if funk.ContainsString([]string{"prod", "production", "release"}, os.Getenv("APP_ENV")) {
		Logger.Formatter = &logrus.JSONFormatter{
			TimestampFormat: layout,
			PrettyPrint:     true,
		}
	} else {
		Logger.Formatter = &logrus.TextFormatter{
			FullTimestamp:             true,
			TimestampFormat:           layout,
			ForceColors:               true,
			EnvironmentOverrideColors: true,
			CallerPrettyfier:          CallerPrettyfier,
		}
	}
}

func getWriter() io.Writer {
	writeToFile, notSet := os.LookupEnv("WRITE_LOG_TO_FILE")
	writeToFileFlag, _ := strconv.ParseBool(writeToFile)

	if notSet || writeToFileFlag == false {
		return os.Stdout
	}

	appEnv := os.Getenv("APP_ENV")
	file, err := os.OpenFile("application.log", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		Logger.Errorf("Failed to open log file: %v", err)
		return os.Stdout
	}
	if appEnv != "local" {
		return file
	}
	return os.Stdout
}

func CallerPrettyfier(frame *runtime.Frame) (function string, file string) {
	fileName := " " + path.Base(frame.File) + ":" + strconv.Itoa(frame.Line)
	//return frame.Function, fileName
	return "", fileName
}
