package cron

import (
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/gorilla/mux"

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

func (cron *ErgoCron) ScanBlockAndUpdateTransactions(w http.ResponseWriter, req *http.Request) {
	vars := mux.Vars(req)

	blockNum := int64(0)
	if vars["blocknum"] != "" {
		blockNum, _ = strconv.ParseInt(vars["blocknum"], 10, 64)
	}
	cron.ScanTransactions(blockNum)
}

func (cron *ErgoCron) ScanTransactions(blockNum int64) (err error) {

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

	}

	// scan by wallet/transaction by minInclusionHeight by block count db - valid block && maxInclusionHeight =get block count
	err = ergo.UnlockWallet()
	if err != nil {
		logger.ErrorLog("scanTransactions unlock wallet. err: " + err.Error())
		return err
	}

	defer ergo.LockWallet()

	// insert transactions
	txCount, blockNumTx, err := cron.saveTransactions(blockDBConv, maxblockCount)
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

func (cron *ErgoCron) saveTransactions(blockDBConv, blockCountNode int64) (txCount int, blockNum int64, err error) {
	if blockDBConv >= 15 {
		blockDBConv -= 15
	}

	transactions, err := ergo.ListTransactions(blockDBConv, blockCountNode)
	if err != nil {
		logger.ErrorLog("scanTransactions lastUpdatedBlockNum convert to int64 err: " + err.Error())
		return txCount, blockNum, err
	}

	if len(transactions.Resp) > 1 {

		for i := len(transactions.Resp) - 1; i >= 0; i-- {
			var (
				transaction = transactions.Resp[i]
				txFound     = false
			)

			if len(transaction.Outputs) > 1 {
				for _, subtx := range transaction.Outputs {

					if !txFound {
						trxDBResp, err := cron.transactionRepo.GetByHashAndAddress(transaction.ID, subtx.Address)
						if err != nil {
							logger.ErrorLog("scanTransactions cron.transactionRepo.GetByHashAndAddress(transaction.ID, subtx.Address), err: " + err.Error())
							return txCount, blockNum, err
						}

						if trxDBResp.Hash != "" {
							logger.Log(" -- Tx with hash: " + trxDBResp.Hash + ", address: " + subtx.Address + " already exists, skipping ...")
							continue
						}

						address, err := cron.addressRepo.GetByAddress(subtx.Address)
						if err != nil {
							logger.ErrorLog("scanTransactions cron.addressRepo.GetByAddress(subtx.Address), err: " + err.Error())
							return txCount, blockNum, err
						}

						if address.Address != "" {

							balTemp := strconv.FormatInt(subtx.Value, 10)
							balance := util.RawToDecimal(balTemp, 9)
							createTx := &mtx.Transaction{
								NumConfirmation: transaction.NumConfirmations,
								BlockNumber:     strconv.FormatInt(transaction.InclusionHeight, 10),
								To:              subtx.Address,
								Amount:          balance,
								Hash:            transaction.ID,
							}

							logger.Log(" -- Inserting new ERGO Transaction (" + transaction.ID + ") ... ")
							err := cron.transactionRepo.Create(createTx)
							if err != nil {
								logger.ErrorLog("scanTransactions ron.transactionRepo.Create(createTx), err: " + err.Error())
								return txCount, blockNum, err
							}
							logger.Log(" -- New ERGO Tx with hash: " + transaction.ID + " inserted successfully.")
							txFound = true
							txCount++
						}

					}

				}
			}
			blockNum = transaction.InclusionHeight
		}
	}
	return txCount, blockNum, nil
}
