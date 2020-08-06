package srvhttp

import (
	"encoding/json"
	"net/http"

	"github.com/gorilla/mux"

	"github.com/akhripko/dummy/src/models"
	log "github.com/sirupsen/logrus"
)

// /api/hello/{name}
func (s *Server) helloHandler(w http.ResponseWriter, r *http.Request) {
	message, err := s.service.Hello(getName(r))
	if err != nil {
		switch err.(type) {
		case models.ErrNotValidRequest:
			http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
			return
		default:
			http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			log.Error("httpsrv: failed to build hello message:", err)
			return
		}
	}

	err = json.NewEncoder(w).Encode(message)
	if err != nil {
		log.Error("httpsrv: failed to build encode message:", err)
	}
}

func getName(r *http.Request) string {
	return mux.Vars(r)["name"]
}
