package cron

import (
	"errors"
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
			)

			// 2 data: luar & indodax = depo
			// —
			// 2 data: indodax & indodax = skip
			// —
			// 3 data: luar & ke 2 akun indodax = depo
			// —
			// 2 data: luar & mainaddress = skip

			valid, err := cron.validateTransactions(transaction)
			if err != nil {
				logger.ErrorLog("scanTransactions cron.validateTransactions(transaction), err: " + err.Error())
				// return txCount, blockNum, err
			}

			if valid {

				done, err := cron.insertTransactions(transaction)
				if err != nil {
					logger.ErrorLog("scanTransactions cron.insertTransactions(transaction), err: " + err.Error())
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

func (cron *ErgoCron) insertTransactions(transaction ergo.ListTransactionResp) (bool, error) {
	var (
		from   string
		to     string
		amount string
	)
	for _, subtx := range transaction.Outputs {

		trxDBResp, err := cron.transactionRepo.GetByHashAndAddress(transaction.ID, subtx.Address)
		if err != nil {
			logger.ErrorLog("scanTransactions cron.transactionRepo.GetByHashAndAddress(transaction.ID, subtx.Address), err: " + err.Error())
			return false, err
		}

		if trxDBResp.Hash != "" {
			logger.Log(" -- Tx with hash: " + trxDBResp.Hash + ", address: " + subtx.Address + " already exists, skipping ...")
			break
		}

		address, err := cron.addressRepo.GetByAddress(subtx.Address)
		if err != nil {
			logger.ErrorLog("scanTransactions cron.addressRepo.GetByAddress(subtx.Address), err: " + err.Error())
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
		err := cron.transactionRepo.Create(createTx)
		if err != nil {
			logger.ErrorLog("scanTransactions ron.transactionRepo.Create(createTx), err: " + err.Error())
			return false, err
		}
		logger.Log(" -- New ERGO Tx with hash: " + transaction.ID + " inserted successfully.")
	}

	return false, nil
}

func (cron *ErgoCron) validateTransactions(req ergo.ListTransactionResp) (result bool, err error) {

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

		address, err := cron.addressRepo.GetByAddress(subtx.Address)
		if err != nil {
			logger.ErrorLog("scanTransactions cron.addressRepo.GetByAddress(subtx.Address), err: " + err.Error())
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
