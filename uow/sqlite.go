package uow

import (
	"database/sql"
	"fmt"

	"github.com/abbasegbeyemi/cosmic-python-go/domain"
	"github.com/abbasegbeyemi/cosmic-python-go/repos"
)

type Repository interface {
	AddBatch(domain.Batch) error
	ListBatches() ([]domain.Batch, error)
	GetBatch(reference domain.Reference) (domain.Batch, error)
	AllocateToBatch(domain.Batch, domain.OrderLine) error
	DeallocateFromBatch(domain.Batch, domain.OrderLine) error
	AddOrderLine(domain.OrderLine) error
}

type SqliteUnitOfWork struct {
	DB          *sql.DB
	batches     Repository
	transaction *sql.Tx
}

func NewSqliteUnitOfWork(db *sql.DB) (*SqliteUnitOfWork, error) {
	batches, err := repos.NewSqliteRepository(repos.WithDBTransaction(db))
	if err != nil {
		return &SqliteUnitOfWork{}, fmt.Errorf("could not instantiate unit of work: %w", err)
	}

	return &SqliteUnitOfWork{batches: batches, DB: db}, nil
}

type QueryFunc func() error

func (s *SqliteUnitOfWork) Batches() Repository {
	return s.batches
}

// Will pass a db transaction to the provided function and perform the db queries within.
// Provided function must return error status. Commits if no error and rolls back if error.
func (s *SqliteUnitOfWork) DBInstruction(queryFunction QueryFunc) error {
	tx, err := s.DB.Begin()

	if err != nil {
		return fmt.Errorf("unable to begin db transaction")
	}
	s.batches, err = repos.NewSqliteRepository(repos.WithDBTransaction(tx))

	// Reset batches repo back to standard db query
	defer func() {
		s.batches, _ = repos.NewSqliteRepository(repos.WithDBTransaction(s.DB))
	}()

	s.transaction = tx

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

func (s *SqliteUnitOfWork) Commit() error {
	return s.transaction.Commit()
}

func (s *SqliteUnitOfWork) Rollback() {
	s.transaction.Rollback()
}
