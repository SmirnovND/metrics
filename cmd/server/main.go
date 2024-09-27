package main

import (
	"github.com/SmirnovND/metrics/internal/router"
	"net/http"
)

func main() {
	if err := run(); err != nil {
		panic(err)
	}
}

func run() error {
	return http.ListenAndServe(`:8080`, router.Router())
}
