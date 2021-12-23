package httpserver

import (
	"encoding/json"
	"github.com/gorilla/mux"
	"net/http"
	_ "net/http/pprof"
)

func (s *Server) GetOrder(w http.ResponseWriter, r *http.Request) {
	orderID, ok := mux.Vars(r)["id"]
	if !ok {
		newErrorResponse(w, http.StatusBadRequest, "empty id order")
	}
	order, err := s.App.GetOrder(r.Context(), orderID)
	if err != nil {
		newErrorResponse(w, http.StatusNotFound, err.Error())
		return
	}

	w.Header().Set("Content-Type", "application/json")
	err = json.NewEncoder(w).Encode(order)
	if err != nil {
		newErrorResponse(w, http.StatusInternalServerError, err.Error())
		return
	}
}
