package cron

import (
	"net/http"
	"strconv"

	"github.com/btcid/ergo-middleware-go/pkg/lib/ergo"
	logger "github.com/btcid/ergo-middleware-go/pkg/logging"
)

func (cron *ErgoCron) UpdateConfirmations(w http.ResponseWriter, req *http.Request) {
	var (
		limit int = 500
		conf  int = 15
	)

	txs, err := cron.transactionRepo.GetConfTransactions(limit, conf)
	if err != nil {
		logger.ErrorLog("UpdateConfirmations GetConfTransactions(limit,conf), err: " + err.Error())
		return
	}

	for _, tx := range txs {
		blckNum, err := strconv.ParseInt(tx.BlockNumber, 10, 64)
		if err != nil {
			logger.ErrorLog("UpdateConfirmations ParseInt(tx.BlockNumber, 10, 64), err: " + err.Error())
			continue
		}

		transactions, err := ergo.ListTransactions(blckNum, blckNum, 0)
		if err != nil {
			logger.ErrorLog("UpdateConfirmations ergo.ListTransactions(blckNum, blckNum, 0) err: " + err.Error())
			continue
		}

		err = cron.transactionRepo.UpdateNumConfirmation(tx.Id, transactions.Resp[0].NumConfirmations)
		if err != nil {
			logger.ErrorLog("UpdateConfirmations cron.transactionRepo.UpdateNumConfirmation err: " + err.Error())
			continue
		}

	}
	logger.InfoLog("UpdateConfirmations confirmations updated", req)
}
