package cron

// import (
// 	"encoding/hex"
// 	"fmt"
// 	"net/http"

// 	"math"

// 	"github.com/btcid/ergo-middleware-go/cmd/config"
// 	"github.com/btcid/ergo-middleware-go/pkg/lib/ergo"
// 	"github.com/btcid/ergo-middleware-go/pkg/lib/util"
// 	logger "github.com/btcid/ergo-middleware-go/pkg/logging"
// )

// func (tc *ergoCron) CollectDust(w http.ResponseWriter, req *http.Request) {

// 	unspents, err := ergo.ListUnspent()
// 	if err != nil {
// 		tc.handleError("ergo.ListUnspent()", err)
// 	}

// 	if len(*(unspents)) <= 0 {
// 		logger.Log("There is no unspent transaction")
// 		fmt.Printf("\n")
// 		return
// 	}

// 	for _, transaction := range *(unspents) {
// 		fmt.Printf("%6f ||%s|\n", transaction.Amount, transaction.Address)
// 	}

// 	Inputs := []ergo.Input{}

// 	limit_tx := 500
// 	unspent_count := 0
// 	var unspents_amount float64 = 0
// 	var amount float64 = 0

// 	for _, transaction := range *(unspents) {
// 		// if transaction.Amount > 1000 { continue }
// 		if transaction.Amount < config.CONF.MinClearing {
// 			continue
// 		}
// 		if transaction.Address == config.CONF.MainAddress {
// 			continue
// 		}

// 		Input := ergo.Input{
// 			Txid: transaction.TxId,
// 			Vout: transaction.Vout,
// 		}
// 		Inputs = append(Inputs, Input)

// 		unspents_amount += transaction.Amount

// 		unspent_count += 1
// 		if unspent_count >= limit_tx {
// 			break
// 		}
// 	}

// 	if len(Inputs) <= 0 {
// 		logger.Log("There is no unspents > fee")
// 		fmt.Printf("\n")
// 		return
// 	}

// 	fmt.Printf("\nCount unspent txs    : %d\n", unspent_count)
// 	fmt.Printf("Total unspent amount : %6f\n", unspents_amount)

// 	fee := config.CONF.FeeDefault * float64(len(Inputs))
// 	feestring := fmt.Sprintf("%.8f", fee)

// 	amount = unspents_amount - fee
// 	amount_round := math.Floor(amount*1000000) / 1000000
// 	address := config.CONF.MainAddress

// 	fmt.Printf("\n=============== Data createrawtransaction ==================\n")
// 	fmt.Printf("\nfee      : %s\namount   : %8f\naddress  : %s\n\n", feestring, amount_round, address)

// 	Outputs := map[string]float64{address: amount_round}

// 	/* Create Raw Transaction */
// 	Raw, err := ergo.CreateRawTransaction(Inputs, Outputs)
// 	if err != nil {
// 		logger.ErrorLog("CreateRawTransaction ergo.CreateRawTransaction(Inputs, Outputs) err: " + err.Error())
// 		return
// 	}

// 	if Raw == "" {
// 		logger.Log("There is no raw")
// 		fmt.Printf("\n")
// 		return
// 	}

// 	logger.Log("CreateRawTransaction Raw :\n" + Raw + " ergo-middleware-go")

// 	fmt.Printf("\n")

// 	/* Sign Raw Transaction */

// 	// unlock wallet
// 	walletPassBytes, _ := hex.DecodeString(config.CONF.EncryptedPassphrase)
// 	walletPass, err := util.Decrypt(walletPassBytes, []byte(config.CONF.MailEncryptionKey))
// 	if err != nil {
// 		logger.ErrorLog("SignRawTransaction util.Decrypt(walletPassBytes, []byte(config.CONF.MailEncryptionKey)) err: " + err.Error())
// 		return
// 	}

// 	_, err = ergo.WalletPassphrase(string(walletPass), 60)
// 	if err != nil {
// 		logger.ErrorLog("SignRawTransaction ergo.WalletPassphrase(walletPass, 60) err: " + err.Error())
// 		return
// 	}

// 	signrawtransactionRes, err := ergo.SignRawTransaction(Raw)
// 	if err != nil {
// 		logger.ErrorLog("SignRawTransaction ergo.SignRawTransaction(Raw) err: " + err.Error())
// 		return
// 	}

// 	complete := signrawtransactionRes.Complete
// 	Hex := signrawtransactionRes.Hex

// 	if !complete {
// 		logger.Log("Sign raw transaction no complete")
// 		fmt.Printf("\n")
// 		return
// 	}

// 	logger.Log("SignRawTransaction Hex :\n" + Hex + " ergo-middleware-go")

// 	fmt.Printf("\n")

// 	/* Send Raw Transaction */
// 	txid, err := ergo.SendRawTransaction(Hex)
// 	if err != nil {
// 		logger.ErrorLog("SendRawTransaction ergo.SendRawTransaction(Hex) err: " + err.Error())
// 		return
// 	}

// 	if txid == "" {
// 		logger.Log("There is no txid")
// 		fmt.Printf("\n")
// 		return
// 	}

// 	amount_string := fmt.Sprintf("%.8f", amount_round)

// 	logger.Log("SendRawTransaction to address : " + address + " with amount " + amount_string + " and fee " + feestring + " returns the txid :\n" + txid)

// 	fmt.Printf("\n")

// }
