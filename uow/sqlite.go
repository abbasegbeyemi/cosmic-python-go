package uow

import (
	"database/sql"
	"fmt"

	"github.com/abbasegbeyemi/cosmic-python-go/domain"
	"github.com/abbasegbeyemi/cosmic-python-go/repos"
)

type ProductRepository interface {
	Add(domain.Product) error
	Get(domain.Sku) (domain.Product, error)
	AddBatch(domain.Batch) error
	ListBatches(domain.Sku) ([]domain.Batch, error)
	GetBatch(reference domain.Reference) (domain.Batch, error)
	AllocateToBatch(domain.Batch, domain.OrderLine) error
	DeallocateFromBatch(domain.Batch, domain.OrderLine) error
	AddOrderLine(domain.OrderLine) error
}

type SqliteUnitOfWork struct {
	db          *sql.DB
	products    ProductRepository
	transaction *sql.Tx
}

func NewSqliteUnitOfWork(db *sql.DB) (*SqliteUnitOfWork, error) {
	productsRepo, err := repos.NewSqliteRepository(repos.WithDBTransaction(db))
	if err != nil {
		return &SqliteUnitOfWork{}, fmt.Errorf("could not instantiate unit of work: %w", err)
	}

	return &SqliteUnitOfWork{products: productsRepo, db: db}, nil
}

type QueryFunc func() error

func (s *SqliteUnitOfWork) Products() ProductRepository {
	return s.products
}

// Will create a transaction and commit only if the provided query function raises no error.
func (s *SqliteUnitOfWork) CommitOnSuccess(queryFunction QueryFunc) error {
	tx, err := s.Transaction()
	if err != nil {
		return fmt.Errorf("could not initialise db transaction: %w", err)
	}

	repo, err := repos.NewSqliteRepository(repos.WithDBTransaction(tx))
	s.products = repo

	// Reset batches repo back to standard db query
	defer func() {
		s.products, _ = repos.NewSqliteRepository(repos.WithDBTransaction(s.db))
	}()

	if err != nil {
		return fmt.Errorf("could not get sqlite repository: %w", err)
	}

	if err = queryFunction(); err != nil {
		s.transaction.Rollback()
		return fmt.Errorf("query function returned an error %w", err)
	}
	if err := s.transaction.Commit(); err != nil {
		s.transaction.Rollback()
		return fmt.Errorf("could not commit transaction: %w", err)
	}

	return nil
}

// Get a transaction from the unit of work
func (s *SqliteUnitOfWork) Transaction() (*sql.Tx, error) {
	tx, err := s.db.Begin()

	if err != nil {
		return &sql.Tx{}, fmt.Errorf("unable to begin db transaction")
	}
	s.transaction = tx
	return tx, nil
}

func (s *SqliteUnitOfWork) Commit() error {
	return s.transaction.Commit()
}

func (s *SqliteUnitOfWork) Rollback() {
	s.transaction.Rollback()
}
