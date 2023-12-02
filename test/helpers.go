package test

import (
	"fmt"
	"math/rand"
	"strings"
	"testing"

	"github.com/abbasegbeyemi/cosmic-python-go/domain"
	"github.com/google/uuid"
)

func RandomSku(t *testing.T, prefix string) domain.Sku {
	t.Helper()
	var sizes = [5]string{"TINY", "SMALL", "MEDIUM", "LARGE", "MASSIVE"}
	var products = [5]string{"TABLE", "CHAIR", "LAMP", "BOTTLE", "KEYRING"}

	genSku := fmt.Sprintf("%s-%s", sizes[rand.Intn(5)], products[rand.Intn(5)])
	if prefix != "" {
		return domain.Sku(fmt.Sprintf("%s-%s", strings.ToUpper(prefix), genSku))
	}
	return domain.Sku(genSku)
}

func RandomBatchRef(t *testing.T, suffix string) domain.Reference {
	t.Helper()
	return domain.Reference(fmt.Sprintf("batch-%s-%s", uuid.New(), suffix))
}

func RandomOrderId(t *testing.T, suffix string) domain.Reference {
	t.Helper()
	return domain.Reference(fmt.Sprintf("order-%s-%s", uuid.New(), suffix))
}
