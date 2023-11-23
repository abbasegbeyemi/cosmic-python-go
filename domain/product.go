package domain

import (
	"fmt"
	"slices"
)

type Sku string
type Reference string

type Product struct {
	Sku     Sku
	Batches []Batch
}

func (p *Product) Allocate(orderLine OrderLine) (Reference, error) {
	slices.SortFunc[[]Batch](p.Batches, func(aBatch, bBatch Batch) int {
		return aBatch.ETA.Compare(bBatch.ETA)
	})

	for _, batch := range p.Batches {
		if err := batch.Allocate(orderLine); err == nil {
			return batch.Reference, nil
		}

	}
	return "", OutOfStockError{orderLine.Sku}
}

type OutOfStockError struct {
	Sku Sku
}

func (o OutOfStockError) Error() string {
	return fmt.Sprintf("%s is out of stock", o.Sku)
}
