package cosmicpythongo

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"
)

func AllocationsServer(w http.ResponseWriter, r *http.Request) {
	db, err := sql.Open("sqlite3", "orders_test.sqlite")
	if err != nil {
		w.WriteHeader(500)
	}
	repo := SQLRepository{db: db}
	batches, err := repo.ListBatches()

	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprintf(w, `{"message": %s}`, err)
		return
	}

	var orderLine OrderLine

	err = json.NewDecoder(r.Body).Decode(&orderLine)

	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprintf(w, `{"message": %s}`, err)
		return
	}

	batchRef, err := Allocate(orderLine, batches)

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
