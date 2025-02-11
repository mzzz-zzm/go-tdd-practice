package main

import (
	"log"
	"net/http"

	"github.com/mzzz-zzm/go-tdd-practice/adapters/httpserver"
)

func main() {
	// handler := http.HandlerFunc(httpserver.Handler)
	// if err := http.ListenAndServe(":8080", handler); err != nil {
	// 	log.Fatal(err)
	// }

	log.Fatal(http.ListenAndServe(":8080", httpserver.NewHandler()))
}
