package mysql

import (
	"database/sql"

	_ "github.com/go-sql-driver/mysql"

	ad "github.com/btcid/ergo-middleware-go/pkg/domain/address"
)

const addressesTable = "addresses"

type addressRepository struct {
	db *sql.DB
}

func NewMysqlAddressRepository(db *sql.DB) ad.AddressRepository {
	return &addressRepository{
		db,
	}
}

func (r *addressRepository) Create(address *ad.Address) error {
	rows, err := r.db.Prepare("INSERT INTO " + addressesTable +
		"(created_at,  address)" +
		" VALUES(?, ?)")
	if err != nil {
		return err
	}
	defer rows.Close()

	_, insErr := rows.Exec(address.Created, address.Address)
	if insErr != nil {
		return err
	}

	return nil
}

func (r *addressRepository) GetAllAddress() ([]string, error) {
	query := "SELECT address FROM " + addressesTable
	addresses := []string{}

	rows, err := r.db.Query(query)
	if err != nil {
		return addresses, err
	}
	defer rows.Close()

	for rows.Next() {
		var address string
		err = rows.Scan(&address)
		if err != nil {
			return addresses, err
		}

		addresses = append(addresses, address)
	}

	return addresses, nil
}

func (r *addressRepository) GetByAddress(address string) (*ad.Address, error) {

	addressObj := ad.Address{}

	query := "SELECT * FROM " + addressesTable + " WHERE `address` = \"" + address + "\""

	rows, err := r.db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		err = mapAddressObj(rows, &addressObj)
		if err != nil {
			return nil, err
		}
	}

	return &addressObj, nil
}

func mapAddressObj(rows *sql.Rows, addressObj *ad.Address) error {
	err := rows.Scan(
		&addressObj.Id,
		&addressObj.Created,
		&addressObj.Address,
	)

	if err != nil {
		return err
	}
	return nil
}
