package repos

import (
	"database/sql"
	"slices"
	"testing"
	"time"

	"github.com/abbasegbeyemi/cosmic-python-go/domain"
	_ "github.com/mattn/go-sqlite3"
	"github.com/stretchr/testify/assert"
)

const createBatchesTableSQL string = `
	CREATE TABLE IF NOT EXISTS batches (
	reference STRING NOT NULL PRIMARY KEY,
	sku STRING NOT NULL,
	quantity INTEGER NOT NULL,
	eta DATETIME
	);
`

const createOrderLinesTableSQL string = `
	CREATE TABLE IF NOT EXISTS order_lines (
	order_id STRING NOT NULL,
	sku STRING NOT NULL,
	quantity INTEGER NOT NULL
	);
`

const createBatchesOrderLinesTableSQL string = `
    CREATE TABLE IF NOT EXISTS batches_order_lines (
    batch_id STRING NOT NULL,
    order_id STRING NOT NULL,
	FOREIGN KEY(batch_id) REFERENCES batches(reference)
    FOREIGN KEY(order_id) REFERENCES order_lines(order_id)
	PRIMARY KEY(batch_id, order_id)
    );
`

const truncateTablesSQL string = `
	DELETE FROM batches;
	DELETE FROM order_lines;
	DELETE FROM batches_order_lines;
`

const testDBFile string = "orders_test.sqlite"

func createTables(t *testing.T, db *sql.DB) {
	t.Helper()
	if _, err := db.Exec(createBatchesTableSQL); err != nil {
		t.Fatalf("could not create batches table %s", err)
	}
	if _, err := db.Exec(createOrderLinesTableSQL); err != nil {
		t.Fatalf("could not create order_lines table %s", err)
	}
	if _, err := db.Exec(createBatchesOrderLinesTableSQL); err != nil {
		t.Fatalf("could not create batches_order_lines table %s", err)
	}
}

func truncateTables(t *testing.T, db *sql.DB) {
	t.Helper()
	if _, err := db.Exec(truncateTablesSQL); err != nil {
		t.Fatalf("could not clear batches table %s", err)
	}
}

func insertBatch(t *testing.T, db *sql.DB, reference domain.Reference, sku domain.Sku, quantity int, eta time.Time) {
	t.Helper()
	if _, err := db.Exec(insertBatchRow, reference, sku, quantity, eta); err != nil {
		t.Fatalf("could not seed the db with batches: %s", err)
	}
}
func insertOrderLine(t *testing.T, db *sql.DB, orderId domain.Reference, sku domain.Sku, quantity int) {
	t.Helper()
	if _, err := db.Exec(insertOrderLineRow, orderId, sku, quantity); err != nil {
		t.Fatalf("could not seed the db with order lines: %s", err)
	}
}

func insertAllocation(t *testing.T, db *sql.DB, batchRef, orderId domain.Reference) {
	t.Helper()
	if _, err := db.Exec(insertBatchOrderLineRow, batchRef, orderId); err != nil {
		t.Fatalf("could not seed the db with allocations: %s", err)
	}
}

func TestSQLRepository_AddBatch(t *testing.T) {
	db, err := sql.Open("sqlite3", testDBFile)
	assert.Nil(t, err)

	createTables(t, db)
	defer truncateTables(t, db)
	t.Run("can store batch", func(t *testing.T) {
		batch := domain.NewBatch(
			"batch-001",
			"SMALL-TABLE",
			10,
			time.Time{},
		)

		repo := SQLRepository{
			db: &DBWrapper{db},
		}
		err = repo.AddBatch(batch)
		assert.Nil(t, err)

		createdBatch := domain.Batch{}
		row := db.QueryRow(selectBatchRow, batch.Reference)
		err = row.Scan(&createdBatch.Reference, &createdBatch.Sku, &createdBatch.Quantity, &createdBatch.ETA)
		assert.Nil(t, err)

		assert.Equal(t, batch.Reference, createdBatch.Reference)
		assert.Equal(t, batch.Sku, createdBatch.Sku)
		assert.Equal(t, batch.Quantity, createdBatch.Quantity)
		assert.Equal(t, batch.ETA, createdBatch.ETA)
	})
}

func TestSQLRepository_GetBatch(t *testing.T) {
	db, err := sql.Open("sqlite3", testDBFile)
	assert.Nil(t, err)

	createTables(t, db)
	defer truncateTables(t, db)
	t.Run("can retrieve batch", func(t *testing.T) {
		repo := SQLRepository{
			db: &DBWrapper{db},
		}

		existingBatch := domain.Batch{
			Reference: "batch-002",
			Sku:       "LARGE-MIRROR",
			Quantity:  23,
			ETA:       time.Now().AddDate(0, 3, 0).UTC(),
		}
		db.Exec(insertBatchRow, existingBatch.Reference, existingBatch.Sku, existingBatch.Quantity, existingBatch.ETA)

		receivedBatch, err := repo.GetBatch(existingBatch.Reference)

		assert.Nil(t, err)
		assert.Equal(t, existingBatch.Reference, receivedBatch.Reference)
		assert.Equal(t, existingBatch.Sku, receivedBatch.Sku)
		assert.Equal(t, existingBatch.Quantity, receivedBatch.Quantity)
		assert.Equal(t, existingBatch.ETA, receivedBatch.ETA)
	})

	t.Run("can retrieve batch with allocation", func(t *testing.T) {
		db, err := sql.Open("sqlite3", testDBFile)
		assert.Nil(t, err)

		createTables(t, db)
		defer truncateTables(t, db)

		batchRef := domain.Reference("batch-001")
		orderId := domain.Reference("order-012")
		sku := domain.Sku("LARGE-MIRROR")

		insertBatch(t, db, batchRef, sku, 50, time.Time{})
		insertOrderLine(t, db, orderId, sku, 3)
		insertAllocation(t, db, batchRef, orderId)

		repo := SQLRepository{
			db: &DBWrapper{db},
		}

		receivedBatch, err := repo.GetBatch(batchRef)

		assert.Nil(t, err)
		allocations := receivedBatch.Allocations.ToSlice()
		assert.Greater(t, len(allocations), 0)
		assert.Equal(t, allocations[0].OrderID, orderId)
		assert.Equal(t, allocations[0].Sku, sku)
	})
}

func TestSQLRepository_ListBatch(t *testing.T) {
	db, err := sql.Open("sqlite3", testDBFile)
	assert.Nil(t, err)

	createTables(t, db)
	defer truncateTables(t, db)

	sku := domain.Sku("BIG-STICKER")

	insertBatch(t, db, "batch-007", sku, 50, time.Time{})
	insertBatch(t, db, "batch-001", sku, 20, time.Time{})
	insertBatch(t, db, "batch-002", sku, 10, time.Time{})

	repo := SQLRepository{
		db: &DBWrapper{db},
	}

	receivedBatches, err := repo.ListBatches()
	assert.Nil(t, err)

	assert.Len(t, receivedBatches, 3)

	assert.True(t, slices.ContainsFunc(receivedBatches, func(b domain.Batch) bool {
		return b.Reference == domain.Reference("batch-001")
	}))

	assert.True(t, slices.ContainsFunc(receivedBatches, func(b domain.Batch) bool {
		return b.Reference == domain.Reference("batch-002")
	}))

	assert.True(t, slices.ContainsFunc(receivedBatches, func(b domain.Batch) bool {
		return b.Reference == domain.Reference("batch-007")
	}))

}

func TestSQLRepository_AddOrderLine(t *testing.T) {
	db, err := sql.Open("sqlite3", testDBFile)
	assert.Nil(t, err)

	createTables(t, db)

	defer truncateTables(t, db)

	repo := SQLRepository{
		db: &DBWrapper{db},
	}

	orderLine := domain.OrderLine{
		OrderID:  "order-001",
		Sku:      "LARGE-TABLE",
		Quantity: 12,
	}

	err = repo.AddOrderLine(orderLine)
	assert.Nil(t, err)

}

func TestSQLRepository_AllocateToBatch(t *testing.T) {
	db, err := sql.Open("sqlite3", testDBFile)
	assert.Nil(t, err)

	createTables(t, db)

	defer truncateTables(t, db)

	repo := SQLRepository{
		db: &DBWrapper{db},
	}

	sku := "LARGE-TABLE"

	batch := domain.Batch{
		Reference: "batch-012",
		Sku:       domain.Sku(sku),
		Quantity:  40,
	}

	orderLine := domain.OrderLine{
		OrderID:  "order-321",
		Quantity: 12,
		Sku:      domain.Sku(sku),
	}

	err = repo.AddBatch(batch)
	assert.Nil(t, err)

	err = repo.AddOrderLine(orderLine)
	assert.Nil(t, err)

	err = repo.AllocateToBatch(batch, orderLine)
	assert.Nil(t, err)

	allocatedBatch, err := repo.GetBatch(batch.Reference)
	assert.Nil(t, err)

	assert.EqualExportedValues(t, allocatedBatch.Allocations.ToSlice()[0], orderLine)
}

func TestSQLRepository_DeallocateFromBatch(t *testing.T) {
	db, err := sql.Open("sqlite3", testDBFile)
	assert.Nil(t, err)

	createTables(t, db)
	defer truncateTables(t, db)

	repo := SQLRepository{
		db: &DBWrapper{db},
	}

	sku := "LARGE-TABLE"

	batch := domain.Batch{
		Reference: "batch-012",
		Sku:       domain.Sku(sku),
		Quantity:  40,
	}

	orderLine := domain.OrderLine{
		OrderID:  "order-321",
		Quantity: 12,
		Sku:      domain.Sku(sku),
	}

	err = repo.AddBatch(batch)
	assert.Nil(t, err)

	err = repo.AddOrderLine(orderLine)
	assert.Nil(t, err)

	err = repo.AllocateToBatch(batch, orderLine)
	assert.Nil(t, err)

	allocatedBatch, err := repo.GetBatch(batch.Reference)
	assert.Nil(t, err)

	assert.EqualExportedValues(t, allocatedBatch.Allocations.ToSlice()[0], orderLine)

	err = repo.DeallocateFromBatch(allocatedBatch, orderLine)
	assert.Nil(t, err)

	deallocatedBatch, err := repo.GetBatch(batch.Reference)
	assert.Nil(t, err)

	assert.Len(t, deallocatedBatch.Allocations.ToSlice(), 0)
}
