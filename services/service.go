package services

import (
	"fmt"
	"time"

	"github.com/abbasegbeyemi/cosmic-python-go/apperrors"
	"github.com/abbasegbeyemi/cosmic-python-go/domain"
	"github.com/abbasegbeyemi/cosmic-python-go/uow"
)

type UnitOfWork interface {
	Products() uow.ProductRepository
	Commit() error
	Rollback()
	CommitOnSuccess(uow.QueryFunc) error
}

type StockService struct {
	UOW UnitOfWork
}

func (s *StockService) AddBatch(reference domain.Reference, sku domain.Sku, quantity int, eta time.Time) error {
	return s.UOW.CommitOnSuccess(func() error {
		product, getProductErr := s.UOW.Products().Get(sku)
		nonExistentProductErr := apperrors.NonExistentProductError{Sku: sku}
		if getProductErr == nonExistentProductErr {
			product = domain.Product{Sku: sku}
			if addProductErr := s.UOW.Products().Add(product); addProductErr != nil {
				return fmt.Errorf("could not create nonexistent product")
			}
		}
		batch := domain.NewBatch(reference, sku, quantity, eta)
		if err := s.UOW.Products().AddBatch(batch); err != nil {
			return fmt.Errorf("could not add batch")
		}
		product.Batches = append(product.Batches, batch)
		return nil
	})
}

func (s *StockService) Allocate(orderId domain.Reference, sku domain.Sku, quantity int) (domain.Reference, error) {
	line := domain.OrderLine{OrderID: orderId, Sku: sku, Quantity: quantity}
	var batchRef domain.Reference
	allocateError := s.UOW.CommitOnSuccess(func() error {
		if err := s.UOW.Products().AddOrderLine(line); err != nil {
			return fmt.Errorf("could not create order line")
		}
		product, getProductError := s.UOW.Products().Get(sku)
		nonExistentProductErr := apperrors.NonExistentProductError{Sku: sku}
		if getProductError == nonExistentProductErr {
			return InvalidSkuError{Sku: sku}
		}
		var err error
		if batchRef, err = product.Allocate(line); err != nil {
			return fmt.Errorf("could not allocate order to a batch: %w", err)
		}
		batchToAllocate, err := s.UOW.Products().GetBatch(batchRef)
		if err != nil {
			return fmt.Errorf("could not get batch marked for allocation")
		}
		if err = s.UOW.Products().AllocateToBatch(batchToAllocate, line); err != nil {
			return fmt.Errorf("error completing allocation to batch")
		}
		return nil
	})
	return batchRef, allocateError

}

func (s *StockService) Deallocate(batch domain.Batch, orderLine domain.OrderLine) error {
	batchEnriched, err := s.UOW.Products().GetBatch(batch.Reference)
	if err != nil {
		return fmt.Errorf("could not retrieve batch")
	}
	if isAllocated := batchEnriched.IsAllocated(orderLine); !isAllocated {
		return fmt.Errorf("order line is not allocated to this batch")
	}
	return s.UOW.Products().DeallocateFromBatch(batch, orderLine)
}

type InvalidSkuError struct {
	Sku domain.Sku
}

func (i InvalidSkuError) Error() string {
	return fmt.Sprintf("%s sku is invalid", i.Sku)
}
