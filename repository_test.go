package cosmicpythongo

import (
	"database/sql"
	"slices"
	"testing"
	"time"

	_ "github.com/mattn/go-sqlite3"
	"github.com/stretchr/testify/assert"
)

const createBatchesTable string = `
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
	if _, err := db.Exec(createBatchesTable); err != nil {
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

func insertBatch(t *testing.T, db *sql.DB, reference Reference, sku Sku, quantity int) {
	t.Helper()
	if _, err := db.Exec("INSERT INTO batches VALUES (?,?,?,?)", reference, sku, quantity, time.Now().AddDate(1, 5, 0)); err != nil {
		t.Fatalf("could not seed the db with batches: %s", err)
	}
}
func insertOrderLine(t *testing.T, db *sql.DB, orderId Reference, sku Sku, quantity int) {
	t.Helper()
	if _, err := db.Exec("INSERT INTO order_lines VALUES (?,?,?)", orderId, sku, quantity); err != nil {
		t.Fatalf("could not seed the db with order lines: %s", err)
	}
}

func insertAllocation(t *testing.T, db *sql.DB, batchRef, orderId Reference) {
	t.Helper()
	if _, err := db.Exec("INSERT INTO batches_order_lines VALUES (?,?)", batchRef, orderId); err != nil {
		t.Fatalf("could not seed the db with allocations: %s", err)
	}
}

func TestSQLRepository_AddBatch(t *testing.T) {
	db, err := sql.Open("sqlite3", testDBFile)
	assert.Nil(t, err)

	createTables(t, db)
	defer truncateTables(t, db)
	t.Run("can store batch", func(t *testing.T) {
		batch := NewBatch(
			"batch-001",
			"SMALL-TABLE",
			10,
		)

		repo := SQLRepository{
			db: db,
		}
		err = repo.AddBatch(batch)
		assert.Nil(t, err)

		createdBatch := Batch{}
		row := db.QueryRow(`SELECT reference, sku, quantity, eta FROM "batches" WHERE reference=?`, batch.reference)
		err = row.Scan(&createdBatch.reference, &createdBatch.sku, &createdBatch.quantity, &createdBatch.ETA)
		assert.Nil(t, err)

		assert.EqualExportedValues(t, batch, createdBatch)
	})
}

func TestSQLRepository_GetBatch(t *testing.T) {
	db, err := sql.Open("sqlite3", testDBFile)
	assert.Nil(t, err)

	createTables(t, db)
	defer truncateTables(t, db)
	t.Run("can retrieve batch", func(t *testing.T) {
		repo := SQLRepository{
			db: db,
		}

		existingBatch := Batch{
			reference: "batch-002",
			sku:       "LARGE-MIRROR",
			quantity:  23,
			ETA:       time.Now().AddDate(0, 3, 0),
		}
		db.Exec(insertBatchRow, existingBatch.reference, existingBatch.sku, existingBatch.quantity, existingBatch.ETA)

		receivedBatch, err := repo.GetBatch(existingBatch.reference)

		assert.Nil(t, err)
		assert.EqualExportedValues(t, existingBatch, receivedBatch)
	})

	t.Run("can retrieve batch with allocation", func(t *testing.T) {
		db, err := sql.Open("sqlite3", testDBFile)
		assert.Nil(t, err)

		createTables(t, db)
		defer truncateTables(t, db)

		batchRef := Reference("batch-001")
		orderId := Reference("order-012")
		sku := Sku("LARGE-MIRROR")

		insertBatch(t, db, batchRef, sku, 50)
		insertOrderLine(t, db, orderId, sku, 3)
		insertAllocation(t, db, batchRef, orderId)

		repo := SQLRepository{
			db: db,
		}

		receivedBatch, err := repo.GetBatch(batchRef)

		assert.Nil(t, err)
		allocations := receivedBatch.allocations.ToSlice()
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

	sku := Sku("BIG-STICKER")

	insertBatch(t, db, "batch-007", sku, 50)
	insertBatch(t, db, "batch-001", sku, 20)
	insertBatch(t, db, "batch-002", sku, 10)

	repo := SQLRepository{
		db: db,
	}

	receivedBatches, err := repo.ListBatches()
	assert.Nil(t, err)

	assert.Len(t, receivedBatches, 3)

	assert.True(t, slices.ContainsFunc(receivedBatches, func(b Batch) bool {
		return b.reference == Reference("batch-001")
	}))

	assert.True(t, slices.ContainsFunc(receivedBatches, func(b Batch) bool {
		return b.reference == Reference("batch-002")
	}))

	assert.True(t, slices.ContainsFunc(receivedBatches, func(b Batch) bool {
		return b.reference == Reference("batch-007")
	}))

}
