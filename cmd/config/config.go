package config

import (
	"encoding/json"
	"fmt"
	"os"
)

var (
	IS_DEV         bool
	CONF           Configuration
	ErrorMailCount int
)

func init() {
	IS_DEV = os.Getenv("PRODUCTION") != "true"

	fmt.Println()
	env := "development"
	if !IS_DEV {
		env = "production"
	}
	fmt.Println("Environment: " + env)

	LoadAppConfig()
}

type Configuration struct {
	RpcUser    string `json:"rpc_user"`
	RpcPass    string `json:"rpc_pass"`
	RpcPort    string `json:"rpc_port"`
	RpcHashkey string `json:"rpc_hashkey"`

	NodeJsonRpcUrl       string `json:"node_jsonrpc_url"`
	NodeRpcUser          string `json:"node_rpc_user"`
	NodeEncryptedRpcPass string `json:"node_encrypted_rpc_pass"`

	NodeJsonHtppUrl    string `json:"node_jsonhttp_url"`
	NodeJsonHtppApiKey string `json:"node_jsonhttp_api_key"`

	WalletPassword string `json:"wallet_pass"`

	SecondNodeJsonRpcUrl       string `json:"second_node_jsonrpc_url"`
	SecondNodeRpcUser          string `json:"second_node_rpc_user"`
	SecondNodeEncryptedRpcPass string `json:"second_node_encrypted_rpc_pass"`
	SecondNodeEnable           bool   `json:"second_node_enable"`

	MainAddress string  `json:"main_address"`
	FeeDefault  float64 `json:"fee_default"`
	MinClearing float64 `json:"min_clearing"`

	EncryptedPassphrase string `json:"encrypted_passphrase"`
	EncryptionKey       string `json:"encryption_key"`

	NotificationEmails []string `json:"notification_emails"`

	AuthorizedIps []string `json:"authorized_ips"`

	MailHost          string `json:"mail_host"`
	MailPort          string `json:"mail_port"`
	MailUser          string `json:"mail_user"`
	MailAddress       string `json:"mail_address"`
	MailEncryptedPass string `json:"mail_encrypted_pass"`
	MailEncryptionKey string `json:"mail_encryption_key"`

	SessionErrorMailNotifLimit int `json:"session_error_mail_notif_limit"`
}

func LoadAppConfig() {
	configFilename := "config.json"
	if IS_DEV {
		configFilename = "config-dev.json"
	}

	fmt.Print("Loading App Configuration ... ")
	gopath := os.Getenv("GOPATH")
	file, _ := os.Open(gopath + "/src/github.com/btcid/ergo-middleware-go/cmd/config/json/" + configFilename)
	defer file.Close()
	decoder := json.NewDecoder(file)
	err := decoder.Decode(&CONF)
	if err != nil {
		fmt.Println("error:", err)
	}
	fmt.Println("Done.")
}
