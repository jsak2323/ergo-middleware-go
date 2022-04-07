package methods

import (
	"fmt"

	"github.com/btcid/ergo-middleware-go/cmd/config"
	"github.com/btcid/ergo-middleware-go/pkg/http/rpc"
)

func GetNewAddress() string {
	reply := rpc.GetNewAddressRes{}

	rpcReq := GenerateRpcReq("", "", "")

	xmlrpcClient := NewXmlRpcClient("localhost", config.CONF.RpcPort, "/xmlrpc")
	err := xmlrpcClient.XmlRpcCall("ERGRpc.GetNewAddress", &rpcReq, &reply)
	if err != nil {
		fmt.Println("xmlrpcClient.XmlRpcCall(\"ERGRpc.GetNewAddress\", &rpcReq, &reply) err: " + err.Error())
	}

	fmt.Println(" - Address: " + reply.Content.Address)

	if reply.Content.Error != "" {
		fmt.Println("Error: " + reply.Content.Error)
	}

	return reply.Content.Address
}
