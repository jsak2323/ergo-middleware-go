package rpc

import (
	"net/http"

	"github.com/btcid/ergo-middleware-go/pkg/lib/ergo"
)

type ValidateAddressRes struct {
	Content ValidateAddressResStruct
}
type ValidateAddressResStruct struct {
	IsValid bool
	Error   string
}

func (rpc *ERGRpc) ValidateAddress(req *http.Request, args *RpcReq, reply *ValidateAddressRes) error {
	defer req.Body.Close()

	reply.Content.IsValid = false

	address := args.Arg1

	validateAddressRes, err := ergo.ValidateAddress(address)
	if err == nil && validateAddressRes.IsValid {
		reply.Content.IsValid = true
	}

	return nil
}
