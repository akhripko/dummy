package postgres

import (
	"database/sql"
	"time"

	"github.com/akhripko/dummy/src/models"
)

func (s *Storage) Hello(name string) (*models.HelloMessage, error) {
	//s.mdSrvClient.
	return &models.HelloMessage{
		Message: "Hello, " + name,
	}, nil
}

func (s *Storage) LogEvent(name string) error {
	const q = `insert into activity_log (name, last_touch) values ($1,$2) 
				ON CONFLICT (name) DO UPDATE 
  					SET last_touch = excluded.last_touch`
	_, err := s.db.Exec(q, name, time.Now().UTC())
	if err != nil {
		return err
	}
	return nil
}

func (s *Storage) ReadEventTime(name string) (*time.Time, error) {
	const q = "select last_touch from activity_log where name=$1"
	var lastTouch *time.Time
	err := s.db.QueryRow(q, name).Scan(&lastTouch)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}
	return lastTouch, nil
}
