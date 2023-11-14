package cosmicpythongo

import (
	"database/sql"
	"testing"

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

const truncateBatchesTable string = `
	DELETE FROM batches;
`

func createTables(t *testing.T, db *sql.DB) {
	t.Helper()
	if _, err := db.Exec(createBatchesTable); err != nil {
		t.Fatalf("could not create batches table %s", err)
	}
}

func truncateTables(t *testing.T, db *sql.DB) {
	t.Helper()
	if _, err := db.Exec(truncateBatchesTable); err != nil {
		t.Fatalf("could not clear batches table %s", err)
	}
}

func TestSQLRepository_Add(t *testing.T) {
	db, err := sql.Open("sqlite3", dbFile)
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
		err = repo.Add(batch)
		assert.Nil(t, err)

		createdBatch := Batch{}
		row := db.QueryRow(`SELECT reference, sku, quantity, eta FROM "batches" WHERE reference=?`, batch.reference)
		err = row.Scan(&createdBatch.reference, &createdBatch.sku, &createdBatch.quantity, &createdBatch.ETA)
		assert.Nil(t, err)

		assert.EqualExportedValues(t, batch, createdBatch)

		// session commit
		// execute query in pure sql to get batches
		// assert batches received as expected
	})

	t.Run("can retrieve batch", func(t *testing.T) {

	})
}
