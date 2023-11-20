package services

import (
	"fmt"
	"slices"
	"time"

	"github.com/abbasegbeyemi/cosmic-python-go/domain"
	"github.com/abbasegbeyemi/cosmic-python-go/uow"
)

type UnitOfWork interface {
	Batches() uow.Repository
	Commit() error
	Rollback()
	DBInstruction(uow.QueryFunc) error
}

type StockService struct {
	UOW UnitOfWork
}

func (s *StockService) AddBatch(reference domain.Reference, sku domain.Sku, quantity int, eta time.Time) error {
	return s.UOW.DBInstruction(func() error {
		return s.UOW.Batches().AddBatch(domain.NewBatch(reference, sku, quantity, eta))
	})
}

func (s *StockService) Allocate(orderId domain.Reference, sku domain.Sku, quantity int) (domain.Reference, error) {
	batches, err := s.UOW.Batches().ListBatches()

	if err != nil {
		return "", fmt.Errorf("could not list batches: %w", err)
	}

	orderLine := domain.OrderLine{
		OrderID:  orderId,
		Sku:      sku,
		Quantity: quantity,
	}

	if !s.isValidSku(orderLine.Sku, batches) {
		return "", InvalidSkuError{sku: orderLine.Sku}
	}

	batchRef, err := domain.Allocate(orderLine, batches)

	if err != nil {
		return "", fmt.Errorf("could not allocate order line to any batch: %w", err)
	}

	if err = s.UOW.Batches().AddOrderLine(orderLine); err != nil {
		return "", fmt.Errorf("could not add order line: %w", err)
	}

	batchToAllocate, err := s.UOW.Batches().GetBatch(batchRef)

	if err != nil {
		return "", fmt.Errorf("could not find batch to allocate order line to: %w", err)
	}

	if err = s.UOW.Batches().AllocateToBatch(batchToAllocate, orderLine); err != nil {
		return "", fmt.Errorf("could not persist order line allocation")
	}
	return batchRef, nil
}

func (s *StockService) Deallocate(batch domain.Batch, orderLine domain.OrderLine) error {
	batchEnriched, err := s.UOW.Batches().GetBatch(batch.Reference)
	if err != nil {
		return fmt.Errorf("could not retrieve batch")
	}
	if isAllocated := batchEnriched.IsAllocated(orderLine); !isAllocated {
		return fmt.Errorf("order line is not allocated to this batch")
	}
	return s.UOW.Batches().DeallocateFromBatch(batch, orderLine)
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
