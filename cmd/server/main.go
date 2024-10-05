package main

import (
	"flag"
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

	var flagRunAddr string
	flag.StringVar(&flagRunAddr, "a", "localhost:8080", "address and port to run server")
	flag.Parse()
	return http.ListenAndServe(flagRunAddr, router.Handler(storage))
}
