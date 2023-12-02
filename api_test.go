package cosmicpythongo

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/abbasegbeyemi/cosmic-python-go/domain"
	"github.com/abbasegbeyemi/cosmic-python-go/services"
	"github.com/abbasegbeyemi/cosmic-python-go/test"
	"github.com/abbasegbeyemi/cosmic-python-go/uow"

	"github.com/stretchr/testify/assert"
)

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

		sku := test.RandomSku(t, "")
		otherSku := test.RandomSku(t, "other")
		earlyBatchRef := test.RandomBatchRef(t, "earlyBatchRef")

		postBatchToServer(t, server, earlyBatchRef, sku, 100, time.Time{}.AddDate(2025, 2, 21))
		postBatchToServer(t, server, test.RandomBatchRef(t, "random"), sku, 100, time.Time{}.AddDate(2025, 4, 22))
		postBatchToServer(t, server, test.RandomBatchRef(t, "random"), otherSku, 100, time.Time{}.AddDate(2025, 5, 21))

		orderJson := generateOrderLineJson(t, test.RandomOrderId(t, "random"), sku, 10)
		request, _ := http.NewRequest(http.MethodPost, "/allocate", bytes.NewReader(orderJson))
		response := httptest.NewRecorder()

		server.AllocationsHandler(response, request)
		if assert.Equal(t, response.Result().StatusCode, http.StatusCreated) {
			batchRef := getBatchRef(t, response)
			assert.Equal(t, string(earlyBatchRef), batchRef)
		}
	})

	t.Run("unhappy path returns 400 and error message", func(t *testing.T) {

		unknownSku := test.RandomSku(t, "unknown")
		orderId := test.RandomOrderId(t, "")
		order1 := generateOrderLineJson(t, orderId, unknownSku, 10)

		uow, err := uow.NewSqliteUnitOfWork(db)
		assert.NoError(t, err)

		service := services.StockService{
			UOW: uow,
		}

		server := Server{
			service: &service,
		}

		postBatchToServer(t, server, test.RandomBatchRef(t, ""), test.RandomSku(t, ""), 10, time.Time{}.AddDate(2025, 2, 21))
		postBatchToServer(t, server, test.RandomBatchRef(t, ""), test.RandomSku(t, ""), 10, time.Time{}.AddDate(2025, 2, 21))

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

		batchJson := generateBatchJson(t, "batch-002", test.RandomSku(t, ""), 100, time.Time{})

		request, err := http.NewRequest(http.MethodPost, "/batches", bytes.NewReader([]byte(batchJson)))
		assert.NoError(t, err)

		response := httptest.NewRecorder()

		server.StocksHandler(response, request)
		assert.Equal(t, response.Result().StatusCode, http.StatusCreated)
	})
}
