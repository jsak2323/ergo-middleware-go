package rpc

import (
	"encoding/json"
	"net/http"
	"strconv"

	logger "github.com/btcid/ergo-middleware-go/pkg/logging"
)

type ListTransactionsRes struct {
	Content ListTransactionsResStruct
}
type ListTransactionsResStruct struct {
	Transactions string
	Error        string
}

type ListTransactionsTx struct {
	From          string `json:"from"`
	To            string `json:"to"`
	Hash          string `json:"hash"`
	Amount        string `json:"amount"`
	Confirmations int    `json:"confirmations"`
}

func (r *ERGORpc) ListTransactions(req *http.Request, args *RpcReq, reply *ListTransactionsRes) error {
	defer req.Body.Close()

	listTransactionTxs := []ListTransactionsTx{}
	maxConfirmations := 20

	ergoTransactions, err := r.transactionRepo.GetAll(maxConfirmations)
	if err != nil {
		logger.ErrorLog("ListTransactions ergo.ListTransactions(limit) err: " + err.Error())
		reply.Content.Error = err.Error()
		return err
	}

	for _, ergoTransaction := range ergoTransactions {

		listTransactionTx := ListTransactionsTx{
			From:          ergoTransaction.From,
			To:            ergoTransaction.To,
			Hash:          ergoTransaction.To,
			Amount:        ergoTransaction.Amount,
			Confirmations: ergoTransaction.NumConfirmation,
		}
		listTransactionTxs = append(listTransactionTxs, listTransactionTx)
	}

	listTransactionTxsJson, err := json.Marshal(listTransactionTxs)
	if err != nil {
		logger.ErrorLog("ListTransactions json.Marshal(listTransactionTxs) err: " + err.Error())
		reply.Content.Error = err.Error()
		return err
	}

	logger.InfoLog(" - ListTransactions Listed "+strconv.Itoa(len(listTransactionTxs))+" transactions.", req)

	reply.Content.Transactions = string(listTransactionTxsJson)

	return nil
}
