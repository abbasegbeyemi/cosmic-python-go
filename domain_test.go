package cosmicpythongo

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestBatch_AvailableQuantity(t *testing.T) {
	batch := NewBatch("batch-001", "SMALL-TABLE", 5)
	assert.Equal(t, batch.quantity, batch.AvailableQuantity())
}

func TestBatch_AllocatedQuantity(t *testing.T) {
	batch := NewBatch("batch-001", "SMALL-TABLE", 5)
	batch.Allocate(OrderLine{
		OrderID:  "order-001",
		Sku:      "SMALL-TABLE",
		Quantity: 2},
	)
	assert.Equal(t, 2, batch.AllocatedQuantity())
}

func TestBatch_Allocate(t *testing.T) {
	t.Run("should be able to allocate an order to a batch", func(t *testing.T) {
		batch := NewBatch("batch-001", "SMALL-TABLE", 5)
		orderLine := OrderLine{
			OrderID:  "order-ref",
			Sku:      "SMALL-TABLE",
			Quantity: 3,
		}

		err := batch.Allocate(orderLine)
		assert.Nil(t, err)
	})

	t.Run("order line allocated to batch decreases available quantity", func(t *testing.T) {
		batch := NewBatch("batch-001", "SMALL-TABLE", 5)
		batch.ETA = time.Now()

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
		batch := NewBatch("batch-001", "SMALL-TABLE", 10)
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
		batch := NewBatch("batch-001", "SMALL-TABLE", 5)
		orderLine := OrderLine{
			OrderID:  "order-ref",
			Sku:      "BIG-TABLE",
			Quantity: 3,
		}

		canAllocate, _ := batch.CanAllocate(orderLine)
		assert.False(t, canAllocate)
	})

	t.Run("should not be able to allocate if available quantity is less than order quantity", func(t *testing.T) {
		batch := NewBatch("batch-001", "SMALL-TABLE", 10)
		orderLine := OrderLine{
			OrderID:  "order-ref",
			Sku:      "SMALL-TABLE",
			Quantity: 11,
		}

		canAllocate, _ := batch.CanAllocate(orderLine)
		assert.False(t, canAllocate)
	})

	t.Run("same order allocated twice should not decrease available quantity", func(t *testing.T) {
		batch := NewBatch("batch-001", "SMALL-TABLE", 10)

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
	batch := NewBatch("batch-001", "SMALL-TABLE", 5)
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
		batch := NewBatch("batch-001", "SMALL-TABLE", 5)

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
		batch := NewBatch("batch-001", "SMALL-TABLE", 5)

		orderLine := OrderLine{
			OrderID:  "order-001",
			Sku:      "SMALL-TABLE",
			Quantity: 3,
		}

		orderAllocated := batch.IsAllocated(orderLine)
		assert.False(t, orderAllocated)
	})

}
