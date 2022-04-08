package methods

import (
	"fmt"

	"github.com/btcid/ergo-middleware-go/cmd/config"
)

func UpdateConfirmations() {
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
			ergoCronService.UpdateConfirmations()

			fmt.Println(" --- ScanBlock end ---")
		}()
	}
}
