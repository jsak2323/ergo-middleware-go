package ergo

import (
	"encoding/json"
	"errors"
	"fmt"

	"github.com/go-resty/resty/v2"

	"github.com/btcid/ergo-middleware-go/cmd/config"
	logger "github.com/btcid/ergo-middleware-go/pkg/logging"
)

func LockWallet() error {

	response := Err{}

	lockWalletURL := fmt.Sprintf("%s/wallet/lock",
		config.CONF.NodeJsonHtppUrl)

	restyClient := resty.New()
	res, err := restyClient.SetCloseConnection(true).R().
		SetHeader("Content-Type", "application/json").
		SetHeader("api_key", config.CONF.NodeJsonHtppApiKey).
		Post(lockWalletURL)

	if err != nil {
		logger.ErrorLog("LockWallet restyClient.R(). err: " + err.Error())
		return err
	}

	unmarshalErr := json.Unmarshal(res.Body(), &response)
	if unmarshalErr != nil {
		logger.ErrorLog("LockWallet json.Unmarshal([]byte(res), &responseExp) err: " + unmarshalErr.Error())
		return unmarshalErr
	}

	if (response.Error >= 300 && response.Error <= 600) && response.Detail != "" {
		logger.ErrorLog("LockWallet, err: " + response.Detail)
		return errors.New(response.Detail)
	}

	return nil
}

// curl -X 'GET' \
//   'http://localhost:9052/wallet/lock' \
//   -H 'accept: application/json' \
//   -H 'api_key: hello'
