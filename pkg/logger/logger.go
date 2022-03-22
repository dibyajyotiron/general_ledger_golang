package logger

import (
	"general_ledger_golang/pkg/util"
	"io"
	"os"
	"path"
	"runtime"
	"strconv"

	"github.com/sirupsen/logrus"
)

var logger = logrus.New()

func Setup() {
	logger.SetOutput(getWriter())
	layout := "2006-Jan-02T15:04:05.000Z"

	if util.Contains(os.Getenv("APP_ENV"), []string{"prod", "production", "release"}) {
		logger.Formatter = &logrus.JSONFormatter{
			TimestampFormat:  layout,
			PrettyPrint:      true,
			CallerPrettyfier: CallerPrettyfier,
		}
	} else {
		logger.Formatter = &logrus.TextFormatter{
			FullTimestamp:             true,
			TimestampFormat:           layout,
			ForceColors:               true,
			EnvironmentOverrideColors: true,
			CallerPrettyfier:          CallerPrettyfier,
		}
	}
}

func getWriter() io.Writer {
	APP_ENV := os.Getenv("APP_ENV")
	file, err := os.OpenFile("application.log", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		logger.Errorf("Failed to open log file: %v", err)
		return os.Stdout
	}
	if APP_ENV != "local" {
		return file
	}
	return os.Stdout
}

func CallerPrettyfier(frame *runtime.Frame) (function string, file string) {
	fileName := " " + path.Base(frame.File) + ":" + strconv.Itoa(frame.Line)
	//return frame.Function, fileName
	return "", fileName
}
