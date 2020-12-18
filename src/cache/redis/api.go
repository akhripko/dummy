package redis

import (
	"fmt"
	"time"

	"github.com/gomodule/redigo/redis"
)

func (c *Cache) Read(key string) (string, error) {
	r := c.pool.Get()
	defer r.Close()

	data, err := redis.String(r.Do("GET", key))
	if err == redis.ErrNil {
		return "", nil
	}
	return data, err
}

func (c *Cache) Write(key, value string) error {
	r := c.pool.Get()
	defer r.Close()

	_, err := r.Do("SET", key, value)
	return err
}

func (c *Cache) WriteWithTTL(key, value string, ttl int) error {
	r := c.pool.Get()
	defer r.Close()

	_, err := r.Do("SETEX", key, ttl, value)
	return err
}

func (c *Cache) WriteWithTime(key, value string) error {
	r := c.pool.Get()
	defer r.Close()

	_, err := r.Do("MSET", key, value, key+":t", time.Now().Unix())
	return err
}

func (c *Cache) ReadWithTime(key string) (v string, t int64, err error) {
	r := c.pool.Get()
	defer r.Close()

	data, err := redis.Values(r.Do("MGET", key, key+":t"))
	if err != nil {
		return v, t, err
	}
	// get value
	switch d0 := data[0].(type) {
	case string:
		v = d0
	case []byte:
		v = string(d0)
	default:
		return v, t, fmt.Errorf("unexpected element type, got type %T", v)
	}
	// get time
	if len(data) < 2 {
		return
	}
	t, _ = data[1].(int64)
	return
}

func (c *Cache) EXPIRE(key string, seconds int) error {
	r := c.pool.Get()
	defer r.Close()

	_, err := r.Do("EXPIRE", key, seconds)
	return err
}

func (c *Cache) Exists(key string) (bool, error) {
	r := c.pool.Get()
	defer r.Close()

	return redis.Bool(r.Do("EXISTS", key))
}

func (c *Cache) ReadSet(key string) ([]string, error) {
	r := c.pool.Get()
	defer r.Close()

	data, err := redis.Strings(r.Do("SMEMBERS", key))
	if err != nil {
		if err == redis.ErrNil {
			return nil, nil
		}
		return nil, err
	}
	return data, nil
}

func (c *Cache) AddItemToSet(key string, value string) error {
	r := c.pool.Get()
	defer r.Close()

	_, err := r.Do("SADD", key, value)
	if err != nil {
		return err
	}
	return nil
}

func (c *Cache) IsItemInSet(key string, value string) (bool, error) {
	r := c.pool.Get()
	defer r.Close()

	return redis.Bool(r.Do("SISMEMBER", key, value))
}

func (c *Cache) RemoveItemFromSet(key string, value string) error {
	r := c.pool.Get()
	defer r.Close()

	if _, err := r.Do("SREM", key, value); err != nil {
		return err
	}
	return nil
}
