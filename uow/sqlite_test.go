package uow

import (
	"database/sql"
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

		err = uow.DBInstruction(func() error {
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
