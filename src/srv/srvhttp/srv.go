package srvhttp

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"sync"
	"time"

	"github.com/gorilla/mux"
	log "github.com/sirupsen/logrus"
)

func New(port int, service Service) (*Server, error) {
	// build http server
	httpSrv := http.Server{
		Addr: fmt.Sprintf(":%d", port),
	}

	// build Server
	var srv Server
	srv.setupHTTP(&httpSrv)
	srv.service = service

	return &srv, nil
}

type Server struct {
	http      *http.Server
	runErr    error
	readiness bool
	service   Service
}

func (s *Server) setupHTTP(srv *http.Server) {
	srv.Handler = s.buildHandler()
	s.http = srv
}

func (s *Server) buildHandler() http.Handler {
	r := mux.NewRouter()
	// path -> handlers

	r.HandleFunc("/api/hello/{name}", s.helloHandler).Methods("GET")
	r.HandleFunc("/api/hello", s.helloHandler).Methods("GET")

	// ==============
	return r
}

func (s *Server) Run(ctx context.Context, wg *sync.WaitGroup) {
	wg.Add(1)
	log.Info("http server: begin run")

	go func() {
		defer wg.Done()
		log.Debug("http server: addr=", s.http.Addr)
		err := s.http.ListenAndServe()
		s.runErr = err
		log.Info("http server: end run > ", err)
	}()

	go func() {
		<-ctx.Done()
		sdCtx, _ := context.WithTimeout(context.Background(), 5*time.Second) // nolint
		msg := s.http.Shutdown(sdCtx)
		log.Info("http server shutdown (", msg, ")")
	}()

	s.readiness = true
}

func (s *Server) HealthCheck() error {
	if !s.readiness {
		return errors.New("http server is not ready yet")
	}
	if s.runErr != nil {
		return errors.New("http server: run issue")
	}
	return nil
}
