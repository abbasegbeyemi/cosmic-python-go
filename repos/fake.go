package repos

import (
	"fmt"
	"slices"
	"time"

	"github.com/abbasegbeyemi/cosmic-python-go/apperrors"
	"github.com/abbasegbeyemi/cosmic-python-go/domain"
	mapset "github.com/deckarep/golang-set/v2"
)

type FakeProductsRepository struct {
	products         mapset.Set[*domain.Product]
	batches          []domain.Batch
	orderLines       []domain.OrderLine
	batchAllocations map[domain.Reference][]domain.OrderLine
}

func (f *FakeProductsRepository) Get(sku domain.Sku) (domain.Product, error) {
	for product := range f.products.Iter() {
		if product.Sku == sku {
			productPopulated := f.populateBatches(*product)
			return productPopulated, nil
		}
	}
	return domain.Product{}, apperrors.NonExistentProductError{Sku: sku}
}

func (f *FakeProductsRepository) populateBatches(product domain.Product) domain.Product {
	for _, batch := range f.batches {
		if batch.Sku == product.Sku {
			product.Batches = append(product.Batches, batch)
		}
	}
	return product
}

func (f *FakeProductsRepository) Add(product domain.Product) error {
	if wasAdded := f.products.Add(&product); !wasAdded {
		return fmt.Errorf("attempted to add a duplicate product")
	}
	return nil
}

func (f *FakeProductsRepository) AddBatch(batch domain.Batch) error {
	f.batches = append(f.batches, batch)
	return nil
}

func (f *FakeProductsRepository) ListBatches(sku domain.Sku) ([]domain.Batch, error) {
	skuBatches := []domain.Batch{}
	for _, batch := range f.batches {
		if batch.Sku == sku {
			skuBatches = append(skuBatches, batch)
		}
	}
	return skuBatches, nil
}

func (f *FakeProductsRepository) GetBatch(reference domain.Reference) (domain.Batch, error) {
	for _, batch := range f.batches {
		if batch.Reference == reference {
			return batch, nil
		}
	}
	return domain.Batch{}, fmt.Errorf("could not find requested batch")
}

func (f *FakeProductsRepository) AddOrderLine(orderLine domain.OrderLine) error {
	f.orderLines = append(f.orderLines, orderLine)
	return nil
}

func (f *FakeProductsRepository) AllocateToBatch(batch domain.Batch, orderLine domain.OrderLine) error {
	_, ok := f.batchAllocations[batch.Reference]
	batch.Allocate(orderLine)
	if !ok {
		f.batchAllocations[batch.Reference] = []domain.OrderLine{orderLine}
		return nil
	}
	f.batchAllocations[batch.Reference] = append(f.batchAllocations[batch.Reference], orderLine)
	return nil
}

func (f *FakeProductsRepository) DeallocateFromBatch(batch domain.Batch, orderLine domain.OrderLine) error {
	allocatedOrderLines := f.batchAllocations[batch.Reference]

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
	f.batchAllocations[batch.Reference] = allocatedOrderLines[:len(allocatedOrderLines)-1]

	return nil
}

type FakeProductsRepositoryOptions func(*FakeProductsRepository)

// Construct a FakeProductsRepository
func NewFakeProductsRepository(options ...FakeProductsRepositoryOptions) *FakeProductsRepository {
	repo := &FakeProductsRepository{
		products:         mapset.NewSet[*domain.Product](),
		batchAllocations: make(map[domain.Reference][]domain.OrderLine),
	}
	for _, o := range options {
		o(repo)
	}
	return repo
}

// Populate a FakeProductsRepository with a batch
func WithBatch(ref domain.Reference, sku domain.Sku, quantity int, eta time.Time) func(*FakeProductsRepository) {
	return func(f *FakeProductsRepository) {
		f.batches = append(f.batches, domain.NewBatch(ref, sku, quantity, eta))
	}
}

// Populate with product
func WithProduct(sku domain.Sku) func(*FakeProductsRepository) {
	return func(f *FakeProductsRepository) {
		f.products.Add(&domain.Product{Sku: sku})
	}
}

// Populate a fake repository with an order line
func WithOrderLine(orderId domain.Reference, sku domain.Sku, quantity int) func(*FakeProductsRepository) {
	return func(f *FakeProductsRepository) {
		f.orderLines = append(f.orderLines, domain.OrderLine{OrderID: orderId, Sku: sku, Quantity: quantity})
	}
}
