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

func TestAPI(t *testing.T) {

	t.Run("happy path returns 201 and allocated batch", func(t *testing.T) {
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

		service := services.NewStockService(repo)

		server := Server{
			service: &service,
		}

		orderJson := generateOrderLineJson(t, randomOrderId(t, "random"), sku, 10)
		request, _ := http.NewRequest(http.MethodPost, "/allocate", bytes.NewReader(orderJson))
		response := httptest.NewRecorder()
		server.AllocationsHandler(response, request)

		assert.Equal(t, response.Result().StatusCode, http.StatusCreated)

		batchRef := getBatchRef(t, response)

		assert.Equal(t, string(earlyBatchRef), batchRef)
	})

	t.Run("unhappy path returns 400 and error message", func(t *testing.T) {

		unknownSku := randomSku(t, "unknown")
		orderId := randomOrderId(t, "")
		order1 := generateOrderLineJson(t, orderId, unknownSku, 10)

		repo := &repos.FakeRepository{
			Batches: []domain.Batch{
				domain.NewBatch(randomBatchRef(t, ""), randomSku(t, ""), 10, time.Time{}.AddDate(2025, 2, 21)),
				domain.NewBatch(randomBatchRef(t, ""), randomSku(t, ""), 10, time.Time{}.AddDate(2025, 2, 21)),
			},
			BatchAllocations: make(map[domain.Reference]domain.OrderLine),
		}

		service := services.NewStockService(repo)

		server := Server{
			service: &service,
		}

		request, _ := http.NewRequest(http.MethodPost, "/allocate", bytes.NewReader(order1))
		response := httptest.NewRecorder()

		server.AllocationsHandler(response, request)

		assert.Equal(t, response.Result().StatusCode, http.StatusUnprocessableEntity)

		responseRecord := make(map[string]any)

		err := json.Unmarshal(response.Body.Bytes(), &responseRecord)
		assert.Nil(t, err)

		assert.Contains(t, responseRecord["message"], fmt.Sprintf("%s sku is invalid", unknownSku))
	})
}
