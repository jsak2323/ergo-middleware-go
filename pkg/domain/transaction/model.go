package transaction

type Transaction struct {
	Id          int
	BlockNumber string
	From        string
	To          string
	Amount      string
	Hash        string
	IsToPool    int
}

type AddressToken struct {
	Id    int
	Token string
	To    string
}
