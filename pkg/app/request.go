package app

import (
	"github.com/astaxie/beego/validation"

	"general_ledger_golang/pkg/logger"
)

// MarkErrors logs error logs
func MarkErrors(errors []*validation.Error) {
	for _, err := range errors {
		logger.Logger.Info(err.Key, err.Message)
	}

	return
}
