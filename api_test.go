package cosmicpythongo

import (
	"bytes"
	"encoding/json"
	"fmt"
	"math/rand"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/abbasegbeyemi/cosmic-python-go/domain"
	"github.com/abbasegbeyemi/cosmic-python-go/services"
	"github.com/abbasegbeyemi/cosmic-python-go/test"
	"github.com/abbasegbeyemi/cosmic-python-go/uow"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

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
	t.Helper()
	return domain.Reference(fmt.Sprintf("batch-%s-%s", uuid.New(), suffix))
}

func randomOrderId(t *testing.T, suffix string) domain.Reference {
	t.Helper()
	return domain.Reference(fmt.Sprintf("order-%s-%s", uuid.New(), suffix))
}

func getBatchRef(t *testing.T, response *httptest.ResponseRecorder) string {
	t.Helper()
	allocatedBatch := make(map[string]interface{})
	assert.NoError(t, json.Unmarshal(response.Body.Bytes(), &allocatedBatch))
	return allocatedBatch["batchRef"].(string)
}

func generateOrderLineJson(t *testing.T, orderId domain.Reference, sku domain.Sku, quantity int) []byte {
	t.Helper()
	orderLine := domain.OrderLine{
		OrderID:  orderId,
		Sku:      sku,
		Quantity: quantity,
	}

	orderJson, err := json.Marshal(orderLine)
	assert.NoError(t, err)

	return orderJson
}

func generateBatchJson(t *testing.T, batchRef domain.Reference, sku domain.Sku, quantity int, eta time.Time) []byte {
	t.Helper()
	batch := domain.Batch{
		Reference: batchRef,
		Sku:       sku,
		Quantity:  quantity,
		ETA:       eta,
	}
	batchJson, err := json.Marshal(batch)
	assert.Nil(t, err)

	return batchJson
}

func postBatchToServer(t *testing.T, server Server, batchRef domain.Reference, sku domain.Sku, quantity int, eta time.Time) {
	t.Helper()
	batchJson := generateBatchJson(t, batchRef, sku, quantity, eta)

	request, err := http.NewRequest(http.MethodPost, "/batches", bytes.NewReader([]byte(batchJson)))
	assert.NoError(t, err)

	response := httptest.NewRecorder()

	server.StocksHandler(response, request)
	assert.Equal(t, response.Result().StatusCode, http.StatusCreated)
}

func TestAPI_AllocationsHandler(t *testing.T) {
	db := test.SqliteDB(t)
	test.CreateTables(t, db)
	defer test.TruncateTables(t, db)
	t.Run("happy path returns 201 and allocated batch", func(t *testing.T) {
		uow, err := uow.NewSqliteUnitOfWork(db)
		assert.NoError(t, err)

		service := services.StockService{
			UOW: uow,
		}

		server := Server{
			service: &service,
		}

		sku := randomSku(t, "")
		otherSku := randomSku(t, "other")
		earlyBatchRef := randomBatchRef(t, "earlyBatchRef")

		postBatchToServer(t, server, earlyBatchRef, sku, 100, time.Time{}.AddDate(2025, 2, 21))
		postBatchToServer(t, server, randomBatchRef(t, "random"), sku, 100, time.Time{}.AddDate(2025, 4, 22))
		postBatchToServer(t, server, randomBatchRef(t, "random"), otherSku, 100, time.Time{}.AddDate(2025, 5, 21))

		orderJson := generateOrderLineJson(t, randomOrderId(t, "random"), sku, 10)
		request, _ := http.NewRequest(http.MethodPost, "/allocate", bytes.NewReader(orderJson))
		response := httptest.NewRecorder()

		server.AllocationsHandler(response, request)
		if assert.Equal(t, response.Result().StatusCode, http.StatusCreated) {
			batchRef := getBatchRef(t, response)
			assert.Equal(t, string(earlyBatchRef), batchRef)
		}
	})

	t.Run("unhappy path returns 400 and error message", func(t *testing.T) {

		unknownSku := randomSku(t, "unknown")
		orderId := randomOrderId(t, "")
		order1 := generateOrderLineJson(t, orderId, unknownSku, 10)

		uow, err := uow.NewSqliteUnitOfWork(db)
		assert.NoError(t, err)

		service := services.StockService{
			UOW: uow,
		}

		server := Server{
			service: &service,
		}

		postBatchToServer(t, server, randomBatchRef(t, ""), randomSku(t, ""), 10, time.Time{}.AddDate(2025, 2, 21))
		postBatchToServer(t, server, randomBatchRef(t, ""), randomSku(t, ""), 10, time.Time{}.AddDate(2025, 2, 21))

		request, _ := http.NewRequest(http.MethodPost, "/allocate", bytes.NewReader(order1))
		response := httptest.NewRecorder()

		server.AllocationsHandler(response, request)

		assert.Equal(t, response.Result().StatusCode, http.StatusUnprocessableEntity)

		responseRecord := make(map[string]any)

		assert.NoError(t, json.Unmarshal(response.Body.Bytes(), &responseRecord))
		assert.Contains(t, responseRecord["message"], fmt.Sprintf("%s sku is invalid", unknownSku))
	})
}

func TestAPI_BatchesHandler(t *testing.T) {
	db := test.SqliteDB(t)
	test.CreateTables(t, db)
	defer test.TruncateTables(t, db)
	t.Run("will add batch", func(t *testing.T) {
		uow, err := uow.NewSqliteUnitOfWork(db)
		assert.NoError(t, err)

		service := services.StockService{
			UOW: uow,
		}

		server := Server{
			service: &service,
		}

		batchJson := generateBatchJson(t, "batch-002", randomSku(t, ""), 100, time.Time{})

		request, err := http.NewRequest(http.MethodPost, "/batches", bytes.NewReader([]byte(batchJson)))
		assert.NoError(t, err)

		response := httptest.NewRecorder()

		server.StocksHandler(response, request)
		assert.Equal(t, response.Result().StatusCode, http.StatusCreated)
	})
}
