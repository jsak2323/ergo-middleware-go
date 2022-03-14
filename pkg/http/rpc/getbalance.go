package rpc

import (
	"fmt"
	"net/http"

	"github.com/btcid/ergo-middleware-go/pkg/lib/ergo"
	logger "github.com/btcid/ergo-middleware-go/pkg/logging"
)

type GetBalanceRes struct {
	Content GetBalanceResStruct
}
type GetBalanceResStruct struct {
	Balance string
	Error   string
}

func (dr *ERGORpc) GetBalance(req *http.Request, args *RpcReq, reply *GetBalanceRes) error {
	defer req.Body.Close()

	reply.Content.Balance = "0"

	balance, err := ergo.GetBalance()
	if err != nil {
		logger.ErrorLog("GetBalance ergo.GetBalance() err: " + err.Error())
		reply.Content.Error = err.Error()
		return err
	}

	balanceString := fmt.Sprintf("%.8f", balance)
	logger.Log(" - GetBalance Balance: " + balanceString + " ergo-middleware-go")

	reply.Content.Balance = balanceString

	return nil
}
