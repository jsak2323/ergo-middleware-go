package rpc

import (
	"errors"
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/btcid/ergo-middleware-go/cmd/config"
	mbl "github.com/btcid/ergo-middleware-go/pkg/domain/blocks"
	mtx "github.com/btcid/ergo-middleware-go/pkg/domain/transaction"
	"github.com/btcid/ergo-middleware-go/pkg/lib/ergo"
	"github.com/btcid/ergo-middleware-go/pkg/lib/util"
	logger "github.com/btcid/ergo-middleware-go/pkg/logging"
)

// get block count
// scan by wallet/transaction by minInclusionHeight by block count db - valid block && maxInclusionHeight =get block count
// insert to transactions
// update last block nums with get block count)

// 2 data: luar & indodax = depo
// —
// 2 data: indodax & indodax = skip
// —
// 3 data: luar & ke 2 akun indodax = depo
// —
// 2 data: luar & mainaddress = skip

type ScanTransactionsRes struct {
	Content ScanTransactionsResStruct
}

type ScanTransactionsResStruct struct {
	CountTx string
	Error   string
}

func (rpc *ERGORpc) ScanTransactions(req *http.Request, args *RpcReq, reply *ScanTransactionsRes) (err error) {
	defer req.Body.Close()
	var (
		start   = time.Now()
		txCount = 0
	)

	// get to blocks db
	lastBlock, err := rpc.blockRepo.Get()
	if err != nil {
		logger.ErrorLog("ScanTransactions rpc.blockRepo.Get() err: " + err.Error())
		reply.Content.Error = err.Error()
		return err
	}
	blockDBConv, _ := strconv.ParseInt(lastBlock.LastUpdatedBlockNum, 10, 64)

	// scan by wallet/transaction by minInclusionHeight by block count db - valid block && maxInclusionHeight =get block count
	err = ergo.UnlockWallet()
	if err != nil {
		logger.ErrorLog("ScanTransactions unlock wallet. err: " + err.Error())
		reply.Content.Error = err.Error()
		return err
	}

	defer ergo.LockWallet()

	// insert transactions
	txCount, blockNumTx, err := rpc.saveTransactions(blockDBConv, 0)
	if err != nil {
		logger.ErrorLog("ScanTransactions rpc.saveTransactions(blockDBConv, blockCountNode) err: " + err.Error())
		reply.Content.Error = err.Error()
		return err
	}

	if blockDBConv != blockNumTx {
		LastUpdateTime := int(time.Now().Unix())

		// update last updated block num
		err = rpc.blockRepo.Update(mbl.Blocks{
			LastUpdateTime:      LastUpdateTime,
			LastUpdatedBlockNum: fmt.Sprintf("%d", blockNumTx),
		})
		fmt.Printf(
			"\n -- Last update time: %d, Last updated block num: %d\n",
			LastUpdateTime,
			blockNumTx,
		)
		if err != nil {
			logger.ErrorLog(err.Error())
			reply.Content.Error = err.Error()
		}
	}

	elapsedMinutes := time.Since(start).Minutes()
	fmt.Println(" - Scan Block Time Elapsed: " + fmt.Sprintf("%f", elapsedMinutes) + " Minutes.")
	logger.Log(fmt.Sprintf(" - Scan transactions Finished with Total Ergo : %d new txs. ---", txCount))
	reply.Content.CountTx = strconv.Itoa(txCount)
	return nil
}

func (rpc *ERGORpc) saveTransactions(blockDBConv, blockCountNode int64) (txCount int, blockNum int64, err error) {
	if blockDBConv >= 15 {
		blockDBConv -= 15
	}

	transactions, err := ergo.ListTransactions(blockDBConv, blockCountNode)
	if err != nil {
		logger.ErrorLog("ScanTransactions lastUpdatedBlockNum convert to int64 err: " + err.Error())
		return txCount, blockNum, err
	}

	if len(transactions.Resp) > 1 {

		for i := len(transactions.Resp) - 1; i >= 0; i-- {
			var (
				transaction = transactions.Resp[i]
			)

			// 2 data: luar & indodax = depo
			// —
			// 2 data: indodax & indodax = skip
			// —
			// 3 data: luar & ke 2 akun indodax = depo
			// —
			// 2 data: luar & mainaddress = skip

			valid, err := rpc.validateTransactions(transaction)
			if err != nil {
				logger.ErrorLog("ScanTransactions rpc.validateTransactions(transaction), err: " + err.Error())
				// return txCount, blockNum, err
			}

			if valid {

				done, err := rpc.insertTransactions(transaction)
				if err != nil {
					logger.ErrorLog("ScanTransactions rpc.insertTransactions(transaction), err: " + err.Error())
					return txCount, blockNum, err
				}
				if done {
					txCount++
				}
			}

			blockNum = transaction.InclusionHeight
		}
	}
	return txCount, blockNum, nil
}

func (rpc *ERGORpc) insertTransactions(transaction ergo.ListTransactionResp) (bool, error) {
	var (
		from   string
		to     string
		amount string
	)
	for _, subtx := range transaction.Outputs {

		trxDBResp, err := rpc.transactionRepo.GetByHashAndAddress(transaction.ID, subtx.Address)
		if err != nil {
			logger.ErrorLog("ScanTransactions rpc.transactionRepo.GetByHashAndAddress(transaction.ID, subtx.Address), err: " + err.Error())
			return false, err
		}

		if trxDBResp.Hash != "" {
			logger.Log(" -- Tx with hash: " + trxDBResp.Hash + ", address: " + subtx.Address + " already exists, skipping ...")
			break
		}

		address, err := rpc.addressRepo.GetByAddress(subtx.Address)
		if err != nil {
			logger.ErrorLog("ScanTransactions rpc.addressRepo.GetByAddress(subtx.Address), err: " + err.Error())
			return false, err
		}

		if address.Address == "" {
			from = subtx.Address
		} else if address.Address != "" && address.Address[0:3] != config.CONF.AddressFeeInit {
			to = subtx.Address
			balTemp := strconv.FormatInt(subtx.Value, 10)
			balance := util.RawToDecimal(balTemp, 9)
			amount = balance
		}
	}
	if from != "" && to != "" && amount != "" {

		createTx := &mtx.Transaction{
			NumConfirmation: transaction.NumConfirmations,
			BlockNumber:     strconv.FormatInt(transaction.InclusionHeight, 10),
			From:            from,
			To:              to,
			Amount:          amount,
			Hash:            transaction.ID,
		}

		logger.Log(" -- Inserting new ERGO Transaction (" + transaction.ID + ") ... ")
		err := rpc.transactionRepo.Create(createTx)
		if err != nil {
			logger.ErrorLog("ScanTransactions ron.transactionRepo.Create(createTx), err: " + err.Error())
			return false, err
		}
		logger.Log(" -- New ERGO Tx with hash: " + transaction.ID + " inserted successfully.")
	}

	return false, nil
}

func (rpc *ERGORpc) validateTransactions(req ergo.ListTransactionResp) (result bool, err error) {

	var (
		externalAddress int
		internalAddress int
	)
	if len(req.Outputs) > 4 {
		err = errors.New("scan transactions, multi payment cant included to transactions")
		return false, err
	}

	for _, subtx := range req.Outputs {

		if subtx.Address[0:3] == config.CONF.AddressFeeInit {
			continue
		}
		if subtx.Address == config.CONF.MainAddress {
			err = errors.New("scan transactions, main address cant be inserted to transactions")
			return false, err
		}

		address, err := rpc.addressRepo.GetByAddress(subtx.Address)
		if err != nil {
			logger.ErrorLog("ScanTransactions rpc.addressRepo.GetByAddress(subtx.Address), err: " + err.Error())
			return false, err
		}

		if address.Address != "" {
			internalAddress++
		} else {
			externalAddress++
		}

	}

	// 2 data: luar & mainaddress = skip 3
	// -
	// 2 data: indodax & indodax = skip 3
	// —
	// 2 data: luar & indodax = depo 3 -> luar from, indodax to
	// —
	// 3 data: luar & ke 2 akun indodax = depo 4 -> luar from , indodax to 1 aja
	// —

	if externalAddress == 1 && internalAddress >= 1 {
		return true, nil
	}

	return false, nil
}
