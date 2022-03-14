package ergo

import (
	"encoding/json"
	"errors"
	"fmt"
	"strconv"

	"github.com/btcid/ergo-middleware-go/cmd/config"
	"github.com/btcid/ergo-middleware-go/pkg/lib/util"
	logger "github.com/btcid/ergo-middleware-go/pkg/logging"
	"github.com/go-resty/resty/v2"
)

type Balances struct {
	Height  int    `json:"height"`
	Balance *int64 `json:"balance"`
	Assets  struct {
	} `json:"assets"`
	Err
}

func GetBalance() (string, error) {
	var balance string

	response := Balances{}

	balancesWalletURL := fmt.Sprintf("%s/wallet/balances",
		config.CONF.NodeJsonHtppUrl)

	restyClient := resty.New()
	res, err := restyClient.SetCloseConnection(true).R().
		SetHeader("Content-Type", "application/json").
		SetHeader("api_key", config.CONF.NodeJsonHtppApiKey).
		Post(balancesWalletURL)

	if err != nil {
		logger.ErrorLog("GetBalance restyClient.R(). err: " + err.Error())
		return balance, err
	}

	unmarshalErr := json.Unmarshal(res.Body(), &response)
	if unmarshalErr != nil {
		logger.ErrorLog("GetBalance json.Unmarshal([]byte(res), &responseExp) err: " + unmarshalErr.Error())
		return balance, unmarshalErr
	}

	if (response.Error >= 300 && response.Error <= 600) && response.Detail != "" {
		logger.ErrorLog("GetBalance, err: " + response.Detail)
		return balance, errors.New(response.Detail)
	}

	if response.Balance == nil {
		return balance, errors.New("GetBalance unexpected error encountered in jsonrpc request")
	}

	balTemp := strconv.FormatInt(*response.Balance, 10)

	// max decimal = 9 (ergo)
	balance = util.RawToDecimal(balTemp, 9)

	return balance, nil
}

// 119994000000
// fee = 1000000
/*

curl command:
  curl -X 'GET' \
  'http://localhost:9052/wallet/balances' \
  -H 'accept: application/json' \
  -H 'api_key: hello'

	response:
	{
  "height" : 180823,
  "balance" : 119994000000,
  "assets" : {

  }
}
*/
