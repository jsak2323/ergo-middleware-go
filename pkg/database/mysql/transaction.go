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
		transactionsTable + "(blockNumber, `to`, amount, hash, numConfirmation) " +
		" VALUES(?, ?, ?, ?, ?)")
	if err != nil {
		return err
	}
	defer rows.Close()

	res, err := rows.Exec(
		transaction.BlockNumber,
		transaction.To,
		transaction.Amount,
		transaction.Hash,
		transaction.NumConfirmation,
	)
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

func (r *transactionRepository) GetConfTransactions(limit, conf int) ([]tx.Transaction, error) {
	query := "SELECT * FROM " + transactionsTable + " WHERE numConfirmation < " + strconv.Itoa(conf) + " ORDER BY id DESC "
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
		" `numConfirmation` = " + strconv.Itoa(numConfirmations) +
		" WHERE `Id` = " + strconv.Itoa(id)

	fmt.Println(query)
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
		&transaction.Hash,
		&transaction.BlockNumber,
		&transaction.To,
		&transaction.Amount,
		&transaction.NumConfirmation,
		// &transaction.IsToPool,
	)
	if err != nil {
		return err
	}
	return nil
}
