package hellosrv

import (
	"github.com/akhripko/dummy/api"
	"github.com/akhripko/dummy/src/models"
	"github.com/pkg/errors"
)

func (s *Client) Hello(name string) (*models.HelloMessage, error) {
	resp, err := s.client.SayHello(s.ctx, &api.HelloRequest{Name: name})
	if err != nil {
		return nil, errors.Wrap(err, "providers.grpc.hellosrv.client.SayHello")
	}
	return &models.HelloMessage{
		Message: resp.Message,
	}, nil
}
