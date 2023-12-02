package repos

import (
	"context"
	"slices"
	"testing"
	"time"

	"github.com/abbasegbeyemi/cosmic-python-go/domain"
	"github.com/abbasegbeyemi/cosmic-python-go/test"
	"github.com/stretchr/testify/assert"
)

func TestORMProductsRepository_AddBatch(t *testing.T) {
	client := test.EntClient(t)
	defer client.Close()
	t.Run("can store batch", func(t *testing.T) {
		batch := domain.NewBatch(
			test.RandomBatchRef(t, ""),
			test.RandomSku(t, ""),
			10,
			time.Time{},
		)
		tx, err := client.Tx(context.Background())
		assert.NoError(t, err)

		repo := ORMProductsRepository{
			client: tx.Client(),
		}

		err = repo.AddBatch(batch)
		assert.NoError(t, err)

		err = tx.Commit()
		assert.NoError(t, err)

		createdBatch := test.EntGetBatch(t, client, batch.Reference)
		assert.NoError(t, err)

		assert.Equal(t, batch, createdBatch)
	})
}

func TestORMProductsRepository_GetBatch(t *testing.T) {
	client := test.EntClient(t)
	defer client.Close()

	t.Run("can retrieve batch", func(t *testing.T) {
		existingBatch := domain.Batch{
			Reference: test.RandomBatchRef(t, ""),
			Sku:       test.RandomSku(t, ""),
			Quantity:  23,
			ETA:       time.Now().AddDate(0, 3, 0).UTC(),
		}
		test.EntInsertBatch(t, client, existingBatch.Reference, existingBatch.Sku, existingBatch.Quantity, existingBatch.ETA)

		repo := ORMProductsRepository{
			client: client,
		}

		receivedBatch, err := repo.GetBatch(existingBatch.Reference)
		assert.NoError(t, err)

		assert.Equal(t, existingBatch.Reference, receivedBatch.Reference)
		assert.Equal(t, existingBatch.Sku, receivedBatch.Sku)
		assert.Equal(t, existingBatch.Quantity, receivedBatch.Quantity)
		assert.Equal(t, existingBatch.ETA, receivedBatch.ETA)
	})

	t.Run("can retrieve batch with allocation", func(t *testing.T) {
		batchRef := test.RandomBatchRef(t, "")
		orderId := test.RandomBatchRef(t, "")
		sku := test.RandomSku(t, "")
		test.EntInsertBatch(t, client, batchRef, sku, 50, time.Time{})
		test.EntInsertAllocation(t, client, batchRef, orderId, sku, 3)

		repo := ORMProductsRepository{
			client: client,
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

func TestORMProductsRepository_ListBatch(t *testing.T) {
	client := test.EntClient(t)
	defer client.Close()

	sku := test.RandomSku(t, "")

	test.EntInsertBatch(t, client, "batch-007", sku, 50, time.Time{})
	test.EntInsertBatch(t, client, "batch-001", sku, 20, time.Time{})
	test.EntInsertBatch(t, client, "batch-002", sku, 10, time.Time{})

	repo := ORMProductsRepository{
		client: client,
	}

	receivedBatches, err := repo.ListBatches(sku)
	assert.NoError(t, err)

	if assert.Len(t, receivedBatches, 3) {
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

}

func TestORMProductsRepository_AllocateToBatch(t *testing.T) {
	client := test.EntClient(t)
	defer client.Close()

	t.Run("can allocate when stock available", func(t *testing.T) {
		repo := ORMProductsRepository{
			client: client,
		}

		sku := test.RandomSku(t, "")

		batch := domain.NewBatch(test.RandomBatchRef(t, ""), sku, 40, time.Time{})

		err := repo.AddBatch(batch)
		assert.NoError(t, err)

		orderLine := domain.OrderLine{
			OrderID:  test.RandomOrderId(t, ""),
			Quantity: 12,
			Sku:      domain.Sku(sku),
		}
		err = repo.AllocateToBatch(batch, orderLine)
		assert.NoError(t, err)

		allocatedBatch, err := repo.GetBatch(batch.Reference)
		assert.NoError(t, err)

		assert.EqualExportedValues(t, allocatedBatch.Allocations.ToSlice()[0], orderLine)
	})

	t.Run("fails to allocate when stock is unavailable", func(t *testing.T) {
		repo := ORMProductsRepository{
			client: client,
		}

		sku := test.RandomSku(t, "")

		batch := domain.NewBatch("batch-012", sku, 2, time.Time{})

		err := repo.AddBatch(batch)
		assert.NoError(t, err)

		orderLine := domain.OrderLine{
			OrderID:  test.RandomOrderId(t, ""),
			Quantity: 12,
			Sku:      sku,
		}
		assert.Error(t, repo.AllocateToBatch(batch, orderLine))
	})

	t.Run("ensures allocation always uses information from db", func(t *testing.T) {
		sku := test.RandomSku(t, "")
		batchRef := test.RandomBatchRef(t, "")

		test.EntInsertBatch(t, client, batchRef, sku, 30, time.Time{})
		test.EntInsertAllocation(t, client, batchRef, test.RandomOrderId(t, ""), sku, 30)

		repo := ORMProductsRepository{
			client: client,
		}

		batch := domain.NewBatch(batchRef, sku, 100, time.Time{})

		orderLine := domain.OrderLine{
			OrderID:  test.RandomOrderId(t, ""),
			Quantity: 12,
			Sku:      domain.Sku(sku),
		}

		assert.Error(t, repo.AllocateToBatch(batch, orderLine))
	})
}

func TestORMProductsRepository_DeallocateFromBatch(t *testing.T) {
	client := test.EntClient(t)
	defer client.Close()

	repo := ORMProductsRepository{
		client: client,
	}

	sku := test.RandomSku(t, "")
	batchRef := test.RandomBatchRef(t, "")

	batch := domain.NewBatch(batchRef, sku, 40, time.Time{})
	orderLine := domain.OrderLine{
		OrderID:  test.RandomOrderId(t, ""),
		Quantity: 12,
		Sku:      domain.Sku(sku),
	}

	err := repo.AddBatch(batch)
	assert.NoError(t, err)

	err = repo.AllocateToBatch(batch, orderLine)
	assert.NoError(t, err)

	err = repo.DeallocateFromBatch(batch, orderLine)
	assert.NoError(t, err)

	deallocatedBatch, err := repo.GetBatch(batch.Reference)
	assert.NoError(t, err)

	assert.Len(t, deallocatedBatch.Allocations.ToSlice(), 0)
}
