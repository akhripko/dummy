package srvgql

import (
	"github.com/akhripko/dummy/src/models"
	models2 "github.com/akhripko/dummy/src/srv/srvgql/models"
)

func helloMessageToMessage(message *models.HelloMessage) (*models2.Message, error) {
	return &models2.Message{
		Data: message.Message,
	}, nil

}
