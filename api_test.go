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

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

func randomSku(t *testing.T, prefix string) Sku {
	t.Helper()
	var sizes = [5]string{"TINY", "SMALL", "MEDIUM", "LARGE", "MASSIVE"}
	var products = [5]string{"TABLE", "CHAIR", "LAMP", "BOTTLE", "KEYRING"}

	genSku := fmt.Sprintf("%s-%s", sizes[rand.Intn(5)], products[rand.Intn(5)])
	if prefix != "" {
		return Sku(fmt.Sprintf("%s-%s", strings.ToUpper(prefix), genSku))
	}
	return Sku(genSku)
}

func randomBatchRef(t *testing.T, suffix string) Reference {
	return Reference(fmt.Sprintf("batch-%s-%s", uuid.New(), suffix))
}

func randomOrderId(t *testing.T, suffix string) Reference {
	return Reference(fmt.Sprintf("order-%s-%s", uuid.New(), suffix))
}

func addStock(t *testing.T, db *sql.DB, batches []Batch) {
	t.Helper()
	for _, batch := range batches {
		insertBatch(t, db, batch.reference, batch.sku, batch.quantity, batch.eta)
	}

}

func getBatchRef(t *testing.T, response *httptest.ResponseRecorder) string {
	allocatedBatch := make(map[string]interface{})
	err := json.Unmarshal(response.Body.Bytes(), &allocatedBatch)
	assert.Nil(t, err)
	return allocatedBatch["batchRef"].(string)
}

func generateOrderLineJson(t *testing.T, orderId Reference, sku Sku, quantity int) []byte {
	orderLine := OrderLine{
		OrderID:  orderId,
		Sku:      sku,
		Quantity: quantity,
	}

	orderJson, err := json.Marshal(orderLine)
	assert.Nil(t, err)

	return orderJson
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

		addStock(t, db, []Batch{
			{reference: earlyBatchRef, sku: sku, quantity: 100, eta: time.Time{}.AddDate(2025, 2, 21)},
			{reference: randomBatchRef(t, "random"), sku: sku, quantity: 100, eta: time.Time{}.AddDate(2025, 2, 22)},
			{reference: randomBatchRef(t, "random"), sku: otherSku, quantity: 100},
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

		addStock(t, db, []Batch{
			{reference: batch1, sku: sku, quantity: 10, eta: time.Time{}.AddDate(2025, 2, 21)},
			{reference: batch2, sku: sku, quantity: 10, eta: time.Time{}.AddDate(2025, 2, 22)},
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
