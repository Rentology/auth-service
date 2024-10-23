package app

import (
	grpcapp "auth-service/internal/app/grpc"
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

	authService := auth.New(log, userRepository, userRepository, appRepository, tokenTTL)

	grpcApp := grpcapp.New(log, authService, grpcPort)
	return &App{
		GRPCSrv: grpcApp,
	}
}
