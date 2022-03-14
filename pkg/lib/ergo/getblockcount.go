package ergo

import (
	"encoding/json"
	"errors"
	"fmt"

	"github.com/btcid/ergo-middleware-go/cmd/config"
	logger "github.com/btcid/ergo-middleware-go/pkg/logging"
	"github.com/go-resty/resty/v2"
)

type InfoResponse struct {
	CurrentTime          int64  `json:"currentTime"`
	Network              string `json:"network"`
	Name                 string `json:"name"`
	StateType            string `json:"stateType"`
	Difficulty           int64  `json:"difficulty"`
	BestFullHeaderID     string `json:"bestFullHeaderId"`
	BestHeaderID         string `json:"bestHeaderId"`
	PeersCount           int    `json:"peersCount"`
	UnconfirmedCount     int    `json:"unconfirmedCount"`
	AppVersion           string `json:"appVersion"`
	StateRoot            string `json:"stateRoot"`
	GenesisBlockID       string `json:"genesisBlockId"`
	PreviousFullHeaderID string `json:"previousFullHeaderId"`
	FullHeight           *int   `json:"fullHeight"`
	HeadersHeight        *int   `json:"headersHeight"`
	StateVersion         string `json:"stateVersion"`
	FullBlocksScore      int64  `json:"fullBlocksScore"`
	LaunchTime           int64  `json:"launchTime"`
	LastSeenMessageTime  int64  `json:"lastSeenMessageTime"`
	HeadersScore         int64  `json:"headersScore"`
	Parameters           struct {
		OutputCost       int `json:"outputCost"`
		TokenAccessCost  int `json:"tokenAccessCost"`
		MaxBlockCost     int `json:"maxBlockCost"`
		Height           int `json:"height"`
		MaxBlockSize     int `json:"maxBlockSize"`
		DataInputCost    int `json:"dataInputCost"`
		BlockVersion     int `json:"blockVersion"`
		InputCost        int `json:"inputCost"`
		StorageFeeFactor int `json:"storageFeeFactor"`
		MinValuePerByte  int `json:"minValuePerByte"`
	} `json:"parameters"`
	IsMining bool `json:"isMining"`
	Err
}

// headers height
func GetBlockCount() (int, error) {
	blockCount := 0

	response := InfoResponse{}

	infoURL := fmt.Sprintf("%s/info",
		config.CONF.NodeJsonHtppUrl)

	restyClient := resty.New()
	res, err := restyClient.SetCloseConnection(true).R().
		SetHeader("Content-Type", "application/json").
		Get(infoURL)

	if err != nil {
		logger.ErrorLog("GetBlockCount restyClient.R(). err: " + err.Error())
		return blockCount, err
	}

	err = json.Unmarshal(res.Body(), &response)
	if err != nil {
		logger.ErrorLog("GetBlockCount json.Unmarshal([]byte(res), &responseExp) err: " + err.Error())
		return blockCount, err
	}

	if (response.Error >= 300 && response.Error <= 600) && response.Detail != "" {
		logger.ErrorLog("GetBlockCount, err: " + response.Detail)
		return blockCount, errors.New(response.Detail)
	}

	if response.HeadersHeight == nil {
		logger.ErrorLog("GetBlockCount, err: node haven't start downloading full blocks")
		return blockCount, errors.New("node haven't start downloading full blocks")
	}
	blockCount = *response.HeadersHeight

	return blockCount, nil

}

/*
headersHeight = total block
curl -X 'GET' \
  'http://localhost:9052/info' \
  -H 'accept: application/json'

*/
