package repos

import (
	"fmt"
	"slices"
	"time"

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

	allocatedOrderLines := f.BatchAllocations[batch.Reference]

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

func NewFakeRepository(options ...func(*FakeRepository)) *FakeRepository {
	repo := &FakeRepository{
		BatchAllocations: make(map[domain.Reference][]domain.OrderLine),
	}
	for _, o := range options {
		o(repo)
	}
	return repo
}

func WithBatch(ref domain.Reference, sku domain.Sku, quantity int, eta time.Time) func(*FakeRepository) {
	return func(f *FakeRepository) {
		f.Batches = append(f.Batches, domain.NewBatch(ref, sku, quantity, eta))
	}
}

func WithOrderLine(orderId domain.Reference, sku domain.Sku, quantity int) func(*FakeRepository) {
	return func(f *FakeRepository) {
		f.OrderLines = append(f.OrderLines, domain.OrderLine{OrderID: orderId, Sku: sku, Quantity: quantity})
	}
}
