package service

import (
	"github.com/akhripko/dummy/src/models"
	"github.com/pkg/errors"
	log "github.com/sirupsen/logrus"
)

func (s *Service) Hello(name string) (*models.HelloMessage, error) {
	if len(name) == 0 {
		return nil, models.ErrNotValidRequest("'name' cannot be empty")
	}

	// try read from cache
	msg, err := s.cache.Read(name)
	if err != nil {
		log.Error("cache err:" + err.Error())
	}
	if len(msg) > 0 {
		return &models.HelloMessage{Message: msg}, nil
	}

	// read from storage
	msgModel, err := s.storage.Hello(name)
	if err != nil {
		return nil, errors.Wrap(err, "storage.hello")
	}

	// cache data
	err = s.cache.WriteWithTTL(name, msgModel.Message, 300)
	if err != nil {
		log.Error("cache err:" + err.Error())
	}

	return msgModel, nil
}
