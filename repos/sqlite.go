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

type SqliteProductsRepository struct {
	db dbTx
}

const insertBatchRow string = `INSERT INTO batches VALUES(?,?,?,?)`
const insertOrderLineRow string = `INSERT INTO order_lines VALUES (?,?,?)`
const insertBatchOrderLineRow string = `INSERT INTO batches_order_lines VALUES (?,?)`
const selectBatchRow string = `SELECT reference, sku, quantity, eta FROM "batches" WHERE reference=?`
const selectBatchAllocations string = `SELECT * FROM batches_order_lines WHERE batch_id=?`
const selectOrderLineRow string = `SELECT * FROM order_lines WHERE order_id=?`

type option func() (func(*SqliteProductsRepository), error)

func success(opt func(*SqliteProductsRepository)) option {
	return func() (func(*SqliteProductsRepository), error) {
		return opt, nil
	}
}

func failure(err error) option {
	return func() (func(*SqliteProductsRepository), error) {
		return nil, err
	}
}

func NewSqliteRepository(options ...option) (*SqliteProductsRepository, error) {
	repo := &SqliteProductsRepository{}
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

	return success(func(s *SqliteProductsRepository) {
		s.db = db
	})
}

// Uses the provided transaction
func WithDBTransaction(dbTransaction dbTx) option {
	return success(func(s *SqliteProductsRepository) {
		s.db = dbTransaction
	})
}

func (s *SqliteProductsRepository) Add(product domain.Product) error {
	return nil
}

func (s *SqliteProductsRepository) Get(sku domain.Sku) (domain.Product, error) {
	productBatches, err := s.ListBatches(sku)
	if len(productBatches) == 0 {
		return domain.Product{}, apperrors.NonExistentProductError{Sku: sku}
	}
	if err != nil {
		return domain.Product{}, fmt.Errorf("could not get batches for product")
	}
	return domain.Product{Sku: sku, Batches: productBatches}, nil
}

func (s *SqliteProductsRepository) AddBatch(batch domain.Batch) error {
	if _, err := s.db.Exec(insertBatchRow, batch.Reference, batch.Sku, batch.Quantity, batch.ETA); err != nil {
		return fmt.Errorf("could not add persist batch to db: %w", err)
	}

	return nil
}

func (s *SqliteProductsRepository) AddOrderLine(orderLine domain.OrderLine) error {
	if _, err := s.db.Exec(insertOrderLineRow, orderLine.OrderID, orderLine.Sku, orderLine.Quantity); err != nil {
		return fmt.Errorf("could not add persist batch to db: %w", err)
	}

	return nil
}

func (s *SqliteProductsRepository) GetBatch(reference domain.Reference) (domain.Batch, error) {
	batch := domain.Batch{
		Allocations: mapset.NewSet[domain.OrderLine](),
	}

	row := s.db.QueryRow(selectBatchRow, reference)

	if err := row.Scan(&batch.Reference, &batch.Sku, &batch.Quantity, &batch.ETA); err != nil {
		return batch, fmt.Errorf("could not find the requested batch %w", err)
	}

	return s.enrichAllocations(batch)
}

func (s *SqliteProductsRepository) enrichAllocations(batch domain.Batch) (domain.Batch, error) {
	allocationsRows, err := s.db.Query(selectBatchAllocations, batch.Reference)

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
		if err := s.db.QueryRow(selectOrderLineRow, orderID).Scan(&orderLine.OrderID, &orderLine.Sku, &orderLine.Quantity); err != nil {
			return batch, fmt.Errorf("could not scan the order line with id %q: %w", orderID, err)
		}
		batch.Allocate(orderLine)
	}

	if err := allocationsRows.Err(); err != nil {
		return batch, fmt.Errorf("an error occurred while iterating over allocations: %w", err)
	}

	return batch, nil
}

func (s *SqliteProductsRepository) ListBatches(sku domain.Sku) ([]domain.Batch, error) {
	var batchList []domain.Batch

	batchRows, err := s.db.Query(`SELECT * FROM batches WHERE sku=?`, sku)

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

func (s *SqliteProductsRepository) AllocateToBatch(batch domain.Batch, orderLine domain.OrderLine) error {
	batch, err := s.GetBatch(batch.Reference)
	if err != nil {
		return fmt.Errorf("could not find batch: %s", err)
	}

	if err = batch.Allocate(orderLine); err != nil {
		return fmt.Errorf("cannot allocate this order to this batch: %s", err)
	}

	if _, err := s.db.Exec(insertBatchOrderLineRow, batch.Reference, orderLine.OrderID); err != nil {
		return fmt.Errorf("failed to store allocation to db: %s", err)
	}

	return nil
}

func (s *SqliteProductsRepository) DeallocateFromBatch(batch domain.Batch, orderLine domain.OrderLine) error {
	batch, err := s.GetBatch(batch.Reference)
	if err != nil {
		return fmt.Errorf("could not find batch: %s", err)
	}

	batch.Deallocate(orderLine)

	deleteQuery := fmt.Sprintf("DELETE FROM batches_order_lines WHERE batch_id=%q AND order_id=%q", batch.Reference, orderLine.OrderID)
	if _, err := s.db.Exec(deleteQuery); err != nil {
		return fmt.Errorf("failed to store allocation to db: %s", err)
	}

	return nil
}
