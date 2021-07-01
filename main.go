package main

import (
	"github.com/jhampac/shortnr/storage/redis"
)

func main() {
	_, err := redis.New("127.0.0.1", "6379", "password")

	if err != nil {
		panic(err)
	}
}
