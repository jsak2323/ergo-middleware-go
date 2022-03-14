package cron

import (
	logger "github.com/btcid/ergo-middleware-go/pkg/logging"
)

type ergoCron struct{}

func NewergoCron() *ergoCron {
	return &ergoCron{}
}

func (fc *ergoCron) handleError(funcName string, err error) {
	logger.ErrorLog(funcName + " err :" + err.Error())
	panic(funcName)
}
