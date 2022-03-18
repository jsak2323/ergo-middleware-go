package mysql

import (
	"database/sql"
	"strconv"

	_ "github.com/go-sql-driver/mysql"

	bl "github.com/btcid/ergo-middleware-go/pkg/domain/blocks"
)

const blocksTable = "blocks"

type blocksRepository struct {
	db *sql.DB
}

func NewMysqlBlocksRepository(db *sql.DB) bl.BlocksRepository {
	return &blocksRepository{
		db,
	}
}

func (r *blocksRepository) Get() (bl.Blocks, error) {
	query := "SELECT * FROM " + blocksTable
	blocks := bl.Blocks{}

	rows, err := r.db.Query(query)
	if err != nil {
		return blocks, err
	}
	defer rows.Close()

	for rows.Next() {
		err = rows.Scan(
			&blocks.LastUpdateTime,
			&blocks.LastUpdatedBlockNum,
		)
		if err != nil {
			return blocks, err
		}
	}

	return blocks, nil
}

func (r *blocksRepository) Update(blocks bl.Blocks) error {
	query := "UPDATE " + blocksTable + " SET " +
		" `lastUpdateTime` = " + strconv.Itoa(blocks.LastUpdateTime) + ", " +
		" `lastUpdatedBlockNum` = \"" + blocks.LastUpdatedBlockNum + "\""

	rows, err := r.db.Query(query)
	if err != nil {
		return err
	}
	defer rows.Close()

	return nil
}
