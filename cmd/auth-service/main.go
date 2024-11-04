package main

import (
	pb "auth-service/gen/go/auth"
	"auth-service/internal/app"
	"auth-service/internal/config"
	"context"
	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"google.golang.org/grpc"
	"log"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"strconv"
	"syscall"
)

const (
	envLocal = "local"
	envDev   = "dev"
	envProd  = "prod"
)

func main() {
	cfg := config.LoadConfig()
	log := setupLogger(cfg.Env)
	log.Info("starting auth-service", slog.String("env", cfg.Env))
	log.Debug("debug messages are enabled")

	application := app.New(log, cfg.GRPC.Port, cfg.DatabaseUrl, cfg.TokenTTL)

	go application.GRPCSrv.MustRun()

	go func() {
		log.Info("starting REST gateway")
		if err := runRESTGateway(strconv.Itoa(cfg.GRPC.Port)); err != nil {
			log.Error("failed to run REST gateway", slog.String("error", err.Error()))
		}
	}()

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, syscall.SIGTERM, syscall.SIGINT)

	sign := <-stop

	log.Info("stopping application", slog.String("signal", os.Signal.String(sign)))

	application.GRPCSrv.Stop()

	log.Info("application stopped")
}

func setupLogger(env string) *slog.Logger {
	var log *slog.Logger

	switch env {
	case envLocal:
		log = slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}))
	case envDev:
		log = slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}))
	case envProd:
		log = slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelInfo}))
	}

	return log
}

func runRESTGateway(grpcPort string) error {
	ctx := context.Background()
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	mux := runtime.NewServeMux()
	opts := []grpc.DialOption{grpc.WithInsecure()}

	err := pb.RegisterAuthHandlerFromEndpoint(ctx, mux, "localhost:"+grpcPort, opts)
	if err != nil {
		return err
	}
	log.Println("REST gateway is running on port 8081")
	return http.ListenAndServe(":8081", mux)
}
