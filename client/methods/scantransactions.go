package methods

import (
	"fmt"

	"github.com/btcid/ergo-middleware-go/cmd/config"
)

func ScanTransactions() {
	mysqlDbConn := config.MysqlDbConn()
	defer mysqlDbConn.Close()
	ergoCronService := NewErgoCron(mysqlDbConn)

	// delay in seconds until next execution
	delay := 31

	for {
		func() {
			fmt.Println("\n\n\n --- ScanBlock begin ---")

			defer HandlePanicAndCountdown(delay)

			// execute
			// ergoCronService.ScanTransactions(0)
			ergoCronService.ScanTransactionsV2(0)

			fmt.Println(" --- ScanBlock end ---")
		}()
	}
}
