package transaction

type Transaction struct {
	Id              int
	BlockNumber     string
	NumConfirmation int
	To              string
	Amount          string
	Hash            string
	IsToPool        int
}
