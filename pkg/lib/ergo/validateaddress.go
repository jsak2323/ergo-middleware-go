package ergo

import (
	"encoding/json"
	"errors"
	"fmt"

	"github.com/btcid/ergo-middleware-go/cmd/config"
	logger "github.com/btcid/ergo-middleware-go/pkg/logging"
	"github.com/go-resty/resty/v2"
)

type ValidateAddressRes struct {
	IsValid bool `json:"isvalid"`
	Err
}

func ValidateAddress(address string) (ValidateAddressRes, error) {
	response := ValidateAddressRes{}

	validateAddressURL := fmt.Sprintf("%s/utils/address/%s",
		config.CONF.NodeJsonHtppUrl,
		address,
	)

	restyClient := resty.New()
	res, err := restyClient.SetCloseConnection(true).R().
		SetHeader("Content-Type", "application/json").
		Get(validateAddressURL)

	if err != nil {
		logger.ErrorLog("ValidateAddress restyClient.R(). err: " + err.Error())
		return response, err
	}

	unmarshalErr := json.Unmarshal(res.Body(), &response)
	if unmarshalErr != nil {
		logger.ErrorLog("ValidateAddress json.Unmarshal([]byte(res), &responseExp) err: " + unmarshalErr.Error())
		return response, unmarshalErr
	}

	if (response.Error >= 300 && response.Error <= 600) && response.Detail != "" {
		logger.ErrorLog("LockWallet, err: " + response.Detail)
		return response, errors.New(response.Detail)
	}

	return response, nil
}

/*

curl command:
curl -X 'GET' \
  'http://localhost:9052/utils/address/{address}' \
  -H 'accept: application/json'

example json response:
{
  "address" : "AfYgQf5PappexKq8Vpig4vwEuZLjrq7gV97BWBVcKymTYqRzCoJLE9cDBpGHvtAAkAgQf8Yyv7NQUjSphKSjYxk3dB3W8VXzHzz5MuCcNbqqKHnMDZAa6dbHH1uyMScq5rXPLFD5P8MWkD5FGE6RbHKrKjANcr6QZHcBpppdjh9r5nra4c7dsCgULFZfWYTaYqHpx646BUHhhp8jDCHzzF33G8XfgKYo93ABqmdqagbYRzrqCgPHv5kxRmFt7Y99z26VQTgXoEmXJ2aRu6LoB59rKN47JxWGos27D79kKzJRiyYNEVzXU8MYCxtAwV",
  "isValid" : true
}


*/
