package test

import (
	"database/sql"
	"testing"
	"time"

	"github.com/abbasegbeyemi/cosmic-python-go/domain"
	_ "github.com/mattn/go-sqlite3"
	"github.com/stretchr/testify/assert"
)

type dbTx interface {
	Exec(string, ...any) (sql.Result, error)
	Query(string, ...any) (*sql.Rows, error)
	QueryRow(string, ...any) *sql.Row
}

// Will return an sqlite db in a temporary file. Asserts that the db is accessible
func SqliteDB(t *testing.T) *sql.DB {
	t.Helper()
	testDBFile := t.TempDir() + "-test.sqlite"
	db, err := sql.Open("sqlite3", testDBFile)
	assert.Nil(t, err)
	return db
}

const CreateBatchesTableSQL string = `
	CREATE TABLE IF NOT EXISTS batches (
	reference STRING NOT NULL PRIMARY KEY,
	sku STRING NOT NULL,
	quantity INTEGER NOT NULL,
	eta DATETIME
	);
`

const CreateOrderLinesTableSQL string = `
	CREATE TABLE IF NOT EXISTS order_lines (
	order_id STRING NOT NULL,
	sku STRING NOT NULL,
	quantity INTEGER NOT NULL
	);
`

const CreateBatchesOrderLinesTableSQL string = `
    CREATE TABLE IF NOT EXISTS batches_order_lines (
    batch_id STRING NOT NULL,
    order_id STRING NOT NULL,
	FOREIGN KEY(batch_id) REFERENCES batches(reference)
    FOREIGN KEY(order_id) REFERENCES order_lines(order_id)
	PRIMARY KEY(batch_id, order_id)
    );
`

const TruncateTablesSQL string = `
	DELETE FROM batches;
	DELETE FROM order_lines;
	DELETE FROM batches_order_lines;
`

const InsertOrderLineRow string = `INSERT INTO order_lines VALUES (?,?,?)`
const InsertBatchRow string = `INSERT INTO batches VALUES(?,?,?,?)`
const InsertBatchOrderLineRow string = `INSERT INTO batches_order_lines VALUES (?,?)`

func CreateTables(t *testing.T, db *sql.DB) {
	t.Helper()
	if _, err := db.Exec(CreateBatchesTableSQL); err != nil {
		t.Fatalf("could not create batches table %s", err)
	}
	if _, err := db.Exec(CreateOrderLinesTableSQL); err != nil {
		t.Fatalf("could not create order_lines table %s", err)
	}
	if _, err := db.Exec(CreateBatchesOrderLinesTableSQL); err != nil {
		t.Fatalf("could not create batches_order_lines table %s", err)
	}
}

func TruncateTables(t *testing.T, db *sql.DB) {
	t.Helper()
	if _, err := db.Exec(TruncateTablesSQL); err != nil {
		t.Fatalf("could not clear batches table %s", err)
	}
}

func InsertBatch(t *testing.T, db dbTx, reference domain.Reference, sku domain.Sku, quantity int, eta time.Time) {
	t.Helper()
	if _, err := db.Exec(InsertBatchRow, reference, sku, quantity, eta); err != nil {
		t.Fatalf("could not seed the db with batches: %s", err)
	}
}
func InsertOrderLine(t *testing.T, db dbTx, orderId domain.Reference, sku domain.Sku, quantity int) {
	t.Helper()
	if _, err := db.Exec(InsertOrderLineRow, orderId, sku, quantity); err != nil {
		t.Fatalf("could not seed the db with order lines: %s", err)
	}
}

func InsertAllocation(t *testing.T, db dbTx, batchRef, orderId domain.Reference) {
	t.Helper()
	if _, err := db.Exec(InsertBatchOrderLineRow, batchRef, orderId); err != nil {
		t.Fatalf("could not seed the db with allocations: %s", err)
	}
}
