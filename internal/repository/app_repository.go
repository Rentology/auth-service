package repository

import (
	"context"
	"database/sql"
	"user-service/internal/domain/models"
	"user-service/internal/repository/postgres"
)

type AppRepository interface {
	App(ctx context.Context, appID int) (models.App, error)
}

type appRepository struct {
	db *sql.DB
}

func NewAppRepository(db *postgres.Database) AppRepository {
	return &appRepository{db: db.Db}
}

func (r *appRepository) App(ctx context.Context, appID int) (models.App, error) {
	return models.App{}, nil
}
