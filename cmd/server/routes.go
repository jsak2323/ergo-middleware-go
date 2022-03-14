package main

import (
	"github.com/divan/gorilla-xmlrpc/xml"
	"github.com/gorilla/mux"
	"github.com/gorilla/rpc"

	httphandler "github.com/btcid/ergo-middleware-go/pkg/http"
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

func SetRoutes(r *mux.Router) {

	// XMLRPC SERVICE
	xmlCodec := xml.NewCodec()
	ergoXmlRpcService := httprpc.NewERGORpc()
	ergoXmlRpcServer := rpc.NewServer()
	ergoXmlRpcServer.RegisterCodec(xmlCodec, "text/xml")
	ergoXmlRpcServer.RegisterBeforeFunc(rpcBeforeFunc)
	ergoXmlRpcServer.RegisterService(ergoXmlRpcService, "")
	ergoXmlRpcServer.RegisterAfterFunc(rpcAfterFunc)
	r.Handle("/xmlrpc", ergoXmlRpcServer)

	// CRON ROUTES
	// ergoCronService := httpcron.NewergoCron()

	// r.HandleFunc("/cron/collect_dust", ergoCronService.CollectDust)

	// HTTP ROUTES
	r.HandleFunc("/log/{date}", httphandler.GetLog)

}
