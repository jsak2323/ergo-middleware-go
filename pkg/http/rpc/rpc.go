package rpc

import (
	ad "github.com/btcid/ergo-middleware-go/pkg/domain/address"
	bl "github.com/btcid/ergo-middleware-go/pkg/domain/blocks"
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
	blockRepo       bl.BlocksRepository
}

func NewERGORpc(
	addressRepo ad.AddressRepository,
	transactionRepo tx.TransactionRepository,
	blockRepo bl.BlocksRepository,
) *ERGORpc {
	return &ERGORpc{
		addressRepo,
		transactionRepo,
		blockRepo,
	}
}
