package cosmicpythongo

import (
	"database/sql"
	"fmt"

	_ "github.com/mattn/go-sqlite3"
)

const dbFile string = "orders.sqlite"

type SQLRepository struct {
	db *sql.DB
}

func NewSQLRepository() SQLRepository {
	return SQLRepository{}
}

const insertBatchRow string = `INSERT INTO batches VALUES(?,?,?,?)`

func (s *SQLRepository) Add(batch Batch) error {
	if _, err := s.db.Exec(insertBatchRow, batch.reference, batch.sku, batch.quantity, batch.ETA); err != nil {
		return fmt.Errorf("could not add persist batch to db %w", err)
	}

	return nil
}
