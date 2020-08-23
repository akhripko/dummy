package redis

import (
	"context"
	"errors"
	"time"

	"github.com/gomodule/redigo/redis"
	log "github.com/sirupsen/logrus"
)

// ErrNil indicates that a reply value is nil.
var ErrNil = errors.New("(nil)")

// Cache describes connection to Cache server
type Cache struct {
	pool *redis.Pool
}

type Config struct {
	Addr               string
	StatsWriteInterval time.Duration
	IdleTimeout        time.Duration // Close connections after remaining idle for this duration.
	MaxActive          int           // When zero, there is no limit on the number of connections in the pool.
	MaxIdle            int
}

/// New returns the initialized Redis object
func New(ctx context.Context, cfg Config) (*Cache, error) {
	log.Info("Redis init: host=", cfg.Addr)
	c := new(Cache)
	c.initNewPool(cfg)
	if err := c.Check(); err != nil {
		return nil, err
	}
	go func() {
		t := time.NewTicker(cfg.StatsWriteInterval)
		defer t.Stop()
		for {
			select {
			case <-ctx.Done():
				if err := c.pool.Close(); err != nil {
					log.Error("close redis connection error:", err.Error())
					return
				}
				log.Debug("close redis connection")
			case <-t.C:
				currStats := c.pool.Stats()
				log.Infof(
					"REDIS_STATS: active connections: %d, idle connections: %d \n",
					currStats.ActiveCount, currStats.IdleCount,
				)
			}
		}
	}()

	return c, nil
}

func (c *Cache) initNewPool(cfg Config) {
	c.pool = &redis.Pool{
		IdleTimeout: cfg.IdleTimeout,
		MaxActive:   cfg.MaxActive,
		MaxIdle:     cfg.MaxIdle,
		Dial:        func() (redis.Conn, error) { return redis.Dial("tcp", cfg.Addr) },
		Wait:        true,
	}
}

// Check checks if connection exists
func (c *Cache) Check() error {
	r := c.pool.Get()
	defer r.Close()
	_, err := r.Do("PING")
	return err
}
