package main

import (
	"github.com/SmirnovND/metrics/internal/repo"
	"github.com/SmirnovND/metrics/internal/router"
	"net/http"
)

func main() {
	if err := Run(); err != nil {
		panic(err)
	}
}

func Run() error {
	storage := repo.NewMetricRepo()
	return http.ListenAndServe(`:8080`, router.Handler(storage))
}
