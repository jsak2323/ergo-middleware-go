package auth

import (
	"net/http"

	logger "github.com/btcid/ergo-middleware-go/pkg/logging"
)

func AuthMiddleware(hf http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {

		// AUTHORIZE IP
		if isIpAuthorized := isIpAuthorized(req); isIpAuthorized { // if ip is authorized, continue
			logger.InfoLog(" - AUTH -- IP is authorized.", req)

		} else { // if not authorized, send notification email and stop request
			logger.InfoLog(" - AUTH -- IP is unauthorized.", req)
			handleUnauthorizedIp(req)
			return
		}

		// AUTHORIZE XML RPC REQUEST
		if req.URL.String() == "/xmlrpc" {
			if err := AuthorizeXmlRequest(req); err == nil { // if xml req is authorized, continue
				logger.InfoLog(" - AUTH -- XML Request is authorized.", req)

			} else { // if not authorized, send notification email and stop request
				logger.InfoLog(" - AUTH -- XML Request is unauthorized. err: "+err.Error(), req)
				handleUnauthorizedXmlRequest(req, err)
				return
			}
		}

		hf.ServeHTTP(w, req) // serve the request

	})
}
