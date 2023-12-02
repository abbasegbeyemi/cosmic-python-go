package repos

import (
	"context"
	"fmt"

	"github.com/abbasegbeyemi/cosmic-python-go/domain"
	"github.com/abbasegbeyemi/cosmic-python-go/ent"
	"github.com/abbasegbeyemi/cosmic-python-go/ent/batch"
	"github.com/abbasegbeyemi/cosmic-python-go/ent/orderline"
)

type ORMProductsRepository struct {
	client *ent.Client
}

func (o *ORMProductsRepository) AddBatch(batch domain.Batch) error {
	_, err := o.client.Batch.
		Create().
		SetReference(string(batch.Reference)).
		SetSku(string(batch.Sku)).
		SetQuantity(batch.Quantity).
		SetEta(batch.ETA).
		Save(context.Background())

	if err != nil {
		return fmt.Errorf("could not create batch: %w", err)
	}
	return nil
}

func (o ORMProductsRepository) GetBatch(reference domain.Reference) (domain.Batch, error) {
	entBatch, err := o.client.Batch.Query().
		Select(batch.FieldReference, batch.FieldSku, batch.FieldQuantity, batch.FieldEta).
		Where(batch.ReferenceEQ(string(reference))).Only(context.Background())
	if err != nil {
		return domain.Batch{}, fmt.Errorf("could not get batch: %w", err)
	}

	domainBatch := domain.NewBatch(domain.Reference(entBatch.Reference), domain.Sku(entBatch.Sku), entBatch.Quantity, entBatch.Eta)

	return entEnrichAllocations(o.client, entBatch, domainBatch)
}

func entEnrichAllocations(client *ent.Client, entBatch *ent.Batch, domainBatch domain.Batch) (domain.Batch, error) {
	allocations, err := client.Batch.QueryOrderLines(entBatch).All(context.Background())
	if err != nil {
		return domainBatch, fmt.Errorf("could not get batch allocations: %w", err)
	}

	for _, allocation := range allocations {
		domainBatch.Allocations.Add(domain.OrderLine{
			OrderID:  domain.Reference(allocation.OrderID),
			Sku:      domain.Sku(allocation.Sku),
			Quantity: allocation.Quantity,
		})
	}
	return domainBatch, nil
}

func (o *ORMProductsRepository) ListBatches(sku domain.Sku) ([]domain.Batch, error) {
	entBatches, err := o.client.Batch.Query().Select(batch.FieldReference, batch.FieldSku, batch.FieldQuantity, batch.FieldEta).All(context.Background())
	if err != nil {
		return []domain.Batch{}, fmt.Errorf("could not query batches: %w", err)
	}
	var enrichedBatches []domain.Batch
	for _, entBatch := range entBatches {
		domainBatch := domain.NewBatch(domain.Reference(entBatch.Reference), domain.Sku(entBatch.Sku), entBatch.Quantity, entBatch.Eta)
		enrichedBatch, err := entEnrichAllocations(o.client, entBatch, domainBatch)
		if err != nil {
			return enrichedBatches, fmt.Errorf("could not enrich allocations for batch ref: %s: %w", domainBatch.Reference, err)
		}
		enrichedBatches = append(enrichedBatches, enrichedBatch)
	}
	return enrichedBatches, nil
}

func entGetBatchIdFromReference(client *ent.Client, reference domain.Reference) (int, error) {
	batchId, err := client.Batch.Query().Select().Where(batch.ReferenceEQ(string(reference))).OnlyID(context.Background())
	if err != nil {
		return 0, err
	}
	return batchId, nil
}

func (o *ORMProductsRepository) AllocateToBatch(domainBatch domain.Batch, domainOrderLine domain.OrderLine) error {
	batchEnriched, err := o.GetBatch(domainBatch.Reference)
	if err != nil {
		return fmt.Errorf("could not enrich batch: %w", err)
	}

	if err = batchEnriched.Allocate(domainOrderLine); err != nil {
		return fmt.Errorf("cannot allocate this order to this batch: %w", err)
	}

	batchId, err := entGetBatchIdFromReference(o.client, batchEnriched.Reference)
	if err != nil {
		return fmt.Errorf("could not get batch id: %w", err)
	}

	_, err = o.client.OrderLine.Create().SetOrderID(string(domainOrderLine.OrderID)).SetSku(string(domainOrderLine.Sku)).SetQuantity(domainOrderLine.Quantity).SetBatchID(batchId).Save(context.Background())
	if err != nil {
		return fmt.Errorf("failed to store allocation: %w", err)
	}
	return nil
}

func (o *ORMProductsRepository) DeallocateFromBatch(domainBatch domain.Batch, domainOrderLine domain.OrderLine) error {
	batchEnriched, err := o.GetBatch(domainBatch.Reference)

	if err != nil {
		return fmt.Errorf("could not enrich batch: %w", err)
	}
	batchEnriched.Deallocate(domainOrderLine)

	entOrderLine, err := o.client.OrderLine.Query().Where(orderline.OrderIDEQ(string(domainOrderLine.OrderID)), orderline.SkuEQ(string(domainOrderLine.Sku)), orderline.QuantityEQ(domainOrderLine.Quantity)).OnlyID(context.Background())
	if err != nil {
		return fmt.Errorf("could not find orderline")
	}

	_, err = o.client.Batch.Update().RemoveOrderLineIDs(entOrderLine).Save(context.Background())
	if err != nil {
		return fmt.Errorf("could not deallocate orderline: %w", err)
	}
	return nil
}
