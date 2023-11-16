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
		repo := repos.FakeRepository{
			Batches:          []domain.Batch{batch},
			BatchAllocations: make(map[domain.Reference]domain.OrderLine),
		}
		service := StockService{
			repo: &repo,
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
		repo := repos.FakeRepository{
			Batches:          []domain.Batch{batch},
			BatchAllocations: make(map[domain.Reference]domain.OrderLine),
		}
		service := StockService{
			repo: &repo,
		}
		_, err := service.Allocate(line)
		assert.ErrorIs(t, err, InvalidSkuError{sku: invalidSku})
	})
}
