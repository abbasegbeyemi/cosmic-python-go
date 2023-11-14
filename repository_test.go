package cosmicpythongo

import (
	"database/sql"
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

const createOrderLinesTable string = `
	CREATE TABLE IF NOT EXISTS order_lines (
	order_id STRING NOT NULL,
	sku STRING NOT NULL,
	quantity INTEGER NOT NULL
	);
`

const createBatchesOrderLinesTable string = `
    CREATE TABLE IF NOT EXISTS batches_order_lines (
    batch_id STRING NOT NULL,
    order_id STRING NOT NULL,
	FOREIGN KEY(batch_id) REFERENCES batches(reference)
    FOREIGN KEY(order_id) REFERENCES order_lines(order_id)
	PRIMARY KEY(batch_id, order_id)
    );
`

const truncateBatchesTable string = `
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
	if _, err := db.Exec(createOrderLinesTable); err != nil {
		t.Fatalf("could not create order_lines table %s", err)
	}
	if _, err := db.Exec(createBatchesOrderLinesTable); err != nil {
		t.Fatalf("could not create batches_order_lines table %s", err)
	}
}

func truncateTables(t *testing.T, db *sql.DB) {
	t.Helper()
	if _, err := db.Exec(truncateBatchesTable); err != nil {
		t.Fatalf("could not clear batches table %s", err)
	}
}

const defaultSku = Sku("LARGE-TABLE")
const defaultBatchRef = Reference("batch-003")
const defaultOrderID = Reference("order-004")

func insertBatch(t *testing.T, db *sql.DB) {
	t.Helper()
	if _, err := db.Exec("INSERT INTO batches VALUES (?,?,?,?)", defaultBatchRef, defaultSku, 50, time.Now().AddDate(1, 5, 0)); err != nil {
		t.Fatalf("could not seed the db with batches: %s", err)
	}
}
func insertOrderLine(t *testing.T, db *sql.DB) {
	t.Helper()
	if _, err := db.Exec("INSERT INTO order_lines VALUES (?,?,?)", defaultOrderID, defaultSku, 5); err != nil {
		t.Fatalf("could not seed the db with order lines: %s", err)
	}
}

func insertAllocation(t *testing.T, db *sql.DB) {
	t.Helper()
	if _, err := db.Exec("INSERT INTO batches_order_lines VALUES (?,?)", defaultBatchRef, defaultOrderID); err != nil {
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

		insertBatch(t, db)
		insertOrderLine(t, db)
		insertAllocation(t, db)

		repo := SQLRepository{
			db: db,
		}

		receivedBatch, err := repo.GetBatch(defaultBatchRef)

		assert.Nil(t, err)
		allocations := receivedBatch.allocations.ToSlice()
		assert.Greater(t, len(allocations), 0)
		assert.Equal(t, allocations[0].OrderID, defaultOrderID)
		assert.Equal(t, allocations[0].Sku, defaultSku)
	})
}
