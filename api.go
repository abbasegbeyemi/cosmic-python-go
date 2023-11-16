package cosmicpythongo

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/abbasegbeyemi/cosmic-python-go/domain"
	"github.com/abbasegbeyemi/cosmic-python-go/repos"
)

func AllocationsServer(w http.ResponseWriter, r *http.Request) {
	repo, err := repos.NewSqliteRepository("orders.sqlite")
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprintf(w, `{"message": %s}`, err)
	}

	batches, err := repo.ListBatches()

	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprintf(w, `{"message": %s}`, err)
		return
	}

	var orderLine domain.OrderLine

	err = json.NewDecoder(r.Body).Decode(&orderLine)

	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprintf(w, `{"message": %s}`, err)
		return
	}

	batchRef, err := domain.Allocate(orderLine, batches)

	if err != nil {
		w.WriteHeader(http.StatusUnprocessableEntity)
		fmt.Fprintf(w, `{"message": %s}`, err)
		return
	}

	if err = repo.AddOrderLine(orderLine); err != nil {
		w.WriteHeader(http.StatusUnprocessableEntity)
		fmt.Fprintf(w, `{"message": %s}`, err)
		return
	}

	batchToAllocate, err := repo.GetBatch(batchRef)

	if err != nil {
		w.WriteHeader(http.StatusUnprocessableEntity)
		fmt.Fprintf(w, `{"message": %s}`, err)
	}

	if err := repo.AllocateToBatch(batchToAllocate, orderLine); err != nil {
		w.WriteHeader(http.StatusUnprocessableEntity)
		fmt.Fprintf(w, `{"message": %s}`, err)
	}

	w.WriteHeader(201)
	fmt.Fprintf(w, `{"batchRef": %q}`, string(batchRef))
}
