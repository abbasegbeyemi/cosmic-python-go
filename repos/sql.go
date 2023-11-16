package repos

import (
	"database/sql"
	"fmt"

	"github.com/abbasegbeyemi/cosmic-python-go/domain"
	mapset "github.com/deckarep/golang-set/v2"
	_ "github.com/mattn/go-sqlite3"
)

type SQLRepository struct {
	DB *sql.DB
}

func NewSQLRepository() SQLRepository {
	return SQLRepository{}
}

const insertBatchRow string = `INSERT INTO batches VALUES(?,?,?,?)`

func (s *SQLRepository) AddBatch(batch domain.Batch) error {
	if _, err := s.DB.Exec(insertBatchRow, batch.Reference, batch.Sku, batch.Quantity, batch.ETA); err != nil {
		return fmt.Errorf("could not add persist batch to db %w", err)
	}

	return nil
}

func (s *SQLRepository) AddOrderLine(orderLine domain.OrderLine) error {
	if _, err := s.DB.Exec(`INSERT INTO order_lines VALUES (?,?,?)`, orderLine.OrderID, orderLine.Sku, orderLine.Quantity); err != nil {
		return fmt.Errorf("could not add persist batch to db %w", err)
	}

	return nil
}

func (s *SQLRepository) GetBatch(reference domain.Reference) (domain.Batch, error) {
	batch := domain.Batch{
		Allocations: mapset.NewSet[domain.OrderLine](),
	}

	row := s.DB.QueryRow(`SELECT * FROM "batches" WHERE reference=?`, reference)

	if err := row.Scan(&batch.Reference, &batch.Sku, &batch.Quantity, &batch.ETA); err != nil {
		return batch, fmt.Errorf("could not find the requested batch %w", err)
	}

	return s.enrichAllocations(batch)
}

func (s *SQLRepository) enrichAllocations(batch domain.Batch) (domain.Batch, error) {
	allocationsRows, err := s.DB.Query(`SELECT * FROM batches_order_lines WHERE batch_id=?`, batch.Reference)

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

		orderLine := domain.OrderLine{}
		if err := s.DB.QueryRow(`SELECT * FROM order_lines WHERE order_id=?`, orderID).Scan(&orderLine.OrderID, &orderLine.Sku, &orderLine.Quantity); err != nil {
			return batch, fmt.Errorf("could not scan the order line with id %q: %w", orderID, err)
		}
		batch.Allocate(orderLine)
	}

	if err := allocationsRows.Err(); err != nil {
		return batch, fmt.Errorf("an error occurred while iterating over allocations: %w", err)
	}

	return batch, nil
}

func (s *SQLRepository) ListBatches() ([]domain.Batch, error) {
	var batchList []domain.Batch

	batchRows, err := s.DB.Query(`SELECT * FROM batches`)

	if err != nil {
		return batchList, fmt.Errorf("could not get batches: %w", err)
	}

	defer batchRows.Close()

	for batchRows.Next() {
		batch := domain.Batch{
			Allocations: mapset.NewSet[domain.OrderLine](),
		}

		if err := batchRows.Scan(&batch.Reference, &batch.Sku, &batch.Quantity, &batch.ETA); err != nil {
			return batchList, fmt.Errorf("could not scan when generating batch list: %s", err)
		}
		batch, err = s.enrichAllocations(batch)
		if err != nil {
			return batchList, fmt.Errorf("could not enrich allocations for batchReference %s: %s", batch.Reference, err)
		}

		batchList = append(batchList, batch)
	}

	return batchList, nil
}

func (s *SQLRepository) AllocateToBatch(batch domain.Batch, orderLine domain.OrderLine) error {
	batch, err := s.GetBatch(batch.Reference)
	if err != nil {
		return fmt.Errorf("could not find batch: %s", err)
	}

	if canAllocate, reason := batch.CanAllocate(orderLine); !canAllocate {
		return fmt.Errorf("cannot allocate this order to this batch: %s", reason)
	}

	if _, err := s.DB.Exec(`INSERT INTO batches_order_lines VALUES (?,?)`, batch.Reference, orderLine.OrderID); err != nil {
		return fmt.Errorf("failed to store allocation to db: %s", err)
	}

	return nil
}
