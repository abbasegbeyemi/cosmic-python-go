package repos

import (
	"fmt"

	"github.com/abbasegbeyemi/cosmic-python-go/domain"
)

type FakeRepository struct {
	Batches          []domain.Batch
	OrderLines       []domain.OrderLine
	BatchAllocations map[domain.Reference]domain.OrderLine
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
	f.BatchAllocations[batch.Reference] = orderLine
	return nil
}
