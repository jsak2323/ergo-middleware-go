package mysql

import (
	"database/sql"
	"fmt"
	"strconv"

	_ "github.com/go-sql-driver/mysql"

	tx "github.com/btcid/ergo-middleware-go/pkg/domain/transaction"
)

const transactionsTable = "transactions"

type transactionRepository struct {
	db *sql.DB
}

func NewMysqlTransactionRepository(db *sql.DB) tx.TransactionRepository {
	return &transactionRepository{
		db,
	}
}

func (r *transactionRepository) Create(transaction *tx.Transaction) error {
	rows, err := r.db.Prepare("INSERT INTO " +
		transactionsTable + "(blockNumber, `from`, `to`, amount, hash) " +
		" VALUES(?, ?, ?, ?, ?)")
	if err != nil {
		return err
	}
	defer rows.Close()

	res, err := rows.Exec(
		transaction.BlockNumber,
		transaction.To,
		transaction.Amount,
		transaction.Hash)
	if err != nil {
		return err
	}

	lastInsertId, _ := res.LastInsertId()
	transaction.Id = int(lastInsertId)

	return nil
}

func (r *transactionRepository) GetByHashAndAddress(hash string, address string) (*tx.Transaction, error) {
	transaction := tx.Transaction{}

	query := "SELECT * FROM " + transactionsTable + " WHERE `hash` = \"" + hash + "\" AND `to` = \"" + address + "\""

	rows, err := r.db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		err = mapTransaction(rows, &transaction)
		if err != nil {
			return nil, err
		}
	}

	return &transaction, nil
}

func (r *transactionRepository) GetAll(limit int) ([]tx.Transaction, error) {
	query := "SELECT * FROM " + transactionsTable + " ORDER BY id DESC "
	limitQuery := "LIMIT " + strconv.Itoa(limit)
	if limit > 0 {
		query += limitQuery
	}

	rows, err := r.db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	transactions := []tx.Transaction{}

	for rows.Next() {
		transaction := tx.Transaction{}
		err = mapTransaction(rows, &transaction)
		if err != nil {
			return nil, err
		}
		transactions = append(transactions, transaction)
	}

	return transactions, nil
}

func (r *transactionRepository) GetAddresses(limit int) ([]string, error) {
	query := "SELECT `to` FROM " + transactionsTable + " WHERE isToPool = 0 ORDER BY id DESC "
	limitQuery := "LIMIT " + strconv.Itoa(limit)
	if limit > 0 {
		query += limitQuery
	}

	rows, err := r.db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	addresses := []string{}

	for rows.Next() {
		var to string

		err = rows.Scan(&to)
		if err != nil {
			return nil, err
		}

		addresses = append(addresses, to)
	}

	return addresses, nil
}

func (r *transactionRepository) GetUnspentAddresses(limit int, exceptions []string) ([]string, error) {
	var exceptionQuery string

	for i, exceptionAddress := range exceptions {
		exceptionQuery += "\"" + exceptionAddress + "\""
		if i != len(exceptions)-1 {
			exceptionQuery += ", "
		}
	}

	query := "SELECT `to` FROM " + transactionsTable + " WHERE "
	query += " `to` NOT IN (" + exceptionQuery + ") AND isToPool = 0 "
	query += " ORDER BY id desc "
	limitQuery := fmt.Sprintf("LIMIT %d", limit)

	if limit > 0 {
		query += limitQuery
	}

	rows, err := r.db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	addresses := []string{}
	addressesExist := map[string]bool{}

	for rows.Next() {
		var to string

		err = rows.Scan(&to)
		if err != nil {
			return nil, err
		}

		if _, addressExists := addressesExist[to]; addressExists {
			continue
		}

		addresses = append([]string{to}, addresses...)
		addressesExist[to] = true
	}

	return addresses, nil
}

func (r *transactionRepository) UpdateUnspentAddresses(addresses string) error {
	query := "UPDATE " + transactionsTable + " SET " +
		" `isToPool` = 1 " +
		" WHERE `to` IN (" + addresses + ") AND `isToPool` = 0"

	rows, err := r.db.Query(query)
	if err != nil {
		return err
	}
	defer rows.Close()

	return nil
}

func (r *transactionRepository) GetLatestNumConfirmations() (int, error) {
	transaction := tx.Transaction{}

	query := "SELECT NumConfirmation FROM " + transactionsTable + " ORDER BY NumConfirmation LIMIT 1"

	rows, err := r.db.Query(query)
	if err != nil {
		return 0, err
	}
	defer rows.Close()

	for rows.Next() {
		err := rows.Scan(
			&transaction.NumConfirmation,
		)
		if err != nil {
			return 0, err
		}
	}

	return transaction.NumConfirmation, nil
}

func (r *transactionRepository) GetConfTransactions(limit, conf int) ([]tx.Transaction, error) {
	query := "SELECT * FROM " + transactionsTable + " WHERE NumConfirmation > " + strconv.Itoa(limit) + " ORDER BY id DESC "
	limitQuery := "LIMIT " + strconv.Itoa(limit)
	if limit > 0 {
		query += limitQuery
	}

	rows, err := r.db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	transactions := []tx.Transaction{}

	for rows.Next() {
		transaction := tx.Transaction{}
		err = mapTransaction(rows, &transaction)
		if err != nil {
			return nil, err
		}
		transactions = append(transactions, transaction)
	}

	return transactions, nil
}

func (r *transactionRepository) UpdateNumConfirmation(id int, numConfirmations int) error {
	query := "UPDATE " + transactionsTable + " SET " +
		" `NumConfirmation` = " + strconv.Itoa(numConfirmations) +
		" WHERE `Id` " + strconv.Itoa(id)

	rows, err := r.db.Query(query)
	if err != nil {
		return err
	}
	defer rows.Close()

	return nil
}

func mapTransaction(rows *sql.Rows, transaction *tx.Transaction) error {
	err := rows.Scan(
		&transaction.Id,
		&transaction.BlockNumber,
		&transaction.NumConfirmation,
		&transaction.To,
		&transaction.Amount,
		&transaction.Hash,
		&transaction.IsToPool,
	)

	if err != nil {
		return err
	}
	return nil
}
