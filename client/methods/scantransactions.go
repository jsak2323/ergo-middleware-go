package methods

import (
	"fmt"

	"github.com/btcid/ergo-middleware-go/cmd/config"
	"github.com/btcid/ergo-middleware-go/pkg/http/rpc"
)

func ScanTransactions() string {
	reply := rpc.ScanTransactionsRes{}

	rpcReq := GenerateRpcReq("", "", "")

	xmlrpcClient := NewXmlRpcClient("localhost", config.CONF.RpcPort, "/xmlrpc")
	err := xmlrpcClient.XmlRpcCall("ERGORpc.ScanTransactions", &rpcReq, &reply)
	if err != nil {
		fmt.Println("xmlrpcClient.XmlRpcCall(\"ERGORpc.ScanTransactions\", &rpcReq, &reply) err: " + err.Error())
	}

	if reply.Content.Error != "" {
		fmt.Println("Error: " + reply.Content.Error)
	}

	return reply.Content.CountTx
}
