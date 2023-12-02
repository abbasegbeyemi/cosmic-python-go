package repos

import (
	"testing"
	"time"

	"github.com/abbasegbeyemi/cosmic-python-go/domain"

	"github.com/abbasegbeyemi/cosmic-python-go/test"
	"github.com/stretchr/testify/assert"
)

func TestSqlProductsRepository_AddBatch(t *testing.T) {
	db := test.SqliteDB(t)
	test.CreateTables(t, db)
	// defer test.TruncateTables(t, db)
	t.Run("can store batch", func(t *testing.T) {
		tx, err := db.Begin()
		assert.Nil(t, err)
		batch := domain.NewBatch(
			"batch-001",
			"SMALL-TABLE",
			10,
			time.Time{},
		)

		repo := SQLProductsRepository{
			db: db,
		}
		err = repo.AddBatch(batch)
		assert.NoError(t, err)

		err = tx.Commit()
		assert.Nil(t, err)

		createdBatch := domain.Batch{}
		row := db.QueryRow(selectBatchRow, batch.Reference)
		err = row.Scan(&createdBatch.Reference, &createdBatch.Sku, &createdBatch.Quantity, &createdBatch.ETA)
		assert.Nil(t, err)

		assert.Equal(t, batch.Reference, createdBatch.Reference)
		assert.Equal(t, batch.Sku, createdBatch.Sku)
		assert.Equal(t, batch.Quantity, createdBatch.Quantity)
		assert.Equal(t, batch.ETA, createdBatch.ETA)
	})
}

func TestSqlProductsRepository_GetBatch(t *testing.T) {
	db := test.SqliteDB(t)
	test.CreateTables(t, db)
	defer test.TruncateTables(t, db)
	t.Run("can retrieve batch", func(t *testing.T) {

		existingBatch := domain.Batch{
			Reference: "batch-002",
			Sku:       "LARGE-MIRROR",
			Quantity:  23,
			ETA:       time.Now().AddDate(0, 3, 0).UTC(),
		}
		db.Exec(insertBatchRow, existingBatch.Reference, existingBatch.Sku, existingBatch.Quantity, existingBatch.ETA)

		repo := SQLProductsRepository{
			db: db,
		}

		receivedBatch, err := repo.GetBatch(existingBatch.Reference)

		assert.Nil(t, err)
		assert.Equal(t, existingBatch.Reference, receivedBatch.Reference)
		assert.Equal(t, existingBatch.Sku, receivedBatch.Sku)
		assert.Equal(t, existingBatch.Quantity, receivedBatch.Quantity)
		assert.Equal(t, existingBatch.ETA, receivedBatch.ETA)
	})

	t.Run("can retrieve batch with allocation", func(t *testing.T) {
		db := test.SqliteDB(t)
		test.CreateTables(t, db)
		defer test.TruncateTables(t, db)

		batchRef := domain.Reference("batch-001")
		orderId := domain.Reference("order-012")
		sku := domain.Sku("LARGE-MIRROR")

		test.InsertBatch(t, db, batchRef, sku, 50, time.Time{})
		test.InsertAllocation(t, db, batchRef, orderId, sku, 3)

		repo := SQLProductsRepository{
			db: db,
		}

		receivedBatch, err := repo.GetBatch(batchRef)

		assert.Nil(t, err)
		allocations := receivedBatch.Allocations.ToSlice()
		if assert.Greater(t, len(allocations), 0) {
			assert.Equal(t, allocations[0].OrderID, orderId)
			assert.Equal(t, allocations[0].Sku, sku)
		}
	})
}

func TestSqlProductsRepository_AllocateToBatch(t *testing.T) {
	db := test.SqliteDB(t)
	test.CreateTables(t, db)
	defer test.TruncateTables(t, db)

	repo := SQLProductsRepository{
		db: db,
	}

	sku := domain.Sku("LARGE-TABLE")

	batch := domain.NewBatch("batch-012", sku, 40, time.Time{})

	err := repo.AddBatch(batch)
	assert.NoError(t, err)

	orderLine := domain.OrderLine{
		OrderID:  "order-321",
		Quantity: 12,
		Sku:      domain.Sku(sku),
	}
	err = repo.AllocateToBatch(batch, orderLine)
	assert.NoError(t, err)

	allocatedBatch, err := repo.GetBatch(batch.Reference)
	assert.NoError(t, err)

	assert.EqualExportedValues(t, allocatedBatch.Allocations.ToSlice()[0], orderLine)
}

func TestSqlProductsRepository_DeallocateFromBatch(t *testing.T) {
	db := test.SqliteDB(t)
	test.CreateTables(t, db)
	defer test.TruncateTables(t, db)

	repo := SQLProductsRepository{
		db: db,
	}

	sku := domain.Sku("LARGE-TABLE")

	batch := domain.NewBatch("batch-012", sku, 40, time.Time{})

	orderLine := domain.OrderLine{
		OrderID:  "order-321",
		Quantity: 12,
		Sku:      domain.Sku(sku),
	}

	err := repo.AddBatch(batch)
	assert.NoError(t, err)

	err = repo.AllocateToBatch(batch, orderLine)
	assert.NoError(t, err)

	allocatedBatch, err := repo.GetBatch(batch.Reference)
	assert.NoError(t, err)

	assert.EqualExportedValues(t, allocatedBatch.Allocations.ToSlice()[0], orderLine)

	err = repo.DeallocateFromBatch(allocatedBatch, orderLine)
	assert.NoError(t, err)

	deallocatedBatch, err := repo.GetBatch(batch.Reference)
	assert.NoError(t, err)

	assert.Len(t, deallocatedBatch.Allocations.ToSlice(), 0)
}
