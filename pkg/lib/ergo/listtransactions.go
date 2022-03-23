package ergo

import (
	"encoding/json"
	"errors"
	"fmt"

	"github.com/btcid/ergo-middleware-go/cmd/config"
	logger "github.com/btcid/ergo-middleware-go/pkg/logging"
	"github.com/go-resty/resty/v2"
)

type (
	ListTransactionsResp struct {
		Resp []ListTransactionResp
		Err
	}

	ListTransactionResp struct {
		ID     string `json:"id"`
		Inputs []struct {
			BoxID         string `json:"boxId"`
			SpendingProof struct {
				ProofBytes string `json:"proofBytes"`
				Extension  struct {
				} `json:"extension"`
			} `json:"spendingProof"`
		} `json:"inputs"`
		DataInputs []interface{} `json:"dataInputs"`
		Outputs    []struct {
			BoxID               string        `json:"boxId"`
			Value               int64         `json:"value"`
			ErgoTree            string        `json:"ergoTree"`
			Address             string        `json:"address"`
			Assets              []interface{} `json:"assets"`
			CreationHeight      int           `json:"creationHeight"`
			AdditionalRegisters struct {
			} `json:"additionalRegisters"`
		} `json:"outputs"`
		Size             int   `json:"size"`
		InclusionHeight  int64 `json:"inclusionHeight"`
		Scans            []int `json:"scans"`
		NumConfirmations int   `json:"numConfirmations"`
	}
)

func ListTransactions(minInclusionHeight, maxInclusionHeight int64) (ListTransactionsResp, error) {
	response := ListTransactionsResp{}

	err := UnlockWallet()
	if err != nil {
		logger.ErrorLog("ListTransactions unlock wallet. err: " + err.Error())
		return response, err
	}

	defer LockWallet()

	transactionsWalletURL := fmt.Sprintf("%s/wallet/transactions?minInclusionHeight=%v", config.CONF.NodeJsonHtppUrl, minInclusionHeight)

	if maxInclusionHeight > 0 {
		// by request block
		transactionsWalletURL += fmt.Sprintf("&maxInclusionHeight=%v", maxInclusionHeight)
	}

	restyClient := resty.New()
	res, err := restyClient.SetCloseConnection(true).R().
		SetHeader("Content-Type", "application/json").
		SetHeader("api_key", config.CONF.NodeJsonHtppApiKey).
		Get(transactionsWalletURL)

	if err != nil {
		logger.ErrorLog("ListTransactions restyClient.R(). err: " + err.Error())
		return response, err
	}

	unmarshalErr := json.Unmarshal(res.Body(), &response.Resp)
	if unmarshalErr != nil {
		logger.ErrorLog("ListTransactions json.Unmarshal([]byte(res), &responseExp) err: " + unmarshalErr.Error())
		return response, unmarshalErr
	}

	if res.StatusCode() != 200 {
		unmarshalErr := json.Unmarshal(res.Body(), &response.Error)
		if unmarshalErr != nil {
			logger.ErrorLog("ListTransactions json.Unmarshal([]byte(res), &responseExp) err: " + unmarshalErr.Error())
			return response, unmarshalErr
		}
		logger.ErrorLog("ListTransactions, err: " + response.Detail)
		return response, errors.New(response.Detail)
	}

	return response, nil
}

// /*

// curl command:

// curl -X 'GET' \
//   'https://editor.swagger.io/wallet/transactions?minConfirmations=0' \
//   -H 'accept: application/json' \
//   -H 'api_key: hello'

// example json response:
// [
//   {
//     "id": "2ab9da11fc216660e974842cc3b7705e62ebb9e0bf5ff78e53f9cd40abadd117",
//     "inputs": [
//       {
//         "boxId": "1ab9da11fc216660e974842cc3b7705e62ebb9e0bf5ff78e53f9cd40abadd117",
//         "spendingProof": {
//           "proofBytes": "4ab9da11fc216660e974842cc3b7705e62ebb9e0bf5ff78e53f9cd40abadd1173ab9da11fc216660e974842cc3b7705e62ebb9e0bf5ff78e53f9cd40abadd1173ab9da11fc216660e974842cc3b7705e62ebb9e0bf5ff78e53f9cd40abadd117",
//           "extension": {
//             "1": "a2aed72ff1b139f35d1ad2938cb44c9848a34d4dcfd6d8ab717ebde40a7304f2541cf628ffc8b5c496e6161eba3f169c6dd440704b1719e0"
//           }
//         }
//       }
//     ],
//     "dataInputs": [
//       {
//         "boxId": "1ab9da11fc216660e974842cc3b7705e62ebb9e0bf5ff78e53f9cd40abadd117"
//       }
//     ],
//     "outputs": [
//       {
//         "boxId": "1ab9da11fc216660e974842cc3b7705e62ebb9e0bf5ff78e53f9cd40abadd117",
//         "value": 147,
//         "ergoTree": "0008cd0336100ef59ced80ba5f89c4178ebd57b6c1dd0f3d135ee1db9f62fc634d637041",
//         "creationHeight": 9149,
//         "assets": [
//           {
//             "tokenId": "4ab9da11fc216660e974842cc3b7705e62ebb9e0bf5ff78e53f9cd40abadd117",
//             "amount": 1000
//           }
//         ],
//         "additionalRegisters": {
//           "R4": "100204a00b08cd0336100ef59ced80ba5f89c4178ebd57b6c1dd0f3d135ee1db9f62fc634d637041ea02d192a39a8cc7a70173007301"
//         },
//         "transactionId": "2ab9da11fc216660e974842cc3b7705e62ebb9e0bf5ff78e53f9cd40abadd117",
//         "index": 0
//       }
//     ],
//     "inclusionHeight": 20998,
//     "numConfirmations": 20998,
//     "scans": [
//       1
//     ],
//     "size": 0
//   }
// ]

// */
