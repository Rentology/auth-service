package app

import (
	grpcapp "auth-service/internal/app/grpc"
	broker "auth-service/internal/broker"
	"auth-service/internal/repository"
	"auth-service/internal/repository/postgres"
	"auth-service/internal/services/auth"
	"log/slog"
	"time"
)

type App struct {
	GRPCSrv *grpcapp.App
}

func New(
	log *slog.Logger,
	grpcPort int,
	databaseUrl string,
	tokenTTL time.Duration,
) *App {
	db, err := postgres.New(databaseUrl)
	if err != nil {
		panic(err)
	}
	userRepository := repository.NewUserRepository(db)
	appRepository := repository.NewAppRepository(db)
	b, err := broker.NewBroker("amqp://guest:guest@localhost:5672/")
	if err != nil {
		panic(err)
	}
	producer, err := broker.NewProducer(b, "test")
	if err != nil {
		panic(err)
	}
	authService := auth.New(log, userRepository, userRepository, appRepository, producer, tokenTTL)

	grpcApp := grpcapp.New(log, authService, grpcPort)
	return &App{
		GRPCSrv: grpcApp,
	}
}
