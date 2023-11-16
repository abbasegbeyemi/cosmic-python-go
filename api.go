package cosmicpythongo

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/abbasegbeyemi/cosmic-python-go/domain"
)

type service interface {
	Allocate(orderLine domain.OrderLine) (domain.Reference, error)
}

type Server struct {
	service service
}

func (s *Server) AllocationsHandler(w http.ResponseWriter, r *http.Request) {
	var orderLine domain.OrderLine

	err := json.NewDecoder(r.Body).Decode(&orderLine)

	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprintf(w, `{"message": %q}`, err)
		return
	}

	batchRef, err := s.service.Allocate(orderLine)

	if err != nil {
		w.WriteHeader(http.StatusUnprocessableEntity)
		fmt.Fprintf(w, `{"message": %q}`, err)
		return
	}

	w.WriteHeader(201)
	fmt.Fprintf(w, `{"batchRef": %q}`, string(batchRef))
}
