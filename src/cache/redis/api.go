package redis

import (
	"github.com/gomodule/redigo/redis"
)

func (c *Cache) Read(name string) (string, error) {
	r := c.pool.Get()
	defer r.Close()

	data, err := redis.String(r.Do("GET", "key_"+name))
	if err == redis.ErrNil {
		return "", nil
	}
	return data, err
}

func (c *Cache) WriteTTL(name, msg string, ttl int) error {
	r := c.pool.Get()
	defer r.Close()

	_, err := r.Do("SETEX", "key_"+name, ttl, msg)
	return err
}
