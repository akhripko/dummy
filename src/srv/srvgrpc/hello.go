package srvgrpc

import (
	"context"

	"github.com/akhripko/dummy/api"
	"github.com/pkg/errors"
)

// Hello request
func (s *Server) SayHello(_ context.Context, req *api.HelloRequest) (*api.HelloResponse, error) {
	message, err := s.service.Hello(req.Name)
	if err != nil {
		return nil, errors.Wrap(err, "failed to make hello")
	}

	return toHelloResp(message)
}
