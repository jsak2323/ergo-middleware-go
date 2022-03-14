package auth

import (
	"bytes"
	"errors"
	"io/ioutil"
	"net/http"
	"strconv"
	"strings"
	"time"

	"crypto/md5"
	"crypto/sha256"
	"encoding/hex"

	"database/sql"

	gorillaxml "github.com/divan/gorilla-xmlrpc/xml"
	_ "github.com/mattn/go-sqlite3"

	config "github.com/btcid/ergo-middleware-go/cmd/config"
	"github.com/btcid/ergo-middleware-go/pkg/http/rpc"
	logger "github.com/btcid/ergo-middleware-go/pkg/logging"
)

const nonceDbPath = "/var/db/nonce.sqlite"

func AuthorizeXmlRequest(req *http.Request) error {
	args := rpc.RpcReq{}
	rawXml, _ := ioutil.ReadAll(req.Body)

	req.Body = ioutil.NopCloser(bytes.NewBuffer(rawXml))
	xmlCodec := gorillaxml.Codec{}
	codecRequest := xmlCodec.NewRequest(req)

	method, _ := codecRequest.Method()
	rpcMethod := strings.SplitAfter(method, ".")[1]

	err := codecRequest.ReadRequest(&args)
	if err != nil {
		logger.ErrorLog("codecRequest.ReadRequest(&args) err : " + err.Error())
	}
	logger.InfoLog(" - AUTH -- XML Request Params => Arg1: "+args.Arg1+", Arg2: "+args.Arg2+", Arg3: "+args.Arg3, req)

	if err := AuthorizeRequest(&args, rpcMethod); err != nil {
		return err
	}

	req.Body = ioutil.NopCloser(bytes.NewBuffer(rawXml))

	return nil
}

func AuthorizeRequest(args *rpc.RpcReq, rpcMethod string) error {

	if args.RpcUser != config.CONF.RpcUser {
		return errors.New("Invalid User.")
	}

	if rpcMethod != "GetBlockCount" {
		if err := validateAndUpdateNonce(args.Nonce); err != nil {
			return err
		}
	}

	if err := validateHashKey(args.Hash, args.Nonce); err != nil {
		return err
	}

	return nil
}

func validateAndUpdateNonce(reqNonce string) error {
	var nonce string

	db, err := sql.Open("sqlite3", nonceDbPath)
	if err != nil {
		logger.ErrorLog(err.Error())
		return err
	}
	defer db.Close()

	rows, err := db.Query("SELECT nonce FROM nonce")
	if err != nil {
		logger.ErrorLog(err.Error())
		return err
	}
	defer rows.Close()

	for rows.Next() {
		err = rows.Scan(&nonce)
		if err != nil {
			logger.ErrorLog(err.Error())
			return err
		}
	}

	reqNonceNum, err := strconv.Atoi(reqNonce)
	nonceNum, err := strconv.Atoi(nonce)

	if reqNonceNum <= nonceNum {
		logger.ErrorLog("Duplicate Request.")
		return errors.New("Duplicate Request.")
	}

	db.Exec("UPDATE nonce SET nonce = '" + reqNonce + "' ")
	return nil
}

func validateHashKey(hash string, nonce string) error {
	unixTime := time.Now().Unix()
	this15m := unixTime / 60

	generateCompareHash := func(timeRef int64) string {
		digest := config.CONF.RpcPass + strconv.FormatInt(timeRef, 10) + config.CONF.RpcHashkey + nonce
		digestMd5 := md5.Sum([]byte(digest))
		digestMd5String := hex.EncodeToString(digestMd5[:])
		digestSha256 := sha256.Sum256([]byte(digestMd5String))
		digestSha256String := hex.EncodeToString(digestSha256[:])
		return digestSha256String
	}

	var cHashes [3]string
	cHashes[0] = generateCompareHash(this15m)
	cHashes[1] = generateCompareHash(this15m - 1)
	cHashes[2] = generateCompareHash(this15m + 1)

	for _, cHash := range cHashes {
		if hash == cHash {
			return nil
		}
	}

	logger.ErrorLog("Invalid Hashkey: " + hash)
	return errors.New("Invalid Hashkey.")
}
