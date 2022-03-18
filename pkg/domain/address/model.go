package address

type Address struct {
	Id      int    `db:"id"`
	Created int    `db:"created"`
	Address string `db:"address"`
}
