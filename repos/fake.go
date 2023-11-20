package repos

import (
	"fmt"
	"slices"

	"github.com/abbasegbeyemi/cosmic-python-go/domain"
)

type FakeRepository struct {
	Batches          []domain.Batch
	OrderLines       []domain.OrderLine
	BatchAllocations map[domain.Reference][]domain.OrderLine
}

func (f *FakeRepository) AddBatch(batch domain.Batch) error {
	f.Batches = append(f.Batches, batch)
	return nil
}

func (f *FakeRepository) ListBatches() ([]domain.Batch, error) {
	return f.Batches, nil
}

func (f *FakeRepository) GetBatch(reference domain.Reference) (domain.Batch, error) {
	for _, batch := range f.Batches {
		if batch.Reference == reference {
			return batch, nil
		}
	}
	return domain.Batch{}, fmt.Errorf("could not find requested batch")
}

func (f *FakeRepository) AddOrderLine(orderLine domain.OrderLine) error {
	f.OrderLines = append(f.OrderLines, orderLine)
	return nil
}

func (f *FakeRepository) AllocateToBatch(batch domain.Batch, orderLine domain.OrderLine) error {
	_, ok := f.BatchAllocations[batch.Reference]
	batch.Allocate(orderLine)
	if !ok {
		f.BatchAllocations[batch.Reference] = []domain.OrderLine{orderLine}
		return nil
	}
	f.BatchAllocations[batch.Reference] = append(f.BatchAllocations[batch.Reference], orderLine)
	return nil
}

func (f *FakeRepository) DeallocateFromBatch(batch domain.Batch, orderLine domain.OrderLine) error {

	allocatedOrderLines, ok := f.BatchAllocations[batch.Reference]
	if !ok {
		return fmt.Errorf("this batch has no allocations")
	}

	batch.Deallocate(orderLine)
	orderLineIndex := slices.IndexFunc[[]domain.OrderLine](allocatedOrderLines, func(ol domain.OrderLine) bool {
		return ol.OrderID == orderLine.OrderID
	})
	if orderLineIndex == -1 {
		return fmt.Errorf("this order line has not been allocated to this batch")
	}
	// Override the order line with the one at the end
	allocatedOrderLines[orderLineIndex] = allocatedOrderLines[len(allocatedOrderLines)-1]

	// Use the list of order lines barring the last one
	f.BatchAllocations[batch.Reference] = allocatedOrderLines[:len(allocatedOrderLines)-1]

	return nil
}

func NewFakeRepository() *FakeRepository {
	return &FakeRepository{
		BatchAllocations: make(map[domain.Reference][]domain.OrderLine),
	}
}
