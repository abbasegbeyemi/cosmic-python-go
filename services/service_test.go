package services

import (
	"testing"
	"time"

	"github.com/abbasegbeyemi/cosmic-python-go/domain"
	"github.com/abbasegbeyemi/cosmic-python-go/repos"
	"github.com/stretchr/testify/assert"
)

func TestService_Allocate(t *testing.T) {
	t.Run("returns allocation", func(t *testing.T) {
		batchRef := domain.Reference("batch-123")
		sku := domain.Sku("MASSIVE-LAMP")
		line := domain.OrderLine{
			OrderID:  "order-1",
			Sku:      sku,
			Quantity: 12,
		}
		batch := domain.NewBatch(batchRef, sku, 100, time.Now())
		repo := repos.NewFakeRepository()
		repo.Batches = []domain.Batch{batch}

		service := StockService{
			repo: repo,
		}
		allocatedBatchRef, err := service.Allocate(line)
		assert.Nil(t, err)
		assert.Equal(t, batchRef, allocatedBatchRef)
	})

	t.Run("returns error for an invalid sku", func(t *testing.T) {
		batchRef := domain.Reference("batch-123")
		invalidSku := domain.Sku("INVALID-SKU")
		line := domain.OrderLine{
			OrderID:  "order-1",
			Sku:      invalidSku,
			Quantity: 12,
		}
		batch := domain.NewBatch(batchRef, "VALID-SKU", 100, time.Now())
		repo := repos.NewFakeRepository()
		repo.Batches = []domain.Batch{batch}
		service := StockService{
			repo: repo,
		}
		_, err := service.Allocate(line)
		assert.ErrorIs(t, err, InvalidSkuError{sku: invalidSku})
	})
}

func TestService_Deallocate(t *testing.T) {
	t.Run("should decrement available quantity", func(t *testing.T) {
		sku := domain.Sku("DISCONTINUED-LAMP")

		allocatedBatch := domain.Batch{Reference: "batch-123", Sku: sku, Quantity: 30, ETA: time.Time{}.AddDate(2025, 10, 2)}
		allocatedOrderLine := domain.OrderLine{OrderID: "order001", Sku: sku, Quantity: 10}

		repo := repos.NewFakeRepository()

		service := StockService{repo: repo}
		service.AddBatch(allocatedBatch.Reference, allocatedBatch.Sku, allocatedBatch.Quantity, allocatedBatch.ETA)

		batchRef, err := service.Allocate(allocatedOrderLine)
		assert.Nil(t, err)

		batch, err := repo.GetBatch(batchRef)
		assert.Nil(t, err)

		assert.Equal(t, (allocatedBatch.Quantity - allocatedOrderLine.Quantity), batch.AvailableQuantity())

		err = service.Deallocate(batch, allocatedOrderLine)
		assert.Nil(t, err)

		deallocatedBatch, err := repo.GetBatch(batchRef)
		assert.Nil(t, err)

		assert.Equal(t, deallocatedBatch.Quantity, deallocatedBatch.AvailableQuantity())
	})

	t.Run("should return error when deallocating an unallocated batch", func(t *testing.T) {
		sku := domain.Sku("DISCONTINUED-LAMP")

		batch := domain.Batch{Reference: "batch-123", Sku: sku, Quantity: 30, ETA: time.Time{}.AddDate(2025, 10, 2)}
		orderLine := domain.OrderLine{OrderID: "order001", Sku: sku, Quantity: 10}

		repo := repos.NewFakeRepository()

		service := StockService{repo: repo}
		service.AddBatch(batch.Reference, batch.Sku, batch.Quantity, batch.ETA)

		err := service.Deallocate(batch, orderLine)
		assert.Error(t, err)
	})
}

func TestService_AddBatch(t *testing.T) {
	batchToAdd := domain.NewBatch("batch-001", "LARGE-TABLE", 30, time.Time{})
	repo := repos.NewFakeRepository()
	service := StockService{
		repo: repo,
	}
	err := service.AddBatch(batchToAdd.Reference, batchToAdd.Sku, batchToAdd.Quantity, batchToAdd.ETA)
	assert.Nil(t, err)

	addedBatch, err := repo.GetBatch(batchToAdd.Reference)
	assert.Nil(t, err)

	assert.EqualExportedValues(t, batchToAdd, addedBatch)
}
