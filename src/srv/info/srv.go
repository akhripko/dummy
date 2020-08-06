package info

import (
	"context"
	"fmt"
	"io/ioutil"
	"net/http"
	"sync"
	"time"

	log "github.com/sirupsen/logrus"
)

type Server struct {
	http *http.Server
}

func New(port int, healthChecks ...func() error) *Server {
	return &Server{
		http: &http.Server{
			Addr:    fmt.Sprintf(":%d", port),
			Handler: buildHandler(healthChecks),
		},
	}
}

func (s *Server) Run(ctx context.Context, wg *sync.WaitGroup) {
	wg.Add(1)
	log.Info("info server: begin run")

	go func() {
		defer wg.Done()
		log.Debug("info server addr:", s.http.Addr)
		err := s.http.ListenAndServe()
		log.Info("info server: end run > ", err)
	}()

	go func() {
		<-ctx.Done()
		sdCtx, _ := context.WithTimeout(context.Background(), 5*time.Second) // nolint
		msg := s.http.Shutdown(sdCtx)
		log.Info("info server shutdown (", msg, ")")
	}()
}

func buildHandler(healthChecks []func() error) http.Handler {
	handler := http.NewServeMux()
	handler.HandleFunc("/version", serveVersion)
	var checks = func(w http.ResponseWriter, _ *http.Request) { serveCheck(w, healthChecks) }
	handler.HandleFunc("/", checks)
	handler.HandleFunc("/health", checks)
	handler.HandleFunc("/ready", checks)
	return handler
}

func writeFile(file string, response http.ResponseWriter) {
	if data, err := ioutil.ReadFile(file); err == nil { // nolint
		response.WriteHeader(http.StatusOK)
		response.Write(data) // nolint
	} else {
		response.WriteHeader(http.StatusNoContent)
	}
}

func serveCheck(w http.ResponseWriter, checks []func() error) {
	writtenHeader := false
	for _, check := range checks {
		if err := check(); err != nil {
			if !writtenHeader {
				w.WriteHeader(http.StatusInternalServerError)
				writtenHeader = true
			}
			w.Write([]byte(err.Error())) // nolint
			w.Write([]byte("\n\n"))      // nolint
		}
	}

	if !writtenHeader {
		w.WriteHeader(http.StatusNoContent)
	}
}
