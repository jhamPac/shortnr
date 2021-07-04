package redis

import (
	"fmt"
	"log"
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

	shortLink := storage.Item{Id: id, URL: url, Expires: expires.Format(time.UnixDate), Visits: 0}

	_, err := conn.Do("HMSET", redis.Args{"Shortener:" + strconv.FormatUint(id, 10)}.AddFlat(shortLink)...)

	if err != nil {
		log.Println(err)
	}

	_, err = conn.Do("EXPIREAT", "Shortener:"+strconv.FormatUint(id, 10), expires.Unix())

	if err != nil {
		return "", err
	}

	return base62.Encode(id), nil
}

func (c *client) Load(code string) (string, error) {
	conn := c.pool.Get()
	defer conn.Close()

	decodedId, err := base62.Decode(code)
	if err != nil {
		return "", err
	}

	urlString, err := redis.String(conn.Do("HGET", "Shortener:"+strconv.FormatUint(decodedId, 10), "url"))
	if err != nil {
		return "", err
	} else if len(urlString) == 0 {
		return "", &storage.LinkError{Msg: "Sorry that link does not"}
	}

	_, err = conn.Do("HINCRBY", "Shortener:"+strconv.FormatUint(decodedId, 10), "visits", 1)

	if err != nil {
		log.Println(err)
	}

	return urlString, nil
}

func (c *client) LoadInfo(code string) (*storage.Item, error) {
	conn := c.pool.Get()
	defer conn.Close()

	decodedId, err := base62.Decode(code)
	if err != nil {
		return nil, err
	}

	values, err := redis.Values(conn.Do("HGETALL", "Shortener:"+strconv.FormatUint(decodedId, 10)))
	if err != nil {
		return nil, err
	} else if len(values) == 0 {
		return nil, &storage.LinkError{Msg: "Sorry no values at that id"}
	}

	var shortLink storage.Item
	err = redis.ScanStruct(values, &shortLink)
	if err != nil {
		return nil, err
	}

	return &shortLink, nil
}

func (c *client) Close() error {
	return c.pool.Close()
}
