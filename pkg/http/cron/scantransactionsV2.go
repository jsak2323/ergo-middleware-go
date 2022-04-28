package cron

import (
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/gorilla/mux"

	"github.com/btcid/ergo-middleware-go/cmd/config"
	mbl "github.com/btcid/ergo-middleware-go/pkg/domain/blocks"
	mtx "github.com/btcid/ergo-middleware-go/pkg/domain/transaction"
	"github.com/btcid/ergo-middleware-go/pkg/lib/ergo"
	"github.com/btcid/ergo-middleware-go/pkg/lib/util"
	logger "github.com/btcid/ergo-middleware-go/pkg/logging"
)

func (cron *ErgoCron) ScanBlockAndUpdateTransactionsV2(w http.ResponseWriter, req *http.Request) {
	vars := mux.Vars(req)

	blockNum := int64(0)
	if vars["blocknum"] != "" {
		blockNum, _ = strconv.ParseInt(vars["blocknum"], 10, 64)
	}
	cron.ScanTransactionsV2(blockNum)
}

func (cron *ErgoCron) ScanTransactionsV2(blockNum int64) (err error) {

	var (
		start         = time.Now()
		txCount       = 0
		blockDBConv   = blockNum
		maxblockCount = blockNum
	)

	if blockNum == 0 {
		// get to blocks db
		lastBlock, err := cron.blockRepo.Get()
		if err != nil {
			logger.ErrorLog("scanTransactions cron.blockRepo.Get() err: " + err.Error())
			return err
		}
		blockDBConv, _ = strconv.ParseInt(lastBlock.LastUpdatedBlockNum, 10, 64)

		maxblockCount, err = ergo.GetBlockCount()
		if err != nil {
			logger.ErrorLog("scanTransactions cron.GetBlockCount()) err: " + err.Error())
			return err
		}

	}

	// scan by wallet/transaction by minInclusionHeight by block count db - valid block && maxInclusionHeight =get block count
	err = ergo.UnlockWallet()
	if err != nil {
		logger.ErrorLog("scanTransactions unlock wallet. err: " + err.Error())
		return err
	}

	defer ergo.LockWallet()

	// insert transactions
	txCount, blockNumTx, err := cron.saveTransactionsv2(blockDBConv, maxblockCount)
	if err != nil {
		logger.ErrorLog("scanTransactions cron.saveTransactions(blockDBConv, blockCountNode) err: " + err.Error())
		return err
	}

	if blockDBConv != blockNumTx {
		LastUpdateTime := int(time.Now().Unix())

		// update last updated block num
		err = cron.blockRepo.Update(mbl.Blocks{
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
		}
	}

	elapsedMinutes := time.Since(start).Minutes()
	fmt.Println(" - Scan Block Time Elapsed: " + fmt.Sprintf("%f", elapsedMinutes) + " Minutes.")
	logger.Log(fmt.Sprintf(" - Scan transactions Finished with Total Ergo : %d new txs. ---", txCount))
	return nil
}

func (cron *ErgoCron) saveTransactionsv2(blockDBConv, blockCountNode int64) (txCount int, blockNum int64, err error) {
	if blockDBConv >= 6 {
		blockDBConv -= 6
	}

	transactions, err := ergo.ListTransactions(blockDBConv, blockCountNode)
	if err != nil {
		logger.ErrorLog("ScanTransactions lastUpdatedBlockNum convert to int64 err: " + err.Error())
		return txCount, blockNum, err
	}
	fmt.Println("transactions in node ", len(transactions.Resp), "minInclusionHeight", blockDBConv, "maxInclusionHeight ", blockCountNode)

	if len(transactions.Resp) > 0 {

		for i := len(transactions.Resp) - 1; i >= 0; i-- {
			var (
				transaction = transactions.Resp[i]
			)

			isDeposit := cron.validateTransactionsv2(transaction)

			if isDeposit {

				done, err := cron.insertTransactionsV2(transaction)
				if err != nil {
					logger.ErrorLog("ScanTransactions cron.insertTransactions(transaction), err: " + err.Error())
					return txCount, blockNum, err
				}
				if done {
					txCount++
				}

			} else {
				// withdraw
			}

			blockNum = transaction.InclusionHeight
		}
	}
	return txCount, blockNum, nil
}

func (cron *ErgoCron) insertTransactionsV2(transaction ergo.ListTransactionResp) (bool, error) {
	var (
		// from   *string
		to     string
		amount string
	)
	for _, subtx := range transaction.Outputs {

		trxDBResp, err := cron.transactionRepo.GetByHashAndAddress(transaction.ID, subtx.Address)
		if err != nil {
			logger.ErrorLog("ScanTransactions cron.transactionRepo.GetByHashAndAddress(transaction.ID, subtx.Address), err: " + err.Error())
		}

		if trxDBResp.Hash != "" {
			logger.Log(" -- Tx with hash: " + trxDBResp.Hash + ", address: " + subtx.Address + " already exists, skipping ...")
			continue
		}

		address, err := cron.addressRepo.GetByAddress(subtx.Address)
		if err != nil {
			logger.ErrorLog("ScanTransactions cron.addressRepo.GetByAddress(subtx.Address), err: " + err.Error())
			return false, err
		}

		if address.Address != "" && subtx.Address[0:2] != config.CONF.AddressFeeInit {
			to = subtx.Address
			balTemp := strconv.FormatInt(subtx.Value, 10)
			balance := util.RawToDecimal(balTemp, 9)
			amount = balance

			createTx := &mtx.Transaction{
				NumConfirmation: transaction.NumConfirmations,
				BlockNumber:     strconv.FormatInt(transaction.InclusionHeight, 10),
				To:              to,
				Amount:          amount,
				Hash:            transaction.ID,
			}

			logger.Log(" -- Inserting new ERGO Transaction (" + transaction.ID + ") ... ")
			err := cron.transactionRepo.Create(createTx)
			if err != nil {
				logger.ErrorLog("ScanTransactions cron.transactionRepo.Create(createTx), err: " + err.Error() + "txid: " + transaction.ID)
			}
			logger.Log(" -- New ERGO Tx with hash: " + transaction.ID + " and address: " + to + "inserted successfully.")
		}
	}

	return true, nil
}

func (cron *ErgoCron) validateTransactionsv2(req ergo.ListTransactionResp) (isDeposit bool) {
	// deposit
	// output harus ada address di table addresses

	// withdraw
	// mainaddress ada di output

	for _, v := range req.Outputs {
		if v.Address == config.CONF.MainAddress {

			return false
		}
	}

	return true
}
