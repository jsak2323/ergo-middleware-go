package rpc

import (
	"net/http"
	"time"

	ad "github.com/btcid/ergo-middleware-go/pkg/domain/address"
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

func (rpc *ERGORpc) GetNewAddress(req *http.Request, args *RpcReq, reply *GetNewAddressRes) error {
	err := ergo.UnlockWallet()
	if err != nil {
		logger.ErrorLog("GetNewAddress -- ergo.UnlockWallet(), err: " + err.Error())
		return err
	}

	defer func() {
		err = ergo.LockWallet()
		if err != nil {
			logger.ErrorLog("GetNewAddress -- ergo.lockWallet(), err: " + err.Error())
		}

		err = req.Body.Close()
		if err != nil {
			logger.ErrorLog("GetNewAddress -- req.Body.Close(), err: " + err.Error())
		}
	}()

	newAddress, err := ergo.GetNewAddress()
	if err != nil {
		logger.ErrorLog("GetNewAddress -- ergo.GetNewAddress(), err: " + err.Error())
		reply.Content.Error = err.Error()
		return err
	}

	addressObj := &ad.Address{
		Created: int(time.Now().Unix()),
		Address: newAddress.Address,
	}

	err = rpc.addressRepo.Create(addressObj)
	if err != nil {
		logger.ErrorLog("GetNewAddress -- saveNewAddress rpc.addressRepo.Create(addressObj), err: " + err.Error())
		return err
	}

	logger.Log("GetNewAddress - New address generated: " + newAddress.Address)

	reply.Content.Address = newAddress.Address

	return nil
}
