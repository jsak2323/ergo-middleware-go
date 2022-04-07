package methods

import (
	"fmt"
	"strconv"

	"github.com/btcid/ergo-middleware-go/cmd/config"
	"github.com/btcid/ergo-middleware-go/pkg/http/rpc"
)

func ValidateAddress(address string) bool {
	reply := rpc.ValidateAddressRes{}

	rpcReq := GenerateRpcReq(address, "", "")

	xmlrpcClient := NewXmlRpcClient("localhost", config.CONF.RpcPort, "/xmlrpc")
	err := xmlrpcClient.XmlRpcCall("ERGRpc.ValidateAddress", &rpcReq, &reply)
	if err != nil {
		fmt.Println("xmlrpcClient.XmlRpcCall(\"ERGRpc.ValidateAddress\", &rpcReq, &reply) err: " + err.Error())
	}

	fmt.Println(" - IsValid: " + strconv.FormatBool(reply.Content.IsValid))

	if reply.Content.Error != "" {
		fmt.Println("Error: " + reply.Content.Error)
	}

	return reply.Content.IsValid
}
