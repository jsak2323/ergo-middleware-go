package ergo

import (
	"encoding/json"
	"errors"
	"fmt"
	"math/big"

	"github.com/btcid/ergo-middleware-go/cmd/config"
	"github.com/btcid/ergo-middleware-go/pkg/lib/util"
	logger "github.com/btcid/ergo-middleware-go/pkg/logging"
	"github.com/go-resty/resty/v2"
)

type (
	paymentSendReq []struct {
		Address string `json:"address"`
		Value   int64  `json:"value"`
	}
)

func SendToAddress(amountInDecimal string, address string) (string, error) {
	var txHash string

	balance := util.DecimalToRaw(amountInDecimal, 9)
	amount, ok := new(big.Int).SetString(balance, 10)
	if !ok {
		return txHash, errors.New("fail big.SetString(" + amountInDecimal + ")")
	}

	reqBody := paymentSendReq{
		{
			Address: address,
			Value:   amount.Int64(),
		},
	}

	fmt.Println("address ", address)
	fmt.Println("amount req ", amountInDecimal)
	fmt.Println("balance ", balance)
	fmt.Println("amount ", amount)
	fmt.Println("amount.INT64 ", amount.Int64())

	reqJson, err := json.Marshal(reqBody)
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
		Post(paymentWalletURL)

	if err != nil {
		logger.ErrorLog("SendToAddress restyClient.R(). err: " + err.Error())
		return txHash, err
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

	unmarshalErr := json.Unmarshal(res.Body(), &txHash)
	if unmarshalErr != nil {
		logger.ErrorLog("SendToAddress json.Unmarshal([]byte(res), &responseExp) err: " + unmarshalErr.Error())
		return txHash, unmarshalErr
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
  'http://localhost:9052/wallet/payment/send' \
  -H 'accept: application/json' \
  -H 'api_key: hello' \
  -H 'Content-Type: application/json' \
  -d '[
  {
    "address": "3WyjzUgxpHC1iEx12CvRNtEGrKT7Qw9mLTcxCGYo4QbWg9m4Cv4N",
    "value": 100000000
  }
]'

example json response:
"b48a63a28eae63634866efcfd0c9d60921d1412adbdcd9e8fec2805c3ad61658"

*/
