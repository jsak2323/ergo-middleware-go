package rpc

import (
	ad "github.com/btcid/ergo-middleware-go/pkg/domain/address"
	tx "github.com/btcid/ergo-middleware-go/pkg/domain/transaction"
)

type RpcReq struct {
	RpcUser string
	Hash    string
	Arg1    string
	Arg2    string
	Arg3    string
	Nonce   string
}

type ERGORpc struct {
	addressRepo     ad.AddressRepository
	transactionRepo tx.TransactionRepository
}

func NewERGORpc(
	addressRepo ad.AddressRepository,
	transactionRepo tx.TransactionRepository,
) *ERGORpc {
	return &ERGORpc{
		addressRepo,
		transactionRepo,
	}
}
