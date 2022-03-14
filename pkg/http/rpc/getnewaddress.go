package rpc

import (
	"net/http"

	"github.com/btcid/ergo-middleware-go/pkg/lib/ergo"
	logger "github.com/btcid/ergo-middleware-go/pkg/logging"
)

type GetNewAddressRes struct {
	Content GetNewAddressResStruct
}
type GetNewAddressResStruct struct {
	Address string
	Error   string
}

func (dr *ERGORpc) GetNewAddress(req *http.Request, args *RpcReq, reply *GetNewAddressRes) error {
	defer req.Body.Close()

	newAddress, err := ergo.GetNewAddress()
	if err != nil {
		logger.ErrorLog("GetNewAddress ergo.GetNewAddress() err: " + err.Error())
		reply.Content.Error = err.Error()
		return err
	}

	logger.Log(" - New address generated: " + newAddress.Address)

	reply.Content.Address = newAddress.Address

	return nil
}
