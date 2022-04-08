package methods

import (
	"database/sql"
	"fmt"
	"time"

	mysqldb "github.com/btcid/ergo-middleware-go/pkg/database/mysql"
	httpcron "github.com/btcid/ergo-middleware-go/pkg/http/cron"
)

func NewErgoCron(mysqlDbConn *sql.DB) *httpcron.ErgoCron {
	transactionRepo := mysqldb.NewMysqlTransactionRepository(mysqlDbConn)
	addressRepo := mysqldb.NewMysqlAddressRepository(mysqlDbConn)
	blocksRepo := mysqldb.NewMysqlBlocksRepository(mysqlDbConn)
	ergoCronService := httpcron.NewErgoCron(addressRepo, transactionRepo, blocksRepo)

	return ergoCronService
}

func HandlePanicAndCountdown(delay int) {
	ticker := time.Tick(time.Second)

	if err := recover(); err != nil {
		fmt.Println(" - panic: ", err)
	}

	for i := delay; i >= 0; i-- {
		<-ticker
		fmt.Printf("\r - Next execution in %d seconds ...", i)
	}
}
