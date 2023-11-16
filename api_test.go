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
	"github.com/abbasegbeyemi/cosmic-python-go/repos"
	"github.com/abbasegbeyemi/cosmic-python-go/services"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

func randomSku(t *testing.T, prefix string) domain.Sku {
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

	t.Run("api should return allocation", func(t *testing.T) {
		sku := randomSku(t, "")
		otherSku := randomSku(t, "other")

		earlyBatchRef := randomBatchRef(t, "earlyBatchRef")

		repo := &repos.FakeRepository{
			Batches: []domain.Batch{
				domain.NewBatch(earlyBatchRef, sku, 100, time.Time{}.AddDate(2025, 2, 21)),
				domain.NewBatch(randomBatchRef(t, "random"), sku, 100, time.Time{}.AddDate(2025, 4, 22)),
				domain.NewBatch(randomBatchRef(t, "random"), otherSku, 100, time.Time{}.AddDate(2025, 5, 21)),
			},
			BatchAllocations: make(map[domain.Reference]domain.OrderLine),
		}

		orderJson := generateOrderLineJson(t, randomOrderId(t, "random"), sku, 10)
		request, _ := http.NewRequest(http.MethodPost, "/allocate", bytes.NewReader(orderJson))
		response := httptest.NewRecorder()

		service := services.NewStockService(repo)

		server := Server{
			service: &service,
		}

		server.AllocationsHandler(response, request)

		assert.Equal(t, response.Result().StatusCode, http.StatusCreated)

		batchRef := getBatchRef(t, response)

		assert.Equal(t, string(earlyBatchRef), batchRef)
	})

	t.Run("allocations are persisted", func(t *testing.T) {

		sku := randomSku(t, "")

		batch1 := randomBatchRef(t, "batch1")
		batch2 := randomBatchRef(t, "batch2")
		repo := &repos.FakeRepository{
			Batches: []domain.Batch{
				domain.NewBatch(batch1, sku, 10, time.Time{}.AddDate(2025, 2, 21)),
				domain.NewBatch(batch2, sku, 10, time.Time{}.AddDate(2025, 2, 21)),
			},
			BatchAllocations: make(map[domain.Reference]domain.OrderLine),
		}

		order1 := generateOrderLineJson(t, randomOrderId(t, "order1"), sku, 10)

		service := services.NewStockService(repo)

		server := Server{
			service: &service,
		}

		request, _ := http.NewRequest(http.MethodPost, "/allocate", bytes.NewReader(order1))
		response := httptest.NewRecorder()

		server.AllocationsHandler(response, request)

		assert.Equal(t, response.Result().StatusCode, http.StatusCreated)
		assert.Equal(t, string(batch1), getBatchRef(t, response))

		// Second order should go to batch 2
		order2 := generateOrderLineJson(t, randomOrderId(t, "order2"), sku, 10)
		request, _ = http.NewRequest(http.MethodPost, "/allocate", bytes.NewReader(order2))
		response = httptest.NewRecorder()
		server.AllocationsHandler(response, request)

		assert.Equal(t, response.Result().StatusCode, http.StatusCreated)
		assert.Equal(t, string(batch2), getBatchRef(t, response))
	})
}
