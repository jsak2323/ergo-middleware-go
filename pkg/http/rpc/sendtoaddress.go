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

func (dr *ERGORpc) SendToAddress(req *http.Request, args *RpcReq, reply *SendToAddressRes) error {
	defer req.Body.Close()

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
