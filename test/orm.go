package test

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/abbasegbeyemi/cosmic-python-go/domain"
	"github.com/abbasegbeyemi/cosmic-python-go/ent"
	"github.com/abbasegbeyemi/cosmic-python-go/ent/batch"
	"github.com/abbasegbeyemi/cosmic-python-go/ent/enttest"
	"github.com/abbasegbeyemi/cosmic-python-go/ent/migrate"
	_ "github.com/mattn/go-sqlite3"
)

func EntClient(t *testing.T) *ent.Client {
	t.Helper()
	testDBFile := t.TempDir() + "-test.sqlite"
	// testDBFile := "f-test.sqlite"

	opts := []enttest.Option{
		enttest.WithOptions(ent.Log(t.Log)),
		enttest.WithMigrateOptions(migrate.WithGlobalUniqueID(true)),
	}

	client := enttest.Open(t, "sqlite3", fmt.Sprintf("file:%s?cache=shared&_fk=1", testDBFile), opts...)
	return client
}

func EntGetBatch(t *testing.T, client *ent.Client, reference domain.Reference) domain.Batch {
	t.Helper()
	batch, err := client.Batch.Query().Where(batch.Reference(string(reference))).Only(context.Background())
	if err != nil {
		t.Fatalf("error occurred when getting batch from database: %s", err)
	}
	return domain.NewBatch(domain.Reference(batch.Reference), domain.Sku(batch.Sku), batch.Quantity, batch.Eta)
}

func EntInsertBatch(t *testing.T, client *ent.Client, reference domain.Reference, sku domain.Sku, quantity int, eta time.Time) {
	t.Helper()
	_, err := client.Batch.Create().SetReference(string(reference)).SetSku(string(sku)).SetEta(eta).SetQuantity(quantity).Save(context.Background())
	if err != nil {
		t.Fatalf("could not seed the database with a batch: %s", err)
	}
}

func EntInsertOrderLine(t *testing.T, client *ent.Client, orderId domain.Reference, sku domain.Sku, quantity int) {
	t.Helper()
	_, err := client.OrderLine.Create().SetOrderID(string(orderId)).SetSku(string(sku)).SetQuantity(quantity).Save(context.Background())
	if err != nil {
		t.Fatalf("could not seed the database with the order line: %s", err)
	}
}

func EntInsertAllocation(t *testing.T, client *ent.Client, batchRef, orderId domain.Reference, sku domain.Sku, quantity int) {
	t.Helper()
	batchId, err := client.Batch.Query().Where(batch.Reference(string(batchRef))).OnlyID(context.Background())
	if err != nil {
		t.Fatalf("could not find batch with ref: %s: %s", batchRef, err)
	}
	_, err = client.OrderLine.Create().SetBatchID(batchId).SetOrderID(string(orderId)).SetSku(string(sku)).SetQuantity(quantity).Save(context.Background())
	if err != nil {
		t.Fatalf("could insert allocation")
	}
}
