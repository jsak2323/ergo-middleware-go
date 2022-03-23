package transaction

type TransactionRepository interface {
	Create(transaction *Transaction) error
	GetByHashAndAddress(hash string, address string) (*Transaction, error)
	GetAll(limit int) ([]Transaction, error)
	GetConfTransactions(limit, conf int) ([]Transaction, error)
	UpdateNumConfirmation(id int, numConfirmations int) error
}
