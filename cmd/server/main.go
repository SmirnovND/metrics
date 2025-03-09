package main

import (
	"context"
	"fmt"
	"github.com/SmirnovND/metrics/internal/controllers"
	"github.com/SmirnovND/metrics/internal/pkg/system"
	"github.com/SmirnovND/metrics/pb"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"google.golang.org/grpc"
	"net"

	_ "github.com/SmirnovND/metrics/docs"
	"github.com/SmirnovND/metrics/internal/interfaces"
	"github.com/SmirnovND/metrics/internal/middleware"
	"github.com/SmirnovND/metrics/internal/pkg/compressor"
	"github.com/SmirnovND/metrics/internal/pkg/container"
	"github.com/SmirnovND/metrics/internal/pkg/crypto"
	"github.com/SmirnovND/metrics/internal/pkg/loggeer"
	"github.com/SmirnovND/metrics/internal/repo"
	"github.com/SmirnovND/metrics/internal/router"
	usecase "github.com/SmirnovND/metrics/internal/usecase/server"
	"github.com/jmoiron/sqlx"
)

var (
	buildVersion string
	buildDate    string
	buildCommit  string
)

func main() {
	if err := Run(); err != nil {
		panic(err)
	}
}

func Run() error {
	// Вывод информации о сборке
	printBuildInfo()

	diContainer := container.NewContainer(container.WithStartCollectionFunc(usecase.RestoreBackup))

	var cf interfaces.ConfigServerInterface
	var storage *repo.MemStorage
	var db *sqlx.DB
	var grpcService *controllers.ServiceServer
	diContainer.Invoke(func(c interfaces.ConfigServerInterface, s *repo.MemStorage, d *sqlx.DB, grpcServerService *controllers.ServiceServer) {
		cf = c
		storage = s
		db = d
		grpcService = grpcServerService
	})

	stopCh := make(chan struct{})
	defer close(stopCh)

	usecase.TimedBackup(cf, storage, db, stopCh)

	// Запуск HTTP сервера
	httpServer := &http.Server{
		Addr: cf.GetFlagRunAddr(),
		Handler: middleware.ChainMiddleware(
			router.Handler(storage, db, cf.GetFlagRunAddr()),
			loggeer.WithLogging,
			compressor.WithDecompression,
			compressor.WithCompression,
			func(next http.Handler) http.Handler {
				return system.TrustedRangeMiddleware(cf, next)
			},
			func(next http.Handler) http.Handler {
				return crypto.WithCryptoKey(cf, next)
			},
			func(next http.Handler) http.Handler {
				return crypto.WithHashMiddleware(cf, next)
			},
		),
	}

	// Запуск gRPC сервера
	grpcServer := grpc.NewServer()
	pb.RegisterMetricsServiceServer(grpcServer, grpcService)

	go func() {
		listener, err := net.Listen("tcp", cf.GetGRPCAddr())
		if err != nil {
			fmt.Printf("Ошибка при запуске gRPC сервера: %v\n", err)
			return
		}
		fmt.Println("Запуск gRPC сервера на", cf.GetGRPCAddr())
		if err := grpcServer.Serve(listener); err != nil {
			fmt.Printf("Ошибка при работе gRPC сервера: %v\n", err)
		}
	}()

	// Канал для перехвата сигналов
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGTERM, syscall.SIGINT, syscall.SIGQUIT)

	// Горутина для graceful shutdown
	go func() {
		<-sigChan
		fmt.Println("Получен сигнал завершения, сервер корректно завершает работу...")

		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()

		if err := httpServer.Shutdown(ctx); err != nil {
			fmt.Printf("Ошибка при завершении работы HTTP сервера: %v\n", err)
		}

		grpcServer.GracefulStop()

		usecase.Backup(cf, storage, db) // Сохранение несохраненных данных перед выходом
	}()

	// Запуск HTTP сервера и обработка ошибки завершения
	err := httpServer.ListenAndServe()
	if err != nil && err != http.ErrServerClosed {
		return fmt.Errorf("ошибка при запуске HTTP сервера: %w", err)
	}

	fmt.Println("Сервер завершил работу.")
	return nil
}

// printBuildInfo выводит информацию о версии сборки
func printBuildInfo() {
	if buildVersion == "" {
		buildVersion = "N/A"
	}
	if buildDate == "" {
		buildDate = "N/A"
	}
	if buildCommit == "" {
		buildCommit = "N/A"
	}

	fmt.Printf("Build version: %s\n", buildVersion)
	fmt.Printf("Build date: %s\n", buildDate)
	fmt.Printf("Build commit: %s\n", buildCommit)
}
