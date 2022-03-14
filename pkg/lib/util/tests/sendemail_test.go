package tests

import (
	"testing"

	"github.com/btcid/ergo-middleware-go/cmd/config"
	"github.com/btcid/ergo-middleware-go/pkg/lib/util"
)

func TestSendEmail(t *testing.T) {
	const emailSubjectPrefix string = "[TEST]"
	subject := emailSubjectPrefix + " Test Send Email"
	message := "Test message"

	recipients := config.CONF.NotificationEmails

	if _, err := util.SendEmail(subject, message, recipients); err != nil {
		t.Errorf("TestSendEmail failed; err: %s", err.Error())
	}

}
