package cosmicpythongo

import (
	"database/sql"
	"fmt"

	mapset "github.com/deckarep/golang-set/v2"
	_ "github.com/mattn/go-sqlite3"
)

type SQLRepository struct {
	db *sql.DB
}

func NewSQLRepository() SQLRepository {
	return SQLRepository{}
}

const insertBatchRow string = `INSERT INTO batches VALUES(?,?,?,?)`

func (s *SQLRepository) AddBatch(batch Batch) error {
	if _, err := s.db.Exec(insertBatchRow, batch.reference, batch.sku, batch.quantity, batch.ETA); err != nil {
		return fmt.Errorf("could not add persist batch to db %w", err)
	}

	return nil
}

func (s *SQLRepository) GetBatch(reference Reference) (Batch, error) {
	batch := Batch{
		allocations: mapset.NewSet[OrderLine](),
	}

	row := s.db.QueryRow(`SELECT * FROM "batches" WHERE reference=?`, reference)
	if err := row.Scan(&batch.reference, &batch.sku, &batch.quantity, &batch.ETA); err != nil {
		return batch, fmt.Errorf("could not find the requested batch %w", err)
	}

	// Get the allocations
	allocationsRows, err := s.db.Query(`SELECT * FROM batches_order_lines WHERE batch_id=?`, reference)

	if err != nil {
		return batch, fmt.Errorf("could not get allocations for batch: %w", err)
	}
	defer allocationsRows.Close()

	for allocationsRows.Next() {
		var orderID string
		var batchID string
		if err := allocationsRows.Scan(&batchID, &orderID); err != nil {
			return batch, fmt.Errorf("could not scan the order id: %w", err)
		}

		orderLine := OrderLine{}
		if err := s.db.QueryRow(`SELECT * FROM order_lines WHERE order_id=?`, orderID).Scan(&orderLine.OrderID, &orderLine.Sku, &orderLine.Quantity); err != nil {
			return batch, fmt.Errorf("could not scan the order line with id %q: %w", orderID, err)
		}
		batch.Allocate(orderLine)
	}

	if err := allocationsRows.Err(); err != nil {
		return batch, fmt.Errorf("an error occured while iterating over allocations: %w", err)
	}

	return batch, nil
}
