package logging

import (
	"fmt"
	"net/http"

	logger "github.com/btcid/ergo-middleware-go/pkg/logging"
)

func LogMiddleware(hf http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		logger.InfoLog(req.URL.String()+" hit. ", req)
		hf.ServeHTTP(w, req)
		logger.InfoLog(req.URL.String()+" done. ", req)
		fmt.Println()
	})
}
