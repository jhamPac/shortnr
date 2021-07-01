package redis

import (
	"fmt"
	"math/rand"
	"strconv"
	"time"

	"github.com/gomodule/redigo/redis"
	"github.com/jhampac/shortnr/base62"
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

func (c *client) isUsed(id uint64) bool {
	conn := c.pool.Get()
	defer conn.Close()

	exists, err := redis.Bool(conn.Do("EXISTS", "Shortener:"+strconv.FormatUint(id, 10)))
	if err != nil {
		return false
	}
	return exists
}

func (c *client) Save(url string, expires time.Time) (string, error) {
	conn := c.pool.Get()
	defer conn.Close()

	var id uint64

	// keeps looping until an id thats not used is found
	for used := true; used; used = c.isUsed(id) {
		id = rand.Uint64()
	}

	shortLink := storage.Item{id, url, expires.Format(time.UnixDate), 0}

	_, err := conn.Do("HMSET", redis.Args{"Shortener:" + strconv.FormatUint(id, 10)}.AddFlat(shortLink)...)

	_, err = conn.Do("EXPIREAT", "Shortener:"+strconv.FormatUint(id, 10), expires.Unix())

	if err != nil {
		return "", err
	}

	return base62.Encode(id), nil
}

func (c *client) Close() error {
	return c.pool.Close()
}
