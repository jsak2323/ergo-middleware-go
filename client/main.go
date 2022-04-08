package main

import (
	"fmt"
	"os"

	"github.com/btcid/ergo-middleware-go/client/methods"
)

func main() {
	fmt.Print("\n\n")
	if len(os.Args) <= 1 {
		fmt.Print(" - Please specify a method.\n\n\n")
		return
	}

	method := os.Args[1]
	switch method {
	case "getblockcount":
		methods.GetBlockCount()

	case "getnewaddress":
		methods.GetNewAddress()

	case "validateaddress":
		address := ""
		if len(os.Args) >= 3 {
			address = os.Args[2] // e.g: qQ91tYLjNe58b8GFHwhFJNiup9qgvfy7PG
		}
		methods.ValidateAddress(address)

	case "getbalance":
		methods.GetBalance()

	case "scantransactions":
		methods.ScanTransactions()

	case "updateconfirmations":
		methods.UpdateConfirmations()

	case "listtransactions":
		limit := ""
		if len(os.Args) >= 3 {
			limit = os.Args[2] // e.g: 5
		}
		methods.ListTransactions(limit)

	case "sendtoaddress":
		if len(os.Args) < 3 {
			fmt.Print(" - Invalid Parameters.\n\n\n")
			return
		}
		amount := os.Args[2] // e.g: 0.01
		to := os.Args[3]     // e.g: RNJz9x7qBkpxxqzKdZryP3SCShnCuDu6og

		methods.SendToAddress(amount, to)

	default:
		fmt.Print(" - Method not found.\n")
	}

	fmt.Print("\n\n")
}
