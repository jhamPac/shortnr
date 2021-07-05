package main

import (
	"log"
	"net/http"

	"github.com/jhampac/shortnr/handler"
	"github.com/jhampac/shortnr/storage/redis"
)

func main() {
	rService, err := redis.New("127.0.0.1", "6379", "password")

	if err != nil {
		panic(err)
	}

	// s, err := r.Save("https://kevia.me", time.Now().Add(time.Hour*1))

	// if err != nil {
	// 	panic(err)
	// }

	// fmt.Println(s)

	r := handler.New("http", "localhost", rService)
	if err := http.ListenAndServe(":8080", r); err != nil {
		log.Fatal(err)
	}
}
