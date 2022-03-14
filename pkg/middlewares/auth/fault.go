package auth

import (
	"net/http"
	"strconv"
	"strings"

	"github.com/btcid/ergo-middleware-go/cmd/config"
	"github.com/btcid/ergo-middleware-go/pkg/lib/util"
	logger "github.com/btcid/ergo-middleware-go/pkg/logging"
)

const emailSubjectPrefix string = "[DGB]"

func handleUnauthorizedIp(req *http.Request) {
	ip := strings.Split(req.RemoteAddr, ":")[0]

	logger.Log(" - AUTH -- Sending notification email ...")

	subject := emailSubjectPrefix + " Request from suspicious IP address: " + ip
	message := "A request from suspicious IP address was recorded with following detail: " +
		"\n IP Address: " + ip +
		"\n URL: " + req.URL.String()

	recipients := config.CONF.NotificationEmails

	isEmailSent, err := util.SendEmail(subject, message, recipients)
	if err != nil {
		logger.ErrorLog(err.Error())
	}
	logger.Log(" - AUTH -- Is unauthorized ip notification email sent: " + strconv.FormatBool(isEmailSent))
}

func handleUnauthorizedXmlRequest(req *http.Request, err error) {
	ip := strings.Split(req.RemoteAddr, ":")[0]

	logger.Log(" - AUTH -- Sending notification email ...")

	subject := emailSubjectPrefix + " Invalid XML RPC Request"
	message := "An invalid XML RPC request was recorded with following detail: " +
		"\n IP Address: " + ip +
		"\n URL: " + req.URL.String() +
		"\n Error: " + err.Error()

	recipients := config.CONF.NotificationEmails

	isEmailSent, err := util.SendEmail(subject, message, recipients)
	if err != nil {
		logger.ErrorLog(err.Error())
	}
	logger.Log(" - AUTH -- Is unauthorized xml request notification email sent: " + strconv.FormatBool(isEmailSent))
}
