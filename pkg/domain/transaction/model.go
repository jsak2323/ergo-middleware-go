package transaction

type Transaction struct {
	Id              int
	Hash            string
	BlockNumber     string
	From            *string
	To              string
	Amount          string
	NumConfirmation int
}
