package ergo

import (
	"encoding/json"
	"errors"
	"fmt"

	"github.com/btcid/ergo-middleware-go/cmd/config"
	logger "github.com/btcid/ergo-middleware-go/pkg/logging"
	"github.com/go-resty/resty/v2"
)

type GetNewAddressResp struct {
	Address string `json:"address"`
	IsValid bool   `json:"isValid"`
	Err
}

func GetNewAddress() (GetNewAddressResp, error) {
	response := GetNewAddressResp{}

	err := UnlockWallet()
	if err != nil {
		logger.ErrorLog("GetNewAddress unlock wallet. err: " + err.Error())
		return response, err
	}

	defer LockWallet()

	GetNewAddressURL := fmt.Sprintf("%s/wallet/deriveNextKey",
		config.CONF.NodeJsonHtppUrl,
	)

	restyClient := resty.New()
	res, err := restyClient.SetCloseConnection(true).R().
		SetHeader("Content-Type", "application/json").
		SetHeader("api_key", config.CONF.NodeJsonHtppApiKey).
		Get(GetNewAddressURL)

	if err != nil {
		logger.ErrorLog("GetNewAddress restyClient.R(). err: " + err.Error())
		return response, err
	}

	unmarshalErr := json.Unmarshal(res.Body(), &response)
	if unmarshalErr != nil {
		logger.ErrorLog("GetNewAddress json.Unmarshal([]byte(res), &responseExp) err: " + unmarshalErr.Error())
		return response, unmarshalErr
	}

	if res.StatusCode() != 200 && response.Detail != "" {
		logger.ErrorLog("GetNewAddress, err: " + response.Detail)
		return response, errors.New(response.Detail)
	}

	// err = LockWallet()
	// if err != nil {
	// 	logger.ErrorLog("GetNewAddress lock wallet. err: " + err.Error())
	// 	return response, err
	// }

	return response, nil
}

/*

curl command:
 curl -X 'GET' \
  'http://localhost:9052/wallet/deriveNextKey' \
  -H 'accept: application/json' \
  -H  "api_key: hello"

example json response:
{
  "result": "RKSMvLG7PnuNbJ6JPKfsvec4HKtBTeRAQB",
  "error": null,
  "id": "curltest"
}

*/
