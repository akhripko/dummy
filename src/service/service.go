package service

import (
	"github.com/akhripko/dummy/src/models"
)

type Storage interface {
	Check() error
	Hello(name string) (*models.HelloMessage, error)
}

type Cache interface {
	Check() error
	Read(key string) (string, error)
	WriteWithTTL(key, value string, ttl int) error
}

func New(storage Storage, cache Cache) (*Service, error) {
	// build service
	srv := Service{
		storage:   storage,
		cache:     cache,
		readiness: true,
	}

	return &srv, nil
}

type Service struct {
	storage   Storage
	cache     Cache
	readiness bool
}
