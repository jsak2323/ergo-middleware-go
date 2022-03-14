package auth

import (
	"net/http"
	"strings"

	"github.com/btcid/ergo-middleware-go/cmd/config"
	"github.com/btcid/ergo-middleware-go/pkg/lib/util"
)

func isIpAuthorized(req *http.Request) bool {
	ip := strings.Split(req.RemoteAddr, ":")[0]

	if isAuthorized, _ := util.InArray(ip, config.CONF.AuthorizedIps); !isAuthorized {
		return false
	}
	return true
}
