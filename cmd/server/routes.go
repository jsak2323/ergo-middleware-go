package main

import (
	"database/sql"

	"github.com/divan/gorilla-xmlrpc/xml"
	"github.com/gorilla/mux"
	"github.com/gorilla/rpc"

	mysqldb "github.com/btcid/ergo-middleware-go/pkg/database/mysql"
	httphandler "github.com/btcid/ergo-middleware-go/pkg/http"
	httpcron "github.com/btcid/ergo-middleware-go/pkg/http/cron"
	httprpc "github.com/btcid/ergo-middleware-go/pkg/http/rpc"
	logger "github.com/btcid/ergo-middleware-go/pkg/logging"
)

func rpcBeforeFunc(ri *rpc.RequestInfo) {
	var req = ri.Request
	logger.InfoLog(ri.Method+" hit.--------------", req)
}

func rpcAfterFunc(ri *rpc.RequestInfo) {
	var req = ri.Request
	logger.InfoLog(ri.Method+" done.--------------", req)
}

func SetRoutes(r *mux.Router, mysqlDbConn *sql.DB) {
	// REPOSITORIES
	addressRepo := mysqldb.NewMysqlAddressRepository(mysqlDbConn)
	transactionRepo := mysqldb.NewMysqlTransactionRepository(mysqlDbConn)
	blocksRepo := mysqldb.NewMysqlBlocksRepository(mysqlDbConn)

	// XMLRPC SERVICE
	xmlCodec := xml.NewCodec()
	ergoXmlRpcService := httprpc.NewERGORpc(addressRepo, transactionRepo, blocksRepo)
	ErgoXmlRpcServer := rpc.NewServer()
	ErgoXmlRpcServer.RegisterCodec(xmlCodec, "text/xml")
	ErgoXmlRpcServer.RegisterBeforeFunc(rpcBeforeFunc)
	ErgoXmlRpcServer.RegisterService(ergoXmlRpcService, "")
	ErgoXmlRpcServer.RegisterAfterFunc(rpcAfterFunc)
	r.Handle("/xmlrpc", ErgoXmlRpcServer)

	// CRON ROUTES
	ergoCronService := httpcron.NewErgoCron(addressRepo, transactionRepo, blocksRepo)

	_ = ergoCronService
	r.HandleFunc("/cron/scan_transactions/{blocknum}", ergoCronService.ScanBlockAndUpdateTransactions)
	r.HandleFunc("/cron/update_confirmations", ergoCronService.UpdateConfirmations)
	// r.HandleFunc("/cron/collect_dust", ergoCronService.CollectDust)

	// HTTP ROUTES
	r.HandleFunc("/log/{date}", httphandler.GetLog)

}
