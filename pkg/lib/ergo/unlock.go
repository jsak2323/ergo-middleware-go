package ergo

import (
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"

	"github.com/go-resty/resty/v2"

	"github.com/btcid/ergo-middleware-go/cmd/config"
	"github.com/btcid/ergo-middleware-go/pkg/lib/util"
	logger "github.com/btcid/ergo-middleware-go/pkg/logging"
)

type unlockReq struct {
	Pass string `json:"pass"`
}

func UnlockWallet() error {

	response := Err{}

	encryptedPassBytes, _ := hex.DecodeString(config.CONF.EncryptedPassphrase)
	decryptedWalletPass, err := util.Decrypt(encryptedPassBytes, []byte(config.CONF.EncryptionKey))
	if err != nil {
		logger.ErrorLog("UnLockWallet util.Decrypt(encryptedPassBytes, []byte(config.CONF.EncryptionKey)) err: " + err.Error())
		return err
	}

	reqJson, err := json.Marshal(unlockReq{
		Pass: string(decryptedWalletPass),
	})
	if err != nil {
		logger.ErrorLog("UnLockWallet json.Marshal(req) err: " + err.Error())
		return err
	}

	unlockWalletURL := fmt.Sprintf("%s/wallet/unlock",
		config.CONF.NodeJsonHtppUrl)

	restyClient := resty.New()
	res, err := restyClient.SetCloseConnection(true).R().
		SetHeader("Content-Type", "application/json").
		SetHeader("api_key", config.CONF.NodeJsonHtppApiKey).
		SetBody(string(reqJson)).
		Post(unlockWalletURL)

	if err != nil {
		logger.ErrorLog("UnLockWallet restyClient.R(). err: " + err.Error())
		return err
	}

	if res.StatusCode() != 200 {
		unmarshalErr := json.Unmarshal(res.Body(), &response)
		if unmarshalErr != nil {
			logger.ErrorLog("UnlockWallet json.Unmarshal([]byte(res), &responseExp) err: " + unmarshalErr.Error())
			return unmarshalErr
		}

		if response.Detail != "Wallet already unlocked" {
			logger.ErrorLog("UnlockWallet, err: " + response.Detail)
			return errors.New(response.Detail)
		}
	}

	return nil
}

// curl -X 'POST' \
//   'http://localhost:9052/wallet/unlock' \
//   -H 'accept: application/json' \
//   -H 'api_key: hello' \
//   -H 'Content-Type: application/json' \
//   -d '{
//   "pass": "ergotest"
// }'
// "OK"
