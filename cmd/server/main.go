package main

import (
	"github.com/SmirnovND/metrics/internal/router"
	"net/http"
)

func main() {
	if err := Run(); err != nil {
		panic(err)
	}
}

func Run() error {
	return http.ListenAndServe(`:8080`, router.Router())
}
