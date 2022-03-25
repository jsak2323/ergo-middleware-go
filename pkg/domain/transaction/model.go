package transaction

type Transaction struct {
	Id              int
	BlockNumber     string
	NumConfirmation int
	From            string
	To              string
	Amount          string
	Hash            string
	IsToPool        int
}
