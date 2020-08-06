package srvgrpc

import "github.com/akhripko/dummy/src/models"

type Service interface {
	Hello(name string) (*models.HelloMessage, error)
}
