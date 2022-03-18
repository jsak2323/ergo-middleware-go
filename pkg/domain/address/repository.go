package address

type AddressRepository interface {
	Create(address *Address) error
	GetAllAddress() ([]string, error)
	GetByAddress(address string) (*Address, error)
}
