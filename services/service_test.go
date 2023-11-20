package services

import (
	"testing"
	"time"

	"github.com/abbasegbeyemi/cosmic-python-go/domain"
	"github.com/abbasegbeyemi/cosmic-python-go/repos"
	"github.com/abbasegbeyemi/cosmic-python-go/uow"
	"github.com/stretchr/testify/assert"
)

func TestService_Allocate(t *testing.T) {
	t.Run("returns allocation", func(t *testing.T) {
		batchRef := domain.Reference("batch-123")
		sku := domain.Sku("MASSIVE-LAMP")

		uow := uow.NewFakeUnitOfWork(repos.NewFakeRepository())
		service := StockService{
			UOW: uow,
		}

		service.AddBatch(batchRef, sku, 100, time.Time{})
		allocatedBatchRef, err := service.Allocate("order-1", sku, 12)

		assert.NoError(t, err)
		assert.Equal(t, batchRef, allocatedBatchRef)
	})

	t.Run("returns error for an invalid sku", func(t *testing.T) {
		invalidSku := domain.Sku("INVALID-SKU")

		uow := uow.NewFakeUnitOfWork(repos.NewFakeRepository(
			repos.WithBatch("batch-002", "VALID-SKU", 100, time.Time{}),
		))

		service := StockService{
			UOW: uow,
		}
		_, err := service.Allocate("order-1", invalidSku, 12)
		assert.ErrorIs(t, err, InvalidSkuError{sku: invalidSku})
	})

	t.Run("allocate prefers warehouse batches to shipments", func(t *testing.T) {

		inStockBatchRef := domain.Reference("in-stock-batch-001")
		shipmentBatchRef := domain.Reference("shipment-batch-001")
		sku := domain.Sku("RETRO-CLOCK")

		uow := uow.NewFakeUnitOfWork(repos.NewFakeRepository(
			repos.WithBatch(inStockBatchRef, sku, 100, time.Time{}),
			repos.WithBatch(shipmentBatchRef, sku, 100, time.Time{}.AddDate(0, 4, 1)),
		))

		service := StockService{
			UOW: uow,
		}

		_, err := service.Allocate("order-002", "RETRO-CLOCK", 10)
		assert.Nil(t, err)

		inStockBatch, _ := uow.Batches().GetBatch(inStockBatchRef)
		shipmentBatch, _ := uow.Batches().GetBatch(shipmentBatchRef)

		assert.Equal(t, 90, inStockBatch.AvailableQuantity())
		assert.Equal(t, 100, shipmentBatch.AvailableQuantity())
	})
}

func TestService_Deallocate(t *testing.T) {
	t.Run("should decrement available quantity", func(t *testing.T) {
		sku := domain.Sku("DISCONTINUED-LAMP")

		allocatedBatch := domain.Batch{Reference: "batch-123", Sku: sku, Quantity: 30, ETA: time.Time{}.AddDate(2025, 10, 2)}
		allocatedOrderLine := domain.OrderLine{OrderID: "order001", Sku: sku, Quantity: 10}

		uow := uow.NewFakeUnitOfWork(repos.NewFakeRepository())

		service := StockService{
			UOW: uow,
		}

		service.AddBatch(allocatedBatch.Reference, allocatedBatch.Sku, allocatedBatch.Quantity, allocatedBatch.ETA)

		batchRef, err := service.Allocate(allocatedOrderLine.OrderID, allocatedBatch.Sku, allocatedOrderLine.Quantity)
		assert.Nil(t, err)

		batch, err := service.UOW.Batches().GetBatch(batchRef)
		assert.Nil(t, err)

		assert.Equal(t, (allocatedBatch.Quantity - allocatedOrderLine.Quantity), batch.AvailableQuantity())

		err = service.Deallocate(batch, allocatedOrderLine)
		assert.Nil(t, err)

		deallocatedBatch, err := service.UOW.Batches().GetBatch(batchRef)
		assert.Nil(t, err)

		assert.Equal(t, deallocatedBatch.Quantity, deallocatedBatch.AvailableQuantity())
	})

	t.Run("should return error when deallocating an unallocated batch", func(t *testing.T) {
		sku := domain.Sku("DISCONTINUED-LAMP")

		batch := domain.Batch{Reference: "batch-123", Sku: sku, Quantity: 30, ETA: time.Time{}.AddDate(2025, 10, 2)}
		orderLine := domain.OrderLine{OrderID: "order001", Sku: sku, Quantity: 10}

		uow := uow.NewFakeUnitOfWork(repos.NewFakeRepository())

		service := StockService{
			UOW: uow,
		}
		service.AddBatch(batch.Reference, batch.Sku, batch.Quantity, batch.ETA)

		err := service.Deallocate(batch, orderLine)
		assert.Error(t, err)
	})
}

func TestService_AddBatch(t *testing.T) {
	uow := uow.NewFakeUnitOfWork(repos.NewFakeRepository())

	service := StockService{
		UOW: uow,
	}

	assert.NoError(t, service.AddBatch("batch-001", "CRUNCHY-NUT", 40, time.Time{}))

	_, err := uow.Batches().GetBatch("batch-001")
	assert.NoError(t, err)

	assert.True(t, uow.Committed)
}
