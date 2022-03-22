package cron

import (
	ad "github.com/btcid/ergo-middleware-go/pkg/domain/address"
	block "github.com/btcid/ergo-middleware-go/pkg/domain/blocks"
	tx "github.com/btcid/ergo-middleware-go/pkg/domain/transaction"
	logger "github.com/btcid/ergo-middleware-go/pkg/logging"
)

type ErgoCron struct {
	addressRepo     ad.AddressRepository
	transactionRepo tx.TransactionRepository
	blockRepo       block.BlocksRepository
}

func NewErgoCron(
	addressRepo ad.AddressRepository,
	transactionRepo tx.TransactionRepository,
	blockRepo block.BlocksRepository,
) *ErgoCron {
	return &ErgoCron{
		addressRepo,
		transactionRepo,
		blockRepo,
	}
}

func (fc *ErgoCron) handleError(funcName string, err error) {
	logger.ErrorLog(funcName + " err :" + err.Error())
	panic(funcName)
}
