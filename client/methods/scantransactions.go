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
	err := xmlrpcClient.XmlRpcCall("ERGRpc.ScanTransactions", &rpcReq, &reply)
	if err != nil {
		fmt.Println("xmlrpcClient.XmlRpcCall(\"ERGRpc.ScanTransactions\", &rpcReq, &reply) err: " + err.Error())
	}

	if reply.Content.Error != "" {
		fmt.Println("Error: " + reply.Content.Error)
	}

	fmt.Println(" - Scan transactions Finished with Total Ergo : " + reply.Content.CountTx + " new txs. ---")

	return reply.Content.CountTx
}
