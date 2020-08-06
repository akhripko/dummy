package srvgql

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"sync"
	"time"

	_ "github.com/99designs/gqlgen/cmd" // nolint
	gqlhandler "github.com/99designs/gqlgen/handler"
	log "github.com/sirupsen/logrus"
)

func New(port int, service Service) (*Server, error) {
	// build http server
	httpSrv := http.Server{
		Addr: fmt.Sprintf(":%d", port),
	}

	// build Server
	var srv Server
	srv.service = service
	srv.setupHTTP(&httpSrv)

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
	handler := http.NewServeMux()
	// path -> handlers

	// build gql resolver
	resolver := Resolver{
		service: s.service,
	}

	graphQLHandler := gqlhandler.GraphQL(
		NewExecutableSchema(
			Config{Resolvers: &resolver},
		),
	)

	handler.Handle("/gql/query", graphQLHandler)

	// ==============
	return handler
}

func (s *Server) Run(ctx context.Context, wg *sync.WaitGroup) {
	wg.Add(1)
	log.Info("gql server: begin run")

	go func() {
		defer wg.Done()
		log.Debug("gql server: addr=", s.http.Addr)
		err := s.http.ListenAndServe()
		s.runErr = err
		log.Info("gql server: end run > ", err)
	}()

	go func() {
		<-ctx.Done()
		sdCtx, _ := context.WithTimeout(context.Background(), 5*time.Second) // nolint
		msg := s.http.Shutdown(sdCtx)
		log.Info("gql server shutdown (", msg, ")")
	}()

	s.readiness = true
}

func (s *Server) HealthCheck() error {
	if !s.readiness {
		return errors.New("gql server is not ready yet")
	}
	if s.runErr != nil {
		return errors.New("gql server: run issue")
	}
	return nil
}
