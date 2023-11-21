package uow

import (
	"database/sql"
	"fmt"
	"testing"
	"time"

	"github.com/abbasegbeyemi/cosmic-python-go/domain"
	"github.com/abbasegbeyemi/cosmic-python-go/test"
	"github.com/stretchr/testify/assert"
)

// Returns the batch ref to which an order line has been allocated
func GetAllocatedBatchRef(t *testing.T, db *sql.DB, orderId domain.Reference, sku domain.Sku) domain.Reference {
	t.Helper()
	var orderLineId string
	err := db.QueryRow("SELECT order_id from order_lines WHERE order_id=? AND sku=?", orderId, sku).Scan(&orderLineId)
	assert.Nil(t, err)

	var batchRef string
	err = db.QueryRow(`SELECT b.reference from batches_order_lines JOIN batches AS b ON batch_id=b.reference WHERE order_id=?`, orderId).Scan(&batchRef)
	assert.Nil(t, err)

	return domain.Reference(batchRef)
}

func TestUOW_Allocate(t *testing.T) {
	db := test.SqliteDB(t)
	test.CreateTables(t, db)
	defer test.TruncateTables(t, db)

	t.Run("can retrieve a batch and allocate to it", func(t *testing.T) {
		sku := domain.Sku("HIPSTER-WORKBENCH")
		test.InsertBatch(t, db, "batch-001", sku, 100, time.Time{})

		uow, err := NewSqliteUnitOfWork(db)
		assert.NoError(t, err)

		err = uow.CommitOnSuccess(func() error {
			batch, err := uow.Batches().GetBatch("batch-001")
			if err != nil {
				return err
			}

			line := domain.OrderLine{OrderID: "order-1", Sku: sku, Quantity: 10}
			if err = uow.Batches().AddOrderLine(line); err != nil {
				return err
			}

			if err = uow.Batches().AllocateToBatch(batch, line); err != nil {
				return err
			}

			return nil
		})

		assert.Nil(t, err)

		batchRef := GetAllocatedBatchRef(t, db, "order-1", sku)
		assert.Equal(t, domain.Reference("batch-001"), batchRef)
	})
}

func TestUOW_Commit(t *testing.T) {

	t.Run("rolls back uncommitted work by default", func(t *testing.T) {
		db := test.SqliteDB(t)
		test.CreateTables(t, db)
		uow, err := NewSqliteUnitOfWork(db)
		assert.NoError(t, err)

		tx, err := uow.Transaction()
		assert.NoError(t, err)

		test.InsertBatch(t, tx, "batch-001", "LARGE-CHAIR", 100, time.Time{})

		var rowCount int
		assert.NoError(t, db.QueryRow("SELECT COUNT(*) FROM batches").Scan(&rowCount))
		assert.Equal(t, 0, rowCount)
	})

	t.Run("rolls back on error", func(t *testing.T) {
		db := test.SqliteDB(t)
		test.CreateTables(t, db)
		uow, err := NewSqliteUnitOfWork(db)
		assert.NoError(t, err)

		uow.CommitOnSuccess(func() error {
			test.InsertBatch(t, uow.transaction, "batch-002", "LARGE-FORK", 100, time.Time{})
			return fmt.Errorf("an error occurred!")
		})
		var rowCount int
		assert.NoError(t, db.QueryRow("SELECT COUNT(*) FROM batches").Scan(&rowCount))
		assert.Equal(t, 0, rowCount)

	})

}