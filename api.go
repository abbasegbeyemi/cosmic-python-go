package cosmicpythongo

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/abbasegbeyemi/cosmic-python-go/domain"
)

type service interface {
	Allocate(orderId domain.Reference, sku domain.Sku, quantity int) (domain.Reference, error)
	AddBatch(reference domain.Reference, sku domain.Sku, quantity int, eta time.Time) error
}

type Server struct {
	service service
}

func (s *Server) AllocationsHandler(w http.ResponseWriter, r *http.Request) {
	var orderLine domain.OrderLine

	err := json.NewDecoder(r.Body).Decode(&orderLine)

	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprint(w, `{"message": "could not decode request data"}`)
		return
	}

	batchRef, err := s.service.Allocate(orderLine.OrderID, orderLine.Sku, orderLine.Quantity)

	if err != nil {
		w.WriteHeader(http.StatusUnprocessableEntity)
		fmt.Fprintf(w, `{"message": %q}`, err)
		return
	}

	w.WriteHeader(201)
	fmt.Fprintf(w, `{"batchRef": %q}`, string(batchRef))
}

func (s *Server) StocksHandler(w http.ResponseWriter, r *http.Request) {
	var batch domain.Batch

	err := json.NewDecoder(r.Body).Decode(&batch)

	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprint(w, `{"message": "could not decode request data"}`)
		return
	}

	if err = s.service.AddBatch(batch.Reference, batch.Sku, batch.Quantity, batch.ETA); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprint(w, `{"message": "could not add batch"}`)
		return
	}

	w.WriteHeader(201)
	fmt.Fprintf(w, `{"message": "ok"}`)
}
