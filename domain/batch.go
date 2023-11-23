package domain

import (
	"fmt"
	"time"

	mapset "github.com/deckarep/golang-set/v2"
)

type Batch struct {
	Reference   Reference
	Sku         Sku
	Quantity    int
	ETA         time.Time
	Allocations mapset.Set[OrderLine]
}

func NewBatch(reference Reference, sku Sku, Quantity int, eta time.Time) Batch {
	return Batch{
		Reference:   reference,
		Sku:         sku,
		Quantity:    Quantity,
		ETA:         eta,
		Allocations: mapset.NewSet[OrderLine](),
	}
}

// Allocate allocates an order line to a batch
func (b *Batch) Allocate(orderLine OrderLine) error {
	canAllocate, reason := b.CanAllocate(orderLine)
	if !canAllocate {
		return reason
	}
	b.Allocations.Add(orderLine)
	return nil
}

// Deallocate removes an order line from a batch
func (b *Batch) Deallocate(orderLine OrderLine) {
	b.Allocations.Remove(orderLine)
}

// CanAllocate returns true if an order can be allocated to the batch and the reason why not if false
func (b *Batch) CanAllocate(orderLine OrderLine) (bool, error) {
	if b.Sku != orderLine.Sku {
		return false, fmt.Errorf("order of %s cannot be allocated to a batch of %s", orderLine.Sku, b.Sku)
	}

	if b.AvailableQuantity() < orderLine.Quantity {
		return false, fmt.Errorf("unable to allocate order to batch, not enough %s left", b.Sku)
	}

	if b.IsAllocated(orderLine) {
		return false, fmt.Errorf("order already allocated")
	}

	return true, nil

}

// AvailableQuantity returns the number of product left after accounting for orders
func (b *Batch) AvailableQuantity() int {
	var orders int
	for order := range b.Allocations.Iter() {
		orders += order.Quantity
	}
	return b.Quantity - orders
}

// AllocatedQuantity returns the quantity that has been allocated to orders
func (b *Batch) AllocatedQuantity() int {
	return b.Quantity - b.AvailableQuantity()
}

// IsAllocated checks if an order line has been allocated to the batch
func (b *Batch) IsAllocated(orderLine OrderLine) bool {
	return b.Allocations.Contains(orderLine)
}
