package cosmicpythongo

import (
	"testing"
	"time"

	mapset "github.com/deckarep/golang-set/v2"
	"github.com/stretchr/testify/assert"
)

func TestBatch_AvailableQuantity(t *testing.T) {
	batch := Batch{reference: "batch-001", sku: "SMALL-TABLE", quantity: 5, allocations: mapset.NewSet[OrderLine]()}
	assert.Equal(t, batch.quantity, batch.AvailableQuantity())
}

func TestBatch_AllocatedQuantity(t *testing.T) {
	batch := Batch{reference: "batch-001", sku: "SMALL-TABLE", quantity: 5, allocations: mapset.NewSet[OrderLine]()}

	batch.Allocate(OrderLine{
		OrderID:  "order-001",
		Sku:      "SMALL-TABLE",
		Quantity: 2},
	)
	assert.Equal(t, 2, batch.AllocatedQuantity())
}

func TestBatch_Allocate(t *testing.T) {
	t.Run("should be able to allocate an order to a batch", func(t *testing.T) {
		batch := Batch{reference: "batch-001", sku: "SMALL-TABLE", quantity: 5, allocations: mapset.NewSet[OrderLine]()}
		orderLine := OrderLine{
			OrderID:  "order-ref",
			Sku:      "SMALL-TABLE",
			Quantity: 3,
		}

		err := batch.Allocate(orderLine)
		assert.Nil(t, err)
	})

	t.Run("order line allocated to batch decreases available quantity", func(t *testing.T) {
		batch := Batch{reference: "batch-001", sku: "SMALL-TABLE", quantity: 5, allocations: mapset.NewSet[OrderLine]()}
		batch.eta = time.Now()

		orderLine := OrderLine{
			OrderID:  "order-ref",
			Sku:      "SMALL-TABLE",
			Quantity: 3,
		}

		err := batch.Allocate(orderLine)
		assert.Nil(t, err)

		assert.Equal(t, batch.quantity-orderLine.Quantity, batch.AvailableQuantity())
	})

	t.Run("should be able to allocate a list of orders", func(t *testing.T) {
		batch := Batch{reference: "batch-001", sku: "SMALL-TABLE", quantity: 20, allocations: mapset.NewSet[OrderLine]()}
		err := batch.Allocate(OrderLine{
			OrderID:  "order-001",
			Quantity: 3,
			Sku:      "SMALL-TABLE",
		})
		assert.Nil(t, err)

		err = batch.Allocate(OrderLine{
			OrderID:  "order-002",
			Quantity: 5,
			Sku:      "SMALL-TABLE",
		})
		assert.Nil(t, err)

		assert.Equal(t, batch.quantity-3-5, batch.AvailableQuantity())
	})
}

func TestBatch_CanAllocate(t *testing.T) {
	t.Run("should not allocate order of different type to a batch", func(t *testing.T) {
		batch := Batch{reference: "batch-001", sku: "SMALL-TABLE", quantity: 5, allocations: mapset.NewSet[OrderLine]()}
		orderLine := OrderLine{
			OrderID:  "order-ref",
			Sku:      "BIG-TABLE",
			Quantity: 3,
		}

		canAllocate, _ := batch.CanAllocate(orderLine)
		assert.False(t, canAllocate)
	})

	t.Run("should not be able to allocate if available quantity is less than order quantity", func(t *testing.T) {
		batch := Batch{reference: "batch-001", sku: "SMALL-TABLE", quantity: 10, allocations: mapset.NewSet[OrderLine]()}
		orderLine := OrderLine{
			OrderID:  "order-ref",
			Sku:      "SMALL-TABLE",
			Quantity: 11,
		}

		canAllocate, _ := batch.CanAllocate(orderLine)
		assert.False(t, canAllocate)
	})

	t.Run("same order allocated twice should not decrease available quantity", func(t *testing.T) {
		batch := Batch{reference: "batch-001", sku: "SMALL-TABLE", quantity: 10, allocations: mapset.NewSet[OrderLine]()}

		orderLine := OrderLine{
			OrderID:  "order-ref",
			Sku:      "SMALL-TABLE",
			Quantity: 3,
		}

		err := batch.Allocate(orderLine)
		assert.Nil(t, err)

		canAllocate, _ := batch.CanAllocate(orderLine)

		assert.False(t, canAllocate)
	})

}

func TestBatch_Deallocate(t *testing.T) {
	batch := Batch{reference: "batch-001", sku: "SMALL-TABLE", quantity: 5, allocations: mapset.NewSet[OrderLine]()}
	orderLine := OrderLine{
		OrderID:  "order-001",
		Sku:      "SMALL-TABLE",
		Quantity: 3,
	}

	err := batch.Allocate(orderLine)
	assert.Nil(t, err)

	assert.Equal(t, batch.AllocatedQuantity(), orderLine.Quantity)

	batch.Deallocate(orderLine)
	assert.Equal(t, batch.AllocatedQuantity(), 0)

}

func TestBatch_OrderAllocated(t *testing.T) {
	t.Run("should return true if batch has order allocated", func(t *testing.T) {
		batch := Batch{reference: "batch-001", sku: "SMALL-TABLE", quantity: 5, allocations: mapset.NewSet[OrderLine]()}

		orderLine := OrderLine{
			OrderID:  "order-001",
			Sku:      "SMALL-TABLE",
			Quantity: 3,
		}

		err := batch.Allocate(orderLine)
		assert.Nil(t, err)

		orderAllocated := batch.IsAllocated(orderLine)
		assert.True(t, orderAllocated)
	})

	t.Run("should return false if batch has order allocated", func(t *testing.T) {
		batch := Batch{reference: "batch-001", sku: "SMALL-TABLE", quantity: 5, allocations: mapset.NewSet[OrderLine]()}

		orderLine := OrderLine{
			OrderID:  "order-001",
			Sku:      "SMALL-TABLE",
			Quantity: 3,
		}

		orderAllocated := batch.IsAllocated(orderLine)
		assert.False(t, orderAllocated)
	})

}

func TestAllocate(t *testing.T) {
	t.Run("allocate prefers current stock batches to shipments", func(t *testing.T) {
		inStockBatch := Batch{reference: "in-stock-batch-001", sku: "RETRO-CLOCK", quantity: 100, allocations: mapset.NewSet[OrderLine]()}
		shipmentBatch := Batch{reference: "shipment-batch-001", sku: "RETRO-CLOCK", quantity: 100, allocations: mapset.NewSet[OrderLine](), eta: time.Time{}.AddDate(0, 4, 1)}

		line := OrderLine{
			OrderID:  "order-002",
			Sku:      "RETRO-CLOCK",
			Quantity: 10,
		}

		_, err := Allocate(line, []Batch{shipmentBatch, inStockBatch})
		assert.Nil(t, err)

		assert.Equal(t, 90, inStockBatch.AvailableQuantity())
		assert.Equal(t, 100, shipmentBatch.AvailableQuantity())
	})

	t.Run("allocate prefers earlier batches to later", func(t *testing.T) {
		earliestBatch := Batch{reference: "earliest-batch-001", sku: "RETRO-CLOCK", quantity: 100, allocations: mapset.NewSet[OrderLine](), eta: time.Time{}.AddDate(0, 2, 0)}
		mediumBatch := Batch{reference: "medium-batch-001", sku: "RETRO-CLOCK", quantity: 100, allocations: mapset.NewSet[OrderLine](), eta: time.Time{}.AddDate(0, 4, 1)}
		laterBatch := Batch{reference: "later-batch-001", sku: "RETRO-CLOCK", quantity: 100, allocations: mapset.NewSet[OrderLine](), eta: time.Time{}.AddDate(1, 0, 1)}

		orderLine := OrderLine{
			OrderID:  "order-123",
			Sku:      "RETRO-CLOCK",
			Quantity: 12,
		}
		_, err := Allocate(orderLine, []Batch{laterBatch, earliestBatch, mediumBatch})

		assert.Nil(t, err)

		assert.Equal(t, 88, earliestBatch.AvailableQuantity())
		assert.Equal(t, 100, mediumBatch.AvailableQuantity())
		assert.Equal(t, 100, laterBatch.AvailableQuantity())
	})

	t.Run("allocate returns allocated batch ref", func(t *testing.T) {
		inStockBatch := Batch{reference: "in-stock-batch-001", sku: "TEDDY-BEAR", quantity: 100, allocations: mapset.NewSet[OrderLine]()}
		shipmentBatch := Batch{reference: "shipment-batch-001", sku: "TEDDY-BEAR", quantity: 100, allocations: mapset.NewSet[OrderLine](), eta: time.Time{}.AddDate(0, 4, 1)}

		line := OrderLine{
			OrderID:  "order-002",
			Sku:      "TEDDY-BEAR",
			Quantity: 10,
		}

		batchRef, err := Allocate(line, []Batch{shipmentBatch, inStockBatch})
		assert.Nil(t, err)

		assert.Equal(t, Reference("in-stock-batch-001"), batchRef)

		assert.Equal(t, 90, inStockBatch.AvailableQuantity())
		assert.Equal(t, 100, shipmentBatch.AvailableQuantity())
	})

	t.Run("allocate will allocate to the soonest available batch", func(t *testing.T) {
		earliestBatch := Batch{reference: "earliest-batch-001", sku: "RETRO-CLOCK", quantity: 4, allocations: mapset.NewSet[OrderLine](), eta: time.Time{}.AddDate(0, 2, 0)}
		mediumBatch := Batch{reference: "medium-batch-001", sku: "RETRO-CLOCK", quantity: 100, allocations: mapset.NewSet[OrderLine](), eta: time.Time{}.AddDate(0, 4, 1)}
		laterBatch := Batch{reference: "later-batch-001", sku: "RETRO-CLOCK", quantity: 100, allocations: mapset.NewSet[OrderLine](), eta: time.Time{}.AddDate(1, 0, 1)}

		orderLine := OrderLine{
			OrderID:  "order-123",
			Sku:      "RETRO-CLOCK",
			Quantity: 12,
		}
		allocatedRef, err := Allocate(orderLine, []Batch{laterBatch, earliestBatch, mediumBatch})

		assert.Nil(t, err)
		assert.Equal(t, Reference("medium-batch-001"), allocatedRef)

		assert.Equal(t, 4, earliestBatch.AvailableQuantity())
		assert.Equal(t, 88, mediumBatch.AvailableQuantity())
		assert.Equal(t, 100, laterBatch.AvailableQuantity())
	})

	t.Run("allocate returns error if unable to allocate", func(t *testing.T) {
		batch := Batch{reference: "earliest-batch-001", sku: "RETRO-CLOCK", quantity: 4, allocations: mapset.NewSet[OrderLine](), eta: time.Time{}.AddDate(0, 2, 0)}
		orderLine := OrderLine{
			OrderID:  "order-123",
			Sku:      "RETRO-CLOCK",
			Quantity: 12,
		}
		_, err := Allocate(orderLine, []Batch{batch})
		assert.Error(t, err)
		assert.ErrorIs(t, OutOfStockError{"RETRO-CLOCK"}, err)
	})
}
