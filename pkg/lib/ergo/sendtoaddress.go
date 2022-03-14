package ergo

import (
	"encoding/json"
	"errors"
	"fmt"
	"math/big"

	"github.com/btcid/ergo-middleware-go/cmd/config"
	logger "github.com/btcid/ergo-middleware-go/pkg/logging"
	"github.com/go-resty/resty/v2"
)

type (
	re1 string

	paymentSendReq struct {
		Address string `json:"address"`
		Value   int64  `json:"value"`
	}
)

func SendToAddress(address string, amountInDecimal string) (string, error) {
	txHash := ""

	err := UnlockWallet()
	if err != nil {
		logger.ErrorLog("SendToAddress unlock wallet. err: " + err.Error())
		return txHash, err
	}

	defer LockWallet()

	amount, ok := new(big.Int).SetString(amountInDecimal, 10)
	if !ok {
		return txHash, errors.New("fail big.SetString(" + amountInDecimal + ")")
	}

	reqJson, err := json.Marshal(paymentSendReq{
		Address: address,
		Value:   amount.Int64(),
	})
	if err != nil {
		logger.ErrorLog("SendToAddress json.Marshal(req) err: " + err.Error())
		return txHash, err
	}

	paymentWalletURL := fmt.Sprintf("%s/wallet/payment/send",
		config.CONF.NodeJsonHtppUrl)

	restyClient := resty.New()
	res, err := restyClient.SetCloseConnection(true).R().
		SetHeader("Content-Type", "application/json").
		SetHeader("api_key", config.CONF.NodeJsonHtppApiKey).
		SetBody(string(reqJson)).
		Get(paymentWalletURL)

	if err != nil {
		logger.ErrorLog("SendToAddress restyClient.R(). err: " + err.Error())
		return txHash, err
	}

	unmarshalErr := json.Unmarshal(res.Body(), &txHash)
	if unmarshalErr != nil {
		logger.ErrorLog("SendToAddress json.Unmarshal([]byte(res), &responseExp) err: " + unmarshalErr.Error())
		return txHash, unmarshalErr
	}

	if res.StatusCode() != 200 {
		var err Err
		unmarshalErr := json.Unmarshal(res.Body(), &err)
		if unmarshalErr != nil {
			logger.ErrorLog("SendToAddress json.Unmarshal([]byte(res), &responseExp) err: " + unmarshalErr.Error())
			return txHash, unmarshalErr
		}

		logger.ErrorLog("SendToAddress, err: " + err.Detail)
		return txHash, errors.New(err.Detail)
	}

	// err = LockWallet()
	// if err != nil {
	// 	logger.ErrorLog("GetNewAddress lock wallet. err: " + err.Error())
	// 	return response, err
	// }

	return txHash, nil
}

/*

curl command:
curl -X 'POST' \
  'https://editor.swagger.io/wallet/payment/send' \
  -H 'accept: application/json' \
  -H 'api_key: hello' \
  -H 'Content-Type: application/json' \
  -d '[
  {
    "address": "3WwbzW6u8hKWBcL1W7kNVMr25s2UHfSBnYtwSHvrRQt7DdPuoXrt",
    "value": 1
  }
]'

example json response:
"b48a63a28eae63634866efcfd0c9d60921d1412adbdcd9e8fec2805c3ad61658"

*/
