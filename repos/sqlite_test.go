package repos

import (
	"slices"
	"testing"
	"time"

	"github.com/abbasegbeyemi/cosmic-python-go/domain"

	"github.com/abbasegbeyemi/cosmic-python-go/test"
	"github.com/stretchr/testify/assert"
)

func TestSqliteRepository_AddBatch(t *testing.T) {
	db := test.SqliteDB(t)
	test.CreateTables(t, db)
	defer test.TruncateTables(t, db)
	t.Run("can store batch", func(t *testing.T) {
		tx, err := db.Begin()
		assert.Nil(t, err)
		batch := domain.NewBatch(
			"batch-001",
			"SMALL-TABLE",
			10,
			time.Time{},
		)

		repo := SqliteRepository{
			db: db,
		}
		err = repo.AddBatch(batch)
		assert.Nil(t, err)

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

func TestSqliteRepository_GetBatch(t *testing.T) {
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

		repo, err := NewSqliteRepository(WithDBTransaction(db))
		assert.Nil(t, err)

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
		test.InsertOrderLine(t, db, orderId, sku, 3)
		test.InsertAllocation(t, db, batchRef, orderId)

		repo, err := NewSqliteRepository(WithDBTransaction(db))
		assert.Nil(t, err)

		receivedBatch, err := repo.GetBatch(batchRef)

		assert.Nil(t, err)
		allocations := receivedBatch.Allocations.ToSlice()
		assert.Greater(t, len(allocations), 0)
		assert.Equal(t, allocations[0].OrderID, orderId)
		assert.Equal(t, allocations[0].Sku, sku)
	})
}

func TestSqliteRepository_ListBatch(t *testing.T) {
	db := test.SqliteDB(t)
	test.CreateTables(t, db)
	defer test.TruncateTables(t, db)

	sku := domain.Sku("BIG-STICKER")

	test.InsertBatch(t, db, "batch-007", sku, 50, time.Time{})
	test.InsertBatch(t, db, "batch-001", sku, 20, time.Time{})
	test.InsertBatch(t, db, "batch-002", sku, 10, time.Time{})

	repo, err := NewSqliteRepository(WithDBTransaction(db))
	assert.Nil(t, err)

	receivedBatches, err := repo.ListBatches()
	assert.Nil(t, err)

	assert.Len(t, receivedBatches, 3)

	assert.True(t, slices.ContainsFunc(receivedBatches, func(b domain.Batch) bool {
		return b.Reference == domain.Reference("batch-001")
	}))

	assert.True(t, slices.ContainsFunc(receivedBatches, func(b domain.Batch) bool {
		return b.Reference == domain.Reference("batch-002")
	}))

	assert.True(t, slices.ContainsFunc(receivedBatches, func(b domain.Batch) bool {
		return b.Reference == domain.Reference("batch-007")
	}))

}

func TestSqliteRepository_AddOrderLine(t *testing.T) {
	db := test.SqliteDB(t)
	test.CreateTables(t, db)
	defer test.TruncateTables(t, db)
	repo, err := NewSqliteRepository(WithDBTransaction(db))
	assert.Nil(t, err)

	orderLine := domain.OrderLine{
		OrderID:  "order-001",
		Sku:      "LARGE-TABLE",
		Quantity: 12,
	}

	err = repo.AddOrderLine(orderLine)
	assert.Nil(t, err)

}

func TestSqliteRepository_AllocateToBatch(t *testing.T) {
	db := test.SqliteDB(t)
	test.CreateTables(t, db)
	defer test.TruncateTables(t, db)

	repo, err := NewSqliteRepository(WithDBTransaction(db))
	assert.Nil(t, err)

	sku := "LARGE-TABLE"

	batch := domain.Batch{
		Reference: "batch-012",
		Sku:       domain.Sku(sku),
		Quantity:  40,
	}

	orderLine := domain.OrderLine{
		OrderID:  "order-321",
		Quantity: 12,
		Sku:      domain.Sku(sku),
	}

	err = repo.AddBatch(batch)
	assert.Nil(t, err)

	err = repo.AddOrderLine(orderLine)
	assert.Nil(t, err)

	err = repo.AllocateToBatch(batch, orderLine)
	assert.Nil(t, err)

	allocatedBatch, err := repo.GetBatch(batch.Reference)
	assert.Nil(t, err)

	assert.EqualExportedValues(t, allocatedBatch.Allocations.ToSlice()[0], orderLine)
}

func TestSqliteRepository_DeallocateFromBatch(t *testing.T) {
	db := test.SqliteDB(t)
	test.CreateTables(t, db)
	defer test.TruncateTables(t, db)

	repo, err := NewSqliteRepository(WithDBTransaction(db))
	assert.Nil(t, err)

	sku := "LARGE-TABLE"

	batch := domain.Batch{
		Reference: "batch-012",
		Sku:       domain.Sku(sku),
		Quantity:  40,
	}

	orderLine := domain.OrderLine{
		OrderID:  "order-321",
		Quantity: 12,
		Sku:      domain.Sku(sku),
	}

	err = repo.AddBatch(batch)
	assert.Nil(t, err)

	err = repo.AddOrderLine(orderLine)
	assert.Nil(t, err)

	err = repo.AllocateToBatch(batch, orderLine)
	assert.Nil(t, err)

	allocatedBatch, err := repo.GetBatch(batch.Reference)
	assert.Nil(t, err)

	assert.EqualExportedValues(t, allocatedBatch.Allocations.ToSlice()[0], orderLine)

	err = repo.DeallocateFromBatch(allocatedBatch, orderLine)
	assert.Nil(t, err)

	deallocatedBatch, err := repo.GetBatch(batch.Reference)
	assert.Nil(t, err)

	assert.Len(t, deallocatedBatch.Allocations.ToSlice(), 0)
}
