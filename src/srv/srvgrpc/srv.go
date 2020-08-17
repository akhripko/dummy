package srvgrpc

import (
	"context"
	"errors"
	"fmt"
	"net"
	"sync"

	"github.com/akhripko/dummy/api"
	grpc_middleware "github.com/grpc-ecosystem/go-grpc-middleware"
	grpc_prometheus "github.com/grpc-ecosystem/go-grpc-prometheus"
	log "github.com/sirupsen/logrus"
	"google.golang.org/grpc"
)

func New(port int, service Service) (*Server, error) {
	// build Service
	srv := Server{
		service: service,
		addr:    fmt.Sprintf(":%d", port),
		server: grpc.NewServer(
			grpc.UnaryInterceptor(grpc_middleware.ChainUnaryServer(grpc_prometheus.UnaryServerInterceptor)),
			grpc.StreamInterceptor(grpc_middleware.ChainStreamServer(grpc_prometheus.StreamServerInterceptor)),
		),
	}
	api.RegisterDummyServiceServer(srv.server, &srv)

	return &srv, nil
}

type Server struct {
	addr      string
	service   Service
	server    *grpc.Server
	runErr    error
	readiness bool
}

func (s *Server) Run(ctx context.Context, wg *sync.WaitGroup) {
	log.Info("grpc server: begin run")

	lis, err := net.Listen("tcp", s.addr)
	if err != nil {
		s.runErr = err
		log.Error("grpc server: run error: ", err)
		return
	}

	wg.Add(1)
	go func() {
		defer wg.Done()
		err := s.server.Serve(lis)
		log.Error("grpc server: end run > ", err)
		s.runErr = err
	}()

	go func() {
		<-ctx.Done()
		s.server.GracefulStop()
		log.Info("grpc server: graceful stop")
	}()

	s.readiness = true
}

func (s *Server) HealthCheck() error {
	if !s.readiness {
		return errors.New("grpcserver is not ready yet")
	}
	if s.runErr != nil {
		return errors.New("grpc service: run issue")
	}
	return nil
}
