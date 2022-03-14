package methods

import (
	"fmt"

	"github.com/btcid/ergo-middleware-go/cmd/config"
	"github.com/btcid/ergo-middleware-go/pkg/http/rpc"
)

func GetBalance() string {
	reply := rpc.GetBalanceRes{}

	rpcReq := GenerateRpcReq("", "", "")

	xmlrpcClient := NewXmlRpcClient("localhost", config.CONF.RpcPort, "/xmlrpc")
	err := xmlrpcClient.XmlRpcCall("ERGORpc.GetBalance", &rpcReq, &reply)
	if err != nil {
		fmt.Println("xmlrpcClient.XmlRpcCall(\"ERGORpc.GetBalance\", &rpcReq, &reply) err: " + err.Error())
	}

	fmt.Println(" - Balance: " + reply.Content.Balance)

	if reply.Content.Error != "" {
		fmt.Println("Error: " + reply.Content.Error)
	}

	return reply.Content.Balance
}
