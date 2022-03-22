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
	GetTransactionByIdResp struct {
		Resp TransactionByIdResp
		Err
	}

	TransactionByIdResp struct {
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

func GetTransactionById(idTx string) (GetTransactionByIdResp, error) {
	response := GetTransactionByIdResp{}

	if idTx == "" {
		logger.ErrorLog("GetTransactionById validation idTx err: empty id tx")
		return response, errors.New("err: empty id tx")
	}

	transactionsWalletURL := fmt.Sprintf("%s/wallet/transactionById?id=", idTx)

	restyClient := resty.New()
	res, err := restyClient.SetCloseConnection(true).R().
		SetHeader("Content-Type", "application/json").
		SetHeader("api_key", config.CONF.NodeJsonHtppApiKey).
		Get(transactionsWalletURL)

	if err != nil {
		logger.ErrorLog("GetTransactionById restyClient.R(). err: " + err.Error())
		return response, err
	}

	unmarshalErr := json.Unmarshal(res.Body(), &response)
	if unmarshalErr != nil {
		logger.ErrorLog("GetTransactionById json.Unmarshal([]byte(res), &responseExp) err: " + unmarshalErr.Error())
		return response, unmarshalErr
	}

	if res.StatusCode() != 200 {
		logger.ErrorLog("GetTransactionById, err: " + response.Detail)
		return response, errors.New(response.Detail)
	}

	return response, nil
}

// /*

// curl command:

// curl -X 'GET' \
//   'http://localhost:9052/wallet/transactionById?id=83ecff3210f9ac37768ae8f0799d5cc4c3a59ff4955e622c99e40108ace1e667' \
//   -H 'accept: application/json' \
//   -H 'api_key: hello'

// example json response:
// {
// 	"id" : "87d6dc51279fab961fb754d202c6ea7977593ce7b502b227286c2865a94f9d59",
// 	"inputs" : [
// 			{
// 					"boxId" : "36dee93684b9cf75f1959dd35d16956217511b52639a183680ed27c69931885c",
// 					"spendingProof" : {
// 							"proofBytes" : "37f4b0e921aef019e3fd23b85a874a562d555678eec1a7c2ca3e44611202fcb33df9887914ac727192cc96ec29ef61293e809829eacca4b5",
// 							"extension" : {

// 							}
// 					}
// 			}
// 	],
// 	"dataInputs" : [
// 	],
// 	"outputs" : [
// 			{
// 					"boxId" : "9d13c08858994fffd23334870aaf1d56ceb43ae9d93318aff069c892ee8bdcdd",
// 					"value" : 2100000000,
// 					"ergoTree" : "0008cd0299cfe61ca153109cd33b69ccad6f5f780e730c1e60fa5f0ebab00b7a27a41314",
// 					"address" : "3Wwmx2AoP5MpLApcpHkdpjdscFcBStSiYo2jQsyYdEc2MaoGaU59",
// 					"assets" : [
// 					],
// 					"creationHeight" : 190589,
// 					"additionalRegisters" : {

// 					}
// 			},
// 			{
// 					"boxId" : "cb69781cd870c7e824214676932a420cff5fd9787d2497ad8ba4accecbc9dfbd",
// 					"value" : 3100000000,
// 					"ergoTree" : "0008cd034260e587c6d647ee40603b1314758fe63d3d447ca006ac4250ac1319d564bef2",
// 					"address" : "3Wy4Bpi1YD74w1CfFNxpoKeGNPSHdQgkz3qxeNuAVahxCjDFQ5Bb",
// 					"assets" : [
// 					],
// 					"creationHeight" : 190589,
// 					"additionalRegisters" : {

// 					}
// 			},
// 			{
// 					"boxId" : "2894c70e970a0a0039d836a8ba9e7f0bf3af088842720653ce1c4fed02a9df39",
// 					"value" : 1000000,
// 					"ergoTree" : "1005040004000e36100204900108cd0279be667ef9dcbbac55a06295ce870b07029bfcdb2dce28d959f2815b16f81798ea02d192a39a8cc7a701730073011001020402d19683030193a38cc7b2a57300000193c2b2a57301007473027303830108cdeeac93b1a57304",
// 					"address" : "Bf1X9JgQTUtgntaepHNF5tbnwto2VpaRo1NkEq3SowNF6Mpv9R9g5kYmRpCXxBkuoZCMmyZNTpVMNyjKJKQ3VYf8T3JyZ724at6VGi6aUq1VyyeHucC7hxnVmCeDCZ3aHtrzVcrCjtqmE6LrPuZ9n3",
// 					"assets" : [
// 					],
// 					"creationHeight" : 190589,
// 					"additionalRegisters" : {

// 					}
// 			},
// 			{
// 					"boxId" : "74b4381d537393cae128ffab5870ef397fcf38b69b0848d1f4811efa85ff946e",
// 					"value" : 54799000000,
// 					"ergoTree" : "0008cd03bf8ae845ab034eb2e639d550bb51fcadec8ce1dc92753845a0368b6c47bb1e64",
// 					"address" : "3Wz1JyMSQBUHX8YADpVJfx1VJvhH11nHmy5mHWcKasmtAj6Vx5h2",
// 					"assets" : [
// 					],
// 					"creationHeight" : 190589,
// 					"additionalRegisters" : {

// 					}
// 			}
// 	],
// 	"size" : 346,
// 	"inclusionHeight" : 190591,
// 	"scans" : [
// 			10
// 	],
// 	"numConfirmations" : 78
// }
