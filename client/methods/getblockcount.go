package methods

import (
	"fmt"

	"github.com/btcid/ergo-middleware-go/cmd/config"
	"github.com/btcid/ergo-middleware-go/pkg/http/rpc"
)

func GetBlockCount() string {
	reply := rpc.GetBlockCountRes{}

	rpcReq := GenerateRpcReq("", "", "")

	xmlrpcClient := NewXmlRpcClient("localhost", config.CONF.RpcPort, "/xmlrpc")
	err := xmlrpcClient.XmlRpcCall("ERGRpc.GetBlockCount", &rpcReq, &reply)
	if err != nil {
		fmt.Println("xmlrpcClient.XmlRpcCall(\"ERGRpc.GetBlockCount\", &rpcReq, &reply) err: " + err.Error())
	}

	fmt.Println(" - Blocks: " + reply.Content.Blocks)

	if reply.Content.Error != "" {
		fmt.Println("Error: " + reply.Content.Error)
	}

	return reply.Content.Blocks
}
