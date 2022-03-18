package blocks

type BlocksRepository interface {
	Get() (Blocks, error)
	Update(blocks Blocks) error
}
