package main

import (
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/gorilla/mux"

	"github.com/btcid/ergo-middleware-go/cmd/config"
	authm "github.com/btcid/ergo-middleware-go/pkg/middlewares/auth"
	logm "github.com/btcid/ergo-middleware-go/pkg/middlewares/logging"
)

func main() {
	r := mux.NewRouter()

	SetRoutes(r)

	r.Use(logm.LogMiddleware)
	r.Use(authm.AuthMiddleware)

	server := &http.Server{
		Handler:      r,
		Addr:         ":" + config.CONF.RpcPort,
		WriteTimeout: 120 * time.Second,
		ReadTimeout:  120 * time.Second,
	}

	fmt.Println()
	log.Println("Running Server on localhost:" + config.CONF.RpcPort)
	fmt.Print("\n\n\n")

	log.Fatal(server.ListenAndServe())
}
