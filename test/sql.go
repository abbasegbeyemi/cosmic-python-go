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
	// testDBFile := "f-test.sqlite"
	db, err := sql.Open("sqlite3", testDBFile)
	assert.Nil(t, err)
	return db
}

const CreateProductsTableSQL string = `
	CREATE TABLE IF NOT EXISTS products (
	id INTEGER PRIMARY KEY,
	sku STRING NOT NULL,
	version_number INTEGER NOT NULL
	);
`

const CreateBatchesTableSQL string = `
	CREATE TABLE IF NOT EXISTS batches (
	id INTEGER PRIMARY KEY,
	reference STRING NOT NULL,
	sku STRING NOT NULL,
	quantity INTEGER NOT NULL,
	eta DATETIME
	);
`

const CreateOrderLinesTableSQL string = `
	CREATE TABLE IF NOT EXISTS order_lines (
	id INTEGER PRIMARY KEY,
	order_id STRING NOT NULL,
	sku STRING NOT NULL,
	quantity INTEGER NOT NULL,
	batch_id INTEGER,
	FOREIGN KEY(batch_id) REFERENCES batches(id)
	);
`

const TruncateTablesSQL string = `
	DELETE FROM products;
	DELETE FROM batches;
	DELETE FROM order_lines;
`

const InsertOrderLineRowSQL string = `INSERT INTO order_lines (order_id, sku, quantity) VALUES (?,?,?)`
const InsertBatchRowSQL string = `INSERT INTO batches (reference, sku, quantity, eta) VALUES(?,?,?,?)`
const InsertBatchOrderLineAllocationSQL string = `INSERT INTO order_lines (order_id, sku, quantity, batch_id) VALUES (?,?,?,?)`

func CreateTables(t *testing.T, db *sql.DB) {
	t.Helper()
	if _, err := db.Exec(CreateProductsTableSQL); err != nil {
		t.Fatalf("could not create products table %s", err)
	}
	if _, err := db.Exec(CreateBatchesTableSQL); err != nil {
		t.Fatalf("could not create batches table %s", err)
	}
	if _, err := db.Exec(CreateOrderLinesTableSQL); err != nil {
		t.Fatalf("could not create order_lines table %s", err)
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
	if _, err := db.Exec(InsertBatchRowSQL, reference, sku, quantity, eta); err != nil {
		t.Fatalf("could not seed the db with batches: %s", err)
	}
}
func InsertOrderLine(t *testing.T, db dbTx, orderId domain.Reference, sku domain.Sku, quantity int) {
	t.Helper()
	if _, err := db.Exec(InsertOrderLineRowSQL, orderId, sku, quantity); err != nil {
		t.Fatalf("could not seed the db with order lines: %s", err)
	}
}

func InsertAllocation(t *testing.T, db dbTx, batchRef, orderId domain.Reference, sku domain.Sku, quantity int) {
	t.Helper()
	var batchId int
	err := db.QueryRow("SELECT id from batches WHERE reference=?", batchRef).Scan(&batchId)
	if err != nil {
		t.Fatalf("could get batch using batchRef: %s: %s", batchRef, err)
	}
	if _, err := db.Exec(InsertBatchOrderLineAllocationSQL, orderId, sku, quantity, batchId); err != nil {
		t.Fatalf("could not seed the db with batches: %s", err)
	}
}
