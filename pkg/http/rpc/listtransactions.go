package rpc

// import (
// 	"encoding/json"
// 	"fmt"
// 	"net/http"
// 	"strconv"

// 	"github.com/btcid/ergo-middleware-go/pkg/lib/ergo"
// 	logger "github.com/btcid/ergo-middleware-go/pkg/logging"
// )

// type ListTransactionsRes struct {
// 	Content ListTransactionsResStruct
// }
// type ListTransactionsResStruct struct {
// 	Transactions string
// 	Error        string
// }

// type ListTransactionsTx struct {
// 	To            string `json:"to"`
// 	Hash          string `json:"hash"`
// 	Amount        string `json:"amount"`
// 	Confirmations int    `json:"confirmations"`
// }

// func (dr *ERGORpc) ListTransactions(req *http.Request, args *RpcReq, reply *ListTransactionsRes) error {
// 	defer req.Body.Close()

// 	listTransactionTxs := []ergo.ListTransactionsResp{}
// 	maxConfirmations := 10

// 	ergoTransactions, err := ergo.ListTransactions(maxConfirmations)
// 	if err != nil {
// 		logger.ErrorLog("ListTransactions ergo.ListTransactions(limit) err: " + err.Error())
// 		reply.Content.Error = err.Error()
// 		return err
// 	}

// 	for _, ergoTransaction := range ergoTransactions.Resp {
// 		// if ergoTransaction.Category != "receive" {
// 		// 	continue
// 		// }

// 		listTransactionTx := ListTransactionsTx{
// 			To:            ergoTransaction.Address,
// 			Hash:          ergoTransaction.TxId,
// 			Amount:        fmt.Sprintf("%.8f", ergoTransaction.Amount),
// 			Confirmations: ergoTransaction.Confirmations,
// 		}
// 		listTransactionTxs = append(listTransactionTxs, listTransactionTx)
// 	}

// 	listTransactionTxsJson, err := json.Marshal(listTransactionTxs)
// 	if err != nil {
// 		logger.ErrorLog("ListTransactions json.Marshal(listTransactionTxs) err: " + err.Error())
// 		reply.Content.Error = err.Error()
// 		return err
// 	}

// 	logger.InfoLog(" - ListTransactions Listed "+strconv.Itoa(len(listTransactionTxs))+" transactions.", req)

// 	reply.Content.Transactions = string(listTransactionTxsJson)

// 	return nil
// }
