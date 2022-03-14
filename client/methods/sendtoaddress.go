package methods

// import (
// 	"fmt"

// 	"github.com/btcid/ergo-middleware-go/cmd/config"
// 	"github.com/btcid/ergo-middleware-go/pkg/http/rpc"
// )

// func SendToAddress(amount string, to string) string {
// 	reply := rpc.SendToAddressRes{}

// 	rpcReq := GenerateRpcReq(amount, to, "")

// 	xmlrpcClient := NewXmlRpcClient("localhost", config.CONF.RpcPort, "/xmlrpc")
// 	err := xmlrpcClient.XmlRpcCall("ERGORpc.SendToAddress", &rpcReq, &reply)
// 	if err != nil {
// 		fmt.Println("xmlrpcClient.XmlRpcCall(\"ERGORpc.SendToAddress\", &rpcReq, &reply) err: " + err.Error())
// 	}

// 	fmt.Println(" - TxHash: " + reply.Content.TxHash)

// 	if reply.Content.Error != "" {
// 		fmt.Println("Error: " + reply.Content.Error)
// 	}

// 	return reply.Content.TxHash
// }
