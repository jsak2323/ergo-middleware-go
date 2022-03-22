package transaction

type TransactionRepository interface {
	Create(transaction *Transaction) error
	GetByHashAndAddress(hash string, address string) (*Transaction, error)
	GetAll(limit int) ([]Transaction, error)
	GetAddresses(limit int) ([]string, error)
	GetUnspentAddresses(limit int, exceptions []string) ([]string, error)
	UpdateUnspentAddresses(addresses string) error
	GetLatestNumConfirmations() (int, error)
	GetConfTransactions(limit, conf int) ([]Transaction, error)
	UpdateNumConfirmation(id int, numConfirmations int) error
}
