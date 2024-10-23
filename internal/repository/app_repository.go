package repository

import (
	"auth-service/internal/domain/models"
	"auth-service/internal/repository/postgres"
	"context"
	"database/sql"
	"errors"
	"fmt"
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
	const op = "repository.app_repository.App"
	query := "SELECT id, name, secret FROM apps WHERE id = $1"
	row := r.db.QueryRowContext(ctx, query, appID)
	var app models.App
	err := row.Scan(&app.ID, &app.Name, &app.Secret)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return models.App{}, fmt.Errorf("%s: %w", op, postgres.ErrAppNotFound)
		}
		return models.App{}, fmt.Errorf("%s: %w", op, err)
	}
	return models.App{}, nil
}
