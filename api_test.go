package cosmicpythongo

import (
	"bytes"
	"database/sql"
	"encoding/json"
	"fmt"
	"math/rand"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/abbasegbeyemi/cosmic-python-go/domain"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

func insertBatch(t *testing.T, db *sql.DB, reference domain.Reference, sku domain.Sku, quantity int, eta time.Time) {
	t.Helper()
	if _, err := db.Exec("INSERT INTO batches VALUES (?,?,?,?)", reference, sku, quantity, eta); err != nil {
		t.Fatalf("could not seed the db with batches: %s", err)
	}
}

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

func randomSku(t *testing.T, prefix string) domain.Sku {
	t.Helper()
	var sizes = [5]string{"TINY", "SMALL", "MEDIUM", "LARGE", "MASSIVE"}
	var products = [5]string{"TABLE", "CHAIR", "LAMP", "BOTTLE", "KEYRING"}

	genSku := fmt.Sprintf("%s-%s", sizes[rand.Intn(5)], products[rand.Intn(5)])
	if prefix != "" {
		return domain.Sku(fmt.Sprintf("%s-%s", strings.ToUpper(prefix), genSku))
	}
	return domain.Sku(genSku)
}

func randomBatchRef(t *testing.T, suffix string) domain.Reference {
	return domain.Reference(fmt.Sprintf("batch-%s-%s", uuid.New(), suffix))
}

func randomOrderId(t *testing.T, suffix string) domain.Reference {
	return domain.Reference(fmt.Sprintf("order-%s-%s", uuid.New(), suffix))
}

func addStock(t *testing.T, db *sql.DB, batches []domain.Batch) {
	t.Helper()
	for _, batch := range batches {
		insertBatch(t, db, batch.Reference, batch.Sku, batch.Quantity, batch.ETA)
	}

}

func getBatchRef(t *testing.T, response *httptest.ResponseRecorder) string {
	allocatedBatch := make(map[string]interface{})
	err := json.Unmarshal(response.Body.Bytes(), &allocatedBatch)
	assert.Nil(t, err)
	return allocatedBatch["batchRef"].(string)
}

func generateOrderLineJson(t *testing.T, orderId domain.Reference, sku domain.Sku, quantity int) []byte {
	orderLine := domain.OrderLine{
		OrderID:  orderId,
		Sku:      sku,
		Quantity: quantity,
	}

	orderJson, err := json.Marshal(orderLine)
	assert.Nil(t, err)

	return orderJson
}

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

func TestAPI(t *testing.T) {
	db, err := sql.Open("sqlite3", testDBFile)
	assert.Nil(t, err)
	createTables(t, db)

	t.Run("api should return allocation", func(t *testing.T) {
		defer truncateTables(t, db)
		sku := randomSku(t, "")
		otherSku := randomSku(t, "other")

		earlyBatchRef := randomBatchRef(t, "earlyBatchRef")

		addStock(t, db, []domain.Batch{
			{Reference: earlyBatchRef, Sku: sku, Quantity: 100, ETA: time.Time{}.AddDate(2025, 2, 21)},
			{Reference: randomBatchRef(t, "random"), Sku: sku, Quantity: 100, ETA: time.Time{}.AddDate(2025, 2, 22)},
			{Reference: randomBatchRef(t, "random"), Sku: otherSku, Quantity: 100},
		})

		orderJson := generateOrderLineJson(t, randomOrderId(t, "random"), sku, 10)
		request, _ := http.NewRequest(http.MethodPost, "/allocate", bytes.NewReader(orderJson))
		response := httptest.NewRecorder()

		AllocationsServer(response, request)

		assert.Equal(t, response.Result().StatusCode, http.StatusCreated)

		batchRef := getBatchRef(t, response)

		assert.Equal(t, string(earlyBatchRef), batchRef)
	})

	t.Run("allocations are persisted", func(t *testing.T) {
		defer truncateTables(t, db)
		sku := randomSku(t, "")

		batch1 := randomBatchRef(t, "batch1")
		batch2 := randomBatchRef(t, "batch2")

		addStock(t, db, []domain.Batch{
			{Reference: batch2, Sku: sku, Quantity: 10, ETA: time.Time{}.AddDate(2025, 2, 22).UTC()},
			{Reference: batch1, Sku: sku, Quantity: 10, ETA: time.Time{}.AddDate(2025, 2, 21).UTC()},
		})

		order1 := generateOrderLineJson(t, randomOrderId(t, "order1"), sku, 10)

		request, _ := http.NewRequest(http.MethodPost, "/allocate", bytes.NewReader(order1))
		response := httptest.NewRecorder()

		AllocationsServer(response, request)

		assert.Equal(t, response.Result().StatusCode, http.StatusCreated)
		assert.Equal(t, string(batch1), getBatchRef(t, response))

		// Second order should go to batch 2
		order2 := generateOrderLineJson(t, randomOrderId(t, "order2"), sku, 10)
		request, _ = http.NewRequest(http.MethodPost, "/allocate", bytes.NewReader(order2))
		response = httptest.NewRecorder()
		AllocationsServer(response, request)

		assert.Equal(t, response.Result().StatusCode, http.StatusCreated)
		assert.Equal(t, string(batch2), getBatchRef(t, response))
	})
}
