package redis

import (
	"context"
	"errors"

	"github.com/gomodule/redigo/redis"
	log "github.com/sirupsen/logrus"
)

// ErrNil indicates that a reply value is nil.
var ErrNil = errors.New("(nil)")

// Cache describes connection to Cache server
type Cache struct {
	pool *redis.Pool
}

// New returns the initialized Cache object
func New(ctx context.Context, redisServer string) (*Cache, error) {
	log.Info("Cache init: host=", redisServer)
	c := new(Cache)
	c.initNewPool(redisServer)
	if err := c.Check(); err != nil {
		return nil, err
	}
	go func() {
		<-ctx.Done()
		err := c.pool.Close()
		if err != nil {
			log.Error("close redis connection error:", err.Error())
			return
		}
		log.Info("close redis connection")
	}()

	return c, nil
}

func (c *Cache) initNewPool(addr string) {
	c.pool = &redis.Pool{
		//MaxIdle:     5,
		//IdleTimeout: 30 * time.Second,
		Dial: func() (redis.Conn, error) { return redis.Dial("tcp", addr) },
		//TestOnBorrow: func(c redis.Conn, t time.Time) error {
		//	_, err := c.Do("PING")
		//	return err
		//},
		//MaxActive: 5,
		//Wait:      true,
	}
}

// Check checks if connection exists
func (c *Cache) Check() error {
	r := c.pool.Get()
	defer r.Close()
	_, err := r.Do("PING")
	return err
}
