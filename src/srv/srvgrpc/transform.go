package srvgrpc

import (
	"github.com/akhripko/dummy/api"
	"github.com/akhripko/dummy/src/models"
)

func toHelloResp(message *models.HelloMessage) (*api.HelloResponse, error) {
	return &api.HelloResponse{
		Message: message.Message,
	}, nil

}
