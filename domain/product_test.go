package domain

import (
	"testing"
	"time"

	mapset "github.com/deckarep/golang-set/v2"
	"github.com/stretchr/testify/assert"
)

func TestProduct_Allocate(t *testing.T) {
	t.Run("allocate prefers current stock batches to shipments", func(t *testing.T) {
		inStockBatch := Batch{Reference: "in-stock-batch-001", Sku: "RETRO-CLOCK", Quantity: 100, Allocations: mapset.NewSet[OrderLine]()}
		shipmentBatch := Batch{Reference: "shipment-batch-001", Sku: "RETRO-CLOCK", Quantity: 100, Allocations: mapset.NewSet[OrderLine](), ETA: time.Time{}.AddDate(0, 4, 1)}

		line := OrderLine{
			OrderID:  "order-002",
			Sku:      "RETRO-CLOCK",
			Quantity: 10,
		}
		product := Product{
			Batches: []Batch{shipmentBatch, inStockBatch},
		}
		_, err := product.Allocate(line)
		assert.Nil(t, err)

		assert.Equal(t, 90, inStockBatch.AvailableQuantity())
		assert.Equal(t, 100, shipmentBatch.AvailableQuantity())
	})

	t.Run("allocate prefers earlier batches to later", func(t *testing.T) {
		earliestBatch := Batch{Reference: "earliest-batch-001", Sku: "RETRO-CLOCK", Quantity: 100, Allocations: mapset.NewSet[OrderLine](), ETA: time.Time{}.AddDate(0, 2, 0)}
		mediumBatch := Batch{Reference: "medium-batch-001", Sku: "RETRO-CLOCK", Quantity: 100, Allocations: mapset.NewSet[OrderLine](), ETA: time.Time{}.AddDate(0, 4, 1)}
		laterBatch := Batch{Reference: "later-batch-001", Sku: "RETRO-CLOCK", Quantity: 100, Allocations: mapset.NewSet[OrderLine](), ETA: time.Time{}.AddDate(1, 0, 1)}

		line := OrderLine{
			OrderID:  "order-123",
			Sku:      "RETRO-CLOCK",
			Quantity: 12,
		}
		product := Product{
			Batches: []Batch{laterBatch, mediumBatch, earliestBatch},
		}
		_, err := product.Allocate(line)
		assert.Nil(t, err)

		assert.Equal(t, 88, earliestBatch.AvailableQuantity())
		assert.Equal(t, 100, mediumBatch.AvailableQuantity())
		assert.Equal(t, 100, laterBatch.AvailableQuantity())
	})

	t.Run("allocate returns allocated batch ref", func(t *testing.T) {
		inStockBatch := Batch{Reference: "in-stock-batch-001", Sku: "TEDDY-BEAR", Quantity: 100, Allocations: mapset.NewSet[OrderLine]()}
		shipmentBatch := Batch{Reference: "shipment-batch-001", Sku: "TEDDY-BEAR", Quantity: 100, Allocations: mapset.NewSet[OrderLine](), ETA: time.Time{}.AddDate(0, 4, 1)}

		line := OrderLine{
			OrderID:  "order-002",
			Sku:      "TEDDY-BEAR",
			Quantity: 10,
		}

		product := Product{
			Batches: []Batch{shipmentBatch, inStockBatch},
		}
		batchRef, err := product.Allocate(line)
		assert.Nil(t, err)

		assert.Equal(t, Reference("in-stock-batch-001"), batchRef)

		assert.Equal(t, 90, inStockBatch.AvailableQuantity())
		assert.Equal(t, 100, shipmentBatch.AvailableQuantity())
	})

	t.Run("allocate will allocate to the soonest available batch", func(t *testing.T) {
		earliestBatch := Batch{Reference: "earliest-batch-001", Sku: "RETRO-CLOCK", Quantity: 4, Allocations: mapset.NewSet[OrderLine](), ETA: time.Time{}.AddDate(0, 2, 0)}
		mediumBatch := Batch{Reference: "medium-batch-001", Sku: "RETRO-CLOCK", Quantity: 100, Allocations: mapset.NewSet[OrderLine](), ETA: time.Time{}.AddDate(0, 4, 1)}
		laterBatch := Batch{Reference: "later-batch-001", Sku: "RETRO-CLOCK", Quantity: 100, Allocations: mapset.NewSet[OrderLine](), ETA: time.Time{}.AddDate(1, 0, 1)}

		line := OrderLine{
			OrderID:  "order-123",
			Sku:      "RETRO-CLOCK",
			Quantity: 12,
		}
		product := Product{
			Batches: []Batch{laterBatch, earliestBatch, mediumBatch},
		}
		allocatedRef, err := product.Allocate(line)
		assert.Nil(t, err)
		assert.Equal(t, Reference("medium-batch-001"), allocatedRef)

		assert.Equal(t, 4, earliestBatch.AvailableQuantity())
		assert.Equal(t, 88, mediumBatch.AvailableQuantity())
		assert.Equal(t, 100, laterBatch.AvailableQuantity())
	})

	t.Run("allocate returns error if unable to allocate", func(t *testing.T) {
		batch := Batch{Reference: "earliest-batch-001", Sku: "RETRO-CLOCK", Quantity: 4, Allocations: mapset.NewSet[OrderLine](), ETA: time.Time{}.AddDate(0, 2, 0)}
		line := OrderLine{
			OrderID:  "order-123",
			Sku:      "RETRO-CLOCK",
			Quantity: 12,
		}
		product := Product{
			Batches: []Batch{batch},
		}
		_, err := product.Allocate(line)
		assert.Error(t, err)
		assert.ErrorIs(t, OutOfStockError{"RETRO-CLOCK"}, err)
	})
}
