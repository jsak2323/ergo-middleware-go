package methods

import (
	"bytes"
	"crypto/md5"
	"crypto/sha256"
	"encoding/hex"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/divan/gorilla-xmlrpc/xml"

	"github.com/btcid/ergo-middleware-go/cmd/config"
	"github.com/btcid/ergo-middleware-go/pkg/lib/util"
	logger "github.com/btcid/ergo-middleware-go/pkg/logging"
)

type RpcConfig struct {
	User     string
	Password string
	Hashkey  string
}

type RpcReq struct {
	RpcUser string
	Hash    string
	Arg1    string
	Arg2    string
	Arg3    string
	Nonce   string
}

type XmlRpc struct {
	Host string
	Port string
	Path string
}

func NewXmlRpcClient(host string, port string, path string) *XmlRpc {
	return &XmlRpc{
		host, port, path,
	}
}
func (xr *XmlRpc) XmlRpcCall(method string, args *RpcReq, reply interface{}) error {
	buf, err := xml.EncodeClientRequest(method, args)
	if err != nil {
		logger.ErrorLog("xml.EncodeClientRequest(method, args) err: " + err.Error())
		return err
	}

	url := "http://" + xr.Host + ":" + xr.Port + xr.Path
	httpClient := &http.Client{
		Timeout: 1000 * time.Second,
	}
	res, err := httpClient.Post(url, "text/xml", bytes.NewBuffer(buf))
	if err != nil {
		logger.ErrorLog("httpClient.Post(url, \"text/xml\", bytes.NewBuffer(buf)) err: " + err.Error())
		return err
	}
	defer res.Body.Close()

	err = xml.DecodeClientResponse(res.Body, reply)
	if err != nil {
		logger.ErrorLog("xml.DecodeClientResponse(res.Body, reply) err: " + err.Error())
		return err
	}

	return nil
}

func GenerateRpcReq(arg1 string, arg2 string, arg3 string) RpcReq {
	rpcConfig := RpcConfig{
		User:     config.CONF.RpcUser,
		Password: config.CONF.RpcPass,
		Hashkey:  config.CONF.RpcHashkey,
	}

	hashkey, nonce := generateHashkey(rpcConfig)

	return RpcReq{
		RpcUser: rpcConfig.User,
		Hash:    hashkey,
		Arg1:    arg1,
		Arg2:    arg2,
		Arg3:    arg3,
		Nonce:   nonce,
	}
}

func generateHashkey(rpcConfig RpcConfig) (digestSha256String string, nonce string) {
	mt := util.Microtime()
	nonce = strings.ReplaceAll(strconv.FormatFloat(mt, 'f', 9, 64), ".", "")

	unixTime := time.Now().Unix()
	this15m := unixTime / 60

	pass := rpcConfig.Password
	hashkey := rpcConfig.Hashkey

	digest := pass + strconv.FormatInt(this15m, 10) + hashkey + nonce
	digestMd5 := md5.Sum([]byte(digest))
	digestMd5String := hex.EncodeToString(digestMd5[:])
	digestSha256 := sha256.Sum256([]byte(digestMd5String))
	digestSha256String = hex.EncodeToString(digestSha256[:])

	return
}
