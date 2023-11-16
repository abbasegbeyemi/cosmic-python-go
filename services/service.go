package services

import (
	"fmt"
	"slices"

	"github.com/abbasegbeyemi/cosmic-python-go/domain"
)

type Repository interface {
	ListBatches() ([]domain.Batch, error)
	GetBatch(reference domain.Reference) (domain.Batch, error)
	AllocateToBatch(domain.Batch, domain.OrderLine) error
	AddOrderLine(domain.OrderLine) error
}

type StockService struct {
	repo Repository
}

func NewStockService(repo Repository) StockService {
	return StockService{
		repo: repo,
	}
}

func (s *StockService) Allocate(orderLine domain.OrderLine) (domain.Reference, error) {
	batches, err := s.repo.ListBatches()

	if err != nil {
		return "", fmt.Errorf("could not list batches: %w", err)
	}

	if !s.isValidSku(orderLine.Sku, batches) {
		return "", InvalidSkuError{sku: orderLine.Sku}
	}

	batchRef, err := domain.Allocate(orderLine, batches)

	if err != nil {
		return "", fmt.Errorf("could not allocate order line to any batch: %w", err)
	}

	if err = s.repo.AddOrderLine(orderLine); err != nil {
		return "", fmt.Errorf("could not add order line: %w", err)
	}

	batchToAllocate, err := s.repo.GetBatch(batchRef)

	if err != nil {
		return "", fmt.Errorf("could not find batch to allocate order line to: %w", err)
	}

	if err = s.repo.AllocateToBatch(batchToAllocate, orderLine); err != nil {
		return "", fmt.Errorf("could not persist order line allocation")
	}
	return batchRef, nil
}

type InvalidSkuError struct {
	sku domain.Sku
}

func (i InvalidSkuError) Error() string {
	return fmt.Sprintf("%s sku is invalid", i.sku)
}

func (s StockService) isValidSku(sku domain.Sku, batches []domain.Batch) bool {
	return slices.ContainsFunc[[]domain.Batch](batches, func(batch domain.Batch) bool {
		return batch.Sku == sku
	})
}
