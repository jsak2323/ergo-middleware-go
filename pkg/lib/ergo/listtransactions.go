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
					Num1 string `json:"1"`
				} `json:"extension"`
			} `json:"spendingProof"`
		} `json:"inputs"`
		DataInputs []struct {
			BoxID string `json:"boxId"`
		} `json:"dataInputs"`
		Outputs []struct {
			BoxID          string `json:"boxId"`
			Value          int    `json:"value"`
			ErgoTree       string `json:"ergoTree"`
			CreationHeight int    `json:"creationHeight"`
			Assets         []struct {
				TokenID string `json:"tokenId"`
				Amount  int    `json:"amount"`
			} `json:"assets"`
			AdditionalRegisters struct {
				R4 string `json:"R4"`
			} `json:"additionalRegisters"`
			TransactionID string `json:"transactionId"`
			Index         int    `json:"index"`
		} `json:"outputs"`
		InclusionHeight  int   `json:"inclusionHeight"`
		NumConfirmations int   `json:"numConfirmations"`
		Scans            []int `json:"scans"`
		Size             int   `json:"size"`
	}
)

func ListTransactions(minConfirmations int) (ListTransactionsResp, error) {
	response := ListTransactionsResp{}

	err := UnlockWallet()
	if err != nil {
		logger.ErrorLog("GetNewAddress unlock wallet. err: " + err.Error())
		return response, err
	}

	defer LockWallet()

	transactionsWalletURL := fmt.Sprintf("%s/wallet/transactions",
		config.CONF.NodeJsonHtppUrl)

	if minConfirmations != 0 {
		transactionsWalletURL += fmt.Sprintf("?minConfirmations=%v", minConfirmations)
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

	unmarshalErr := json.Unmarshal(res.Body(), &response)
	if unmarshalErr != nil {
		logger.ErrorLog("ListTransactions json.Unmarshal([]byte(res), &responseExp) err: " + unmarshalErr.Error())
		return response, unmarshalErr
	}

	if (response.Error >= 300 && response.Error <= 600) && response.Detail != "" {
		logger.ErrorLog("ListTransactions, err: " + response.Detail)
		return response, errors.New(response.Detail)
	}

	// err = LockWallet()
	// if err != nil {
	// 	logger.ErrorLog("GetNewAddress lock wallet. err: " + err.Error())
	// 	return response, err
	// }

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
