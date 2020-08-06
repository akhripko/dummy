package prometheus

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"sync"
	"time"

	"github.com/prometheus/client_golang/prometheus/promhttp"
	log "github.com/sirupsen/logrus"
)

type Server struct {
	http      *http.Server
	runErr    error
	readiness bool
}

func New(port int) *Server {
	return &Server{
		http: &http.Server{
			Addr:    fmt.Sprintf(":%d", port),
			Handler: handler(),
		},
	}
}

func (s *Server) Run(ctx context.Context, wg *sync.WaitGroup) {
	wg.Add(1)
	log.Info("prometheus server: begin run")

	go func() {
		defer wg.Done()
		log.Debug("prometheus server addr:", s.http.Addr)
		err := s.http.ListenAndServe()
		s.runErr = err
		log.Info("prometheus server: end run >", err)
	}()

	go func() {
		<-ctx.Done()
		sdCtx, _ := context.WithTimeout(context.Background(), 5*time.Second) // nolint
		msg := s.http.Shutdown(sdCtx)
		log.Info("prometheus server shutdown (", msg, ")")
	}()

	s.readiness = true
}

func handler() http.Handler {
	handler := http.NewServeMux()
	handler.Handle("/metrics", promhttp.Handler())
	return handler
}

func (s *Server) HealthCheck() error {
	if !s.readiness {
		return errors.New("prometheus server is not ready yet")
	}
	if s.runErr != nil {
		return errors.New("prometheus server: run issue")
	}
	return nil
}
