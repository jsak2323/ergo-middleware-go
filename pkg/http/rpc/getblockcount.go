package rpc

import (
	"net/http"
	"strconv"

	"github.com/btcid/ergo-middleware-go/pkg/lib/ergo"
	logger "github.com/btcid/ergo-middleware-go/pkg/logging"
)

type GetBlockCountRes struct {
	Content GetBlockCountResStruct
}
type GetBlockCountResStruct struct {
	Blocks string
	Error  string
}

func (r *ERGORpc) GetBlockCount(req *http.Request, args *RpcReq, reply *GetBlockCountRes) error {
	defer req.Body.Close()

	reply.Content.Blocks = "0"

	blockCount, err := ergo.GetBlockCount()
	if err != nil {
		logger.ErrorLog("GetBlockCount ergo.GetBlockCount() err: " + err.Error())
		reply.Content.Error = err.Error()
		return err
	}

	reply.Content.Blocks = strconv.FormatInt(blockCount, 10)

	return nil
}
