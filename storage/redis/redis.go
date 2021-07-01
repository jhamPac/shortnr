package redis

import (
	"fmt"
	"time"

	"github.com/gomodule/redigo/redis"
	"github.com/jhampac/shortnr/storage"
)

type client struct {
	pool *redis.Pool
}

func New(host, port, password string) (storage.Service, error) {
	pool := &redis.Pool{
		MaxIdle:     10,
		IdleTimeout: 240 * time.Second,
		Dial: func() (redis.Conn, error) {
			return redis.Dial("tcp", fmt.Sprintf("%s:%s", host, port))
		},
	}
	return &client{pool}, nil
}

func (c *client) Close() error {
	return c.pool.Close()
}
