package repos

import (
	"database/sql"
	"fmt"

	"github.com/abbasegbeyemi/cosmic-python-go/domain"
	mapset "github.com/deckarep/golang-set/v2"
	_ "github.com/mattn/go-sqlite3"
)

type DB interface {
	Exec(string, ...any) (sql.Result, error)
	Query(string, ...any) (DBRows, error)
	QueryRow(string, ...any) DBRow
}

type DBWrapper struct {
	DB *sql.DB
}

func (db *DBWrapper) Exec(query string, args ...any) (sql.Result, error) {
	return db.DB.Exec(query, args...)
}

func (db *DBWrapper) Query(query string, args ...any) (DBRows, error) {
	return db.DB.Query(query, args...)
}

func (db *DBWrapper) QueryRow(query string, args ...any) DBRow {
	return db.DB.QueryRow(query, args...)
}

type DBRow interface {
	Scan(...any) error
}

type DBRows interface {
	Next() bool
	Close() error
	Scan(...any) error
	Err() error
}

type SQLRepository struct {
	DB DB
}

const insertBatchRow string = `INSERT INTO batches VALUES(?,?,?,?)`
const insertOrderLineRow string = `INSERT INTO order_lines VALUES (?,?,?)`
const insertBatchOrderLineRow string = `INSERT INTO batches_order_lines VALUES (?,?)`
const selectBatchRow string = `SELECT reference, sku, quantity, eta FROM "batches" WHERE reference=?`
const selectAllBatches string = `SELECT * FROM batches`
const selectBatchAllocations string = `SELECT * FROM batches_order_lines WHERE batch_id=?`
const selectOrderLineRow string = `SELECT * FROM order_lines WHERE order_id=?`

func NewSqliteRepository(filepath string) (SQLRepository, error) {
	db, err := sql.Open("sqlite3", filepath)
	if err != nil {
		return SQLRepository{}, fmt.Errorf("could not open sqlite filepath: %w", err)
	}
	return SQLRepository{
		DB: &DBWrapper{DB: db},
	}, nil
}

func (s *SQLRepository) AddBatch(batch domain.Batch) error {
	if _, err := s.DB.Exec(insertBatchRow, batch.Reference, batch.Sku, batch.Quantity, batch.ETA); err != nil {
		return fmt.Errorf("could not add persist batch to db %w", err)
	}

	return nil
}

func (s *SQLRepository) AddOrderLine(orderLine domain.OrderLine) error {
	if _, err := s.DB.Exec(insertOrderLineRow, orderLine.OrderID, orderLine.Sku, orderLine.Quantity); err != nil {
		return fmt.Errorf("could not add persist batch to db %w", err)
	}

	return nil
}

func (s *SQLRepository) GetBatch(reference domain.Reference) (domain.Batch, error) {
	batch := domain.Batch{
		Allocations: mapset.NewSet[domain.OrderLine](),
	}

	row := s.DB.QueryRow(selectBatchRow, reference)

	if err := row.Scan(&batch.Reference, &batch.Sku, &batch.Quantity, &batch.ETA); err != nil {
		return batch, fmt.Errorf("could not find the requested batch %w", err)
	}

	return s.enrichAllocations(batch)
}

func (s *SQLRepository) enrichAllocations(batch domain.Batch) (domain.Batch, error) {
	allocationsRows, err := s.DB.Query(selectBatchAllocations, batch.Reference)

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
		if err := s.DB.QueryRow(selectOrderLineRow, orderID).Scan(&orderLine.OrderID, &orderLine.Sku, &orderLine.Quantity); err != nil {
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

	batchRows, err := s.DB.Query(selectAllBatches)

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

	if _, err := s.DB.Exec(insertBatchOrderLineRow, batch.Reference, orderLine.OrderID); err != nil {
		return fmt.Errorf("failed to store allocation to db: %s", err)
	}

	return nil
}
