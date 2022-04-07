package methods

import (
	"fmt"

	"github.com/btcid/ergo-middleware-go/cmd/config"
	"github.com/btcid/ergo-middleware-go/pkg/http/rpc"
)

func ListTransactions(limit string) string {
	reply := rpc.ListTransactionsRes{}

	rpcReq := GenerateRpcReq(limit, "", "")

	xmlrpcClient := NewXmlRpcClient("localhost", config.CONF.RpcPort, "/xmlrpc")
	err := xmlrpcClient.XmlRpcCall("ERGRpc.ListTransactions", &rpcReq, &reply)
	if err != nil {
		fmt.Println("xmlrpcClient.XmlRpcCall(\"ERGRpc.ListTransactions\", &rpcReq, &reply) err: " + err.Error())
	}

	fmt.Println(" - Transactions: " + reply.Content.Transactions)

	if reply.Content.Error != "" {
		fmt.Println("Error: " + reply.Content.Error)
	}

	return reply.Content.Transactions
}
