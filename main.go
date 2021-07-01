package main

import (
	"fmt"
	"time"

	"github.com/jhampac/shortnr/storage/redis"
)

func main() {
	r, err := redis.New("127.0.0.1", "6379", "password")

	if err != nil {
		panic(err)
	}

	s, err := r.Save("https://kevia.me", time.Now().Add(time.Hour*1))

	if err != nil {
		panic(err)
	}

	fmt.Println(s)
}
