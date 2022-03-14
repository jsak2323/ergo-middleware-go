package rpc

type RpcReq struct {
	RpcUser string
	Hash    string
	Arg1    string
	Arg2    string
	Arg3    string
	Nonce   string
}

type ERGORpc struct{}

func NewERGORpc() *ERGORpc {
	return &ERGORpc{}
}
