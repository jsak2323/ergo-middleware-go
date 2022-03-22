package rpc

import (
	"net/http"

	"github.com/btcid/ergo-middleware-go/pkg/lib/ergo"
	logger "github.com/btcid/ergo-middleware-go/pkg/logging"
)

type SendToAddressRes struct {
	Content SendToAddressResStruct
}
type SendToAddressResStruct struct {
	TxHash string
	Error  string
}

func (r *ERGORpc) SendToAddress(req *http.Request, args *RpcReq, reply *SendToAddressRes) error {
	err := ergo.UnlockWallet()
	if err != nil {
		logger.ErrorLog("SendToAddress -- ergo.UnlockWallet(), err: " + err.Error())
		return err
	}

	defer func() {
		err = ergo.LockWallet()
		if err != nil {
			logger.ErrorLog("SendToAddress -- ergo.lockWallet(), err: " + err.Error())
		}

		err = req.Body.Close()
		if err != nil {
			logger.ErrorLog("SendToAddress -- req.Body.Close(), err: " + err.Error())
		}
	}()

	amountInDecimal := args.Arg1
	toAddr := args.Arg2

	sendMsg := " - Transferring " + amountInDecimal + " ergo to " + toAddr
	logger.Log(sendMsg + " ...")

	txHash, err := ergo.SendToAddress(toAddr, amountInDecimal)
	if err != nil {
		logger.ErrorLog("SendToAddress ergo.SendToAddress(toAddr, amountInDecimal) err: " + err.Error())
		reply.Content.Error = err.Error()
		return err
	}
	logger.Log(sendMsg + " Successful. TxHash: " + txHash)

	reply.Content.TxHash = txHash

	return nil
}
