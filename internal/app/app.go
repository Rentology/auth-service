package app

import (
	"log/slog"
	"time"
	grpcapp "user-service/internal/app/grpc"
	"user-service/internal/repository"
	"user-service/internal/repository/postgres"
	"user-service/internal/services/auth"
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
