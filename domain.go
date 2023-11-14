package cosmicpythongo

import (
	"fmt"
	"slices"
	"time"

	mapset "github.com/deckarep/golang-set/v2"
)

type Sku string
type Reference string

type Batch struct {
	reference   Reference
	sku         Sku
	quantity    int
	eta         time.Time
	allocations mapset.Set[OrderLine]
}

func NewBatch(reference Reference, sku Sku, Quantity int, eta time.Time) Batch {
	return Batch{
		reference:   reference,
		sku:         sku,
		quantity:    Quantity,
		allocations: mapset.NewSet[OrderLine](),
	}
}

type OrderLine struct {
	OrderID  Reference
	Sku      Sku
	Quantity int
}

// Allocate allocates an order line to a batch
func (b *Batch) Allocate(orderLine OrderLine) error {
	canAllocate, reason := b.CanAllocate(orderLine)
	if !canAllocate {
		return reason
	}
	b.allocations.Add(orderLine)
	return nil
}

// Deallocate removes an order line from a batch
func (b *Batch) Deallocate(orderLine OrderLine) {
	b.allocations.Remove(orderLine)
}

// CanAllocate returns true if an order can be allocated to the batch and the reason why not if false
func (b *Batch) CanAllocate(orderLine OrderLine) (bool, error) {
	if b.sku != orderLine.Sku {
		return false, fmt.Errorf("order of %s cannot be allocated to a batch of %s", orderLine.Sku, b.sku)
	}

	if b.AvailableQuantity() < orderLine.Quantity {
		return false, fmt.Errorf("unable to allocate order to batch, not enough %s left", b.sku)
	}

	if b.IsAllocated(orderLine) {
		return false, fmt.Errorf("order already allocated")
	}

	return true, nil

}

// AvailableQuantity returns the number of product left after accounting for orders
func (b *Batch) AvailableQuantity() int {
	var orders int
	for order := range b.allocations.Iter() {
		orders += order.Quantity
	}
	return b.quantity - orders
}

// AllocatedQuantity returns the quantity that has been allocated to orders
func (b *Batch) AllocatedQuantity() int {
	return b.quantity - b.AvailableQuantity()
}

// IsAllocated checks if an order line has been allocated to the batch
func (b *Batch) IsAllocated(orderLine OrderLine) bool {
	return b.allocations.Contains(orderLine)
}

func Allocate(orderLine OrderLine, batches []Batch) (Reference, error) {
	slices.SortFunc[[]Batch](batches, func(aBatch, bBatch Batch) int {
		return aBatch.eta.Compare(bBatch.eta)
	})
	var batchCheckError error
	for _, batch := range batches {
		batchCheckError = batch.Allocate(orderLine)
		if batchCheckError == nil {
			return batch.reference, nil
		}
	}
	return "", OutOfStockError{orderLine.Sku}
}

type OutOfStockError struct {
	sku Sku
}

func (o OutOfStockError) Error() string {
	return fmt.Sprintf("%s is out of stock", o.sku)
}
