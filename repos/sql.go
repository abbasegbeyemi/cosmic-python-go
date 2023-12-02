package repos

import (
	"database/sql"
	"fmt"

	"github.com/abbasegbeyemi/cosmic-python-go/apperrors"
	"github.com/abbasegbeyemi/cosmic-python-go/domain"
	mapset "github.com/deckarep/golang-set/v2"
	_ "github.com/mattn/go-sqlite3"
)

type dbTx interface {
	Exec(string, ...any) (sql.Result, error)
	Query(string, ...any) (*sql.Rows, error)
	QueryRow(string, ...any) *sql.Row
}

type SQLProductsRepository struct {
	db dbTx
}

const insertBatchRow string = `INSERT INTO batches (reference, sku, quantity, eta) VALUES (?,?,?,?)`
const insertBatchOrderLineAllocation = `INSERT INTO order_lines (order_id, sku, quantity, batch_id) VALUES (?,?,?,?)`
const selectBatchRow string = `SELECT reference, sku, quantity, eta FROM "batches" WHERE reference=?`
const selectBatchAllocations string = `SELECT order_id, sku, quantity FROM order_lines WHERE batch_id=?`

type option func() (func(*SQLProductsRepository), error)

func success(opt func(*SQLProductsRepository)) option {
	return func() (func(*SQLProductsRepository), error) {
		return opt, nil
	}
}

func failure(err error) option {
	return func() (func(*SQLProductsRepository), error) {
		return nil, err
	}
}

func NewSqliteRepository(options ...option) (*SQLProductsRepository, error) {
	repo := &SQLProductsRepository{}
	for _, option := range options {
		opt, err := option()
		if err != nil {
			return nil, err
		}
		opt(repo)
	}
	return repo, nil
}

// Uses the sqlite database file and creates the transaction
func WithDBFile(filepath string) option {
	db, err := sql.Open("sqlite3", filepath)

	if err != nil {
		return failure(fmt.Errorf("could not open sqlite filepath: %w", err))
	}

	return success(func(s *SQLProductsRepository) {
		s.db = db
	})
}

// Uses the provided transaction
func WithDBTransaction(dbTransaction dbTx) option {
	return success(func(s *SQLProductsRepository) {
		s.db = dbTransaction
	})
}

func (s *SQLProductsRepository) Add(product domain.Product) error {
	return nil
}

func (s *SQLProductsRepository) Get(sku domain.Sku) (domain.Product, error) {
	productBatches, err := s.ListBatches(sku)
	if len(productBatches) == 0 {
		return domain.Product{}, apperrors.NonExistentProductError{Sku: sku}
	}
	if err != nil {
		return domain.Product{}, fmt.Errorf("could not get batches for product")
	}
	return domain.Product{Sku: sku, Batches: productBatches}, nil
}

func (s *SQLProductsRepository) AddBatch(batch domain.Batch) error {
	if _, err := s.db.Exec(insertBatchRow, batch.Reference, batch.Sku, batch.Quantity, batch.ETA); err != nil {
		return fmt.Errorf("could not persist batch to db: %w", err)
	}

	return nil
}

func (s *SQLProductsRepository) GetBatch(reference domain.Reference) (domain.Batch, error) {
	batch := domain.Batch{
		Allocations: mapset.NewSet[domain.OrderLine](),
	}

	row := s.db.QueryRow(selectBatchRow, reference)

	if err := row.Scan(&batch.Reference, &batch.Sku, &batch.Quantity, &batch.ETA); err != nil {
		return batch, fmt.Errorf("could not find the requested batch %w", err)
	}

	return enrichAllocations(s.db, batch)
}

func getBatchIdFromReference(db dbTx, batchRef domain.Reference) (int, error) {
	var batchId int
	if err := db.QueryRow("SELECT id FROM batches WHERE reference=?", batchRef).Scan(&batchId); err != nil {
		return 0, fmt.Errorf("could not get batch id: %w", err)
	}
	return batchId, nil
}

func enrichAllocations(db dbTx, batch domain.Batch) (domain.Batch, error) {
	batchId, err := getBatchIdFromReference(db, batch.Reference)
	if err != nil {
		return domain.Batch{}, err
	}

	allocationsRows, err := db.Query(selectBatchAllocations, batchId)
	if err != nil {
		return batch, fmt.Errorf("could not get allocations for batch: %w", err)
	}

	defer allocationsRows.Close()

	for allocationsRows.Next() {
		orderLine := domain.OrderLine{}
		if err := allocationsRows.Scan(&orderLine.OrderID, &orderLine.Sku, &orderLine.Quantity); err != nil {
			return batch, fmt.Errorf("could not scan the order line: %w", err)
		}

		batch.Allocate(orderLine)
	}

	if err := allocationsRows.Err(); err != nil {
		return batch, fmt.Errorf("an error occurred while iterating over allocations: %w", err)
	}

	return batch, nil
}

func (s *SQLProductsRepository) ListBatches(sku domain.Sku) ([]domain.Batch, error) {
	var batchList []domain.Batch

	batchRows, err := s.db.Query(`SELECT reference, sku, quantity, eta FROM batches WHERE sku=?`, sku)

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
		batch, err = enrichAllocations(s.db, batch)
		if err != nil {
			return batchList, fmt.Errorf("could not enrich allocations for batch reference %s: %s", batch.Reference, err)
		}

		batchList = append(batchList, batch)
	}

	return batchList, nil
}

func (s *SQLProductsRepository) AllocateToBatch(batch domain.Batch, orderLine domain.OrderLine) error {
	batchEnriched, err := s.GetBatch(batch.Reference)
	if err != nil {
		return fmt.Errorf("could not enrich batch: %w", err)
	}

	if err = batchEnriched.Allocate(orderLine); err != nil {
		return fmt.Errorf("cannot allocate this order to this batch: %s", err)
	}

	batchId, err := getBatchIdFromReference(s.db, batchEnriched.Reference)
	if err != nil {
		return fmt.Errorf("could not find batch")
	}

	if _, err := s.db.Exec(insertBatchOrderLineAllocation, orderLine.OrderID, orderLine.Sku, orderLine.Quantity, batchId); err != nil {
		return fmt.Errorf("failed to store allocation to db: %s", err)
	}

	return nil
}

func (s *SQLProductsRepository) DeallocateFromBatch(batch domain.Batch, orderLine domain.OrderLine) error {
	batchEnriched, err := s.GetBatch(batch.Reference)
	if err != nil {
		return fmt.Errorf("could not enrich batch: %w", err)
	}

	batchEnriched.Deallocate(orderLine)

	batchId, err := getBatchIdFromReference(s.db, batchEnriched.Reference)
	if err != nil {
		return fmt.Errorf("could not find batch: %s", err)
	}

	if _, err = s.db.Exec("UPDATE order_lines SET batch_id=NULL WHERE order_id=? AND sku=? AND quantity=? AND batch_id=?", orderLine.OrderID, orderLine.Sku, orderLine.Quantity, batchId); err != nil {
		return fmt.Errorf("failed to deallocate order line from batch: %w", err)
	}

	return nil
}
