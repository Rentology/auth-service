package repository

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"github.com/jackc/pgx/v5/pgconn"
	"user-service/internal/domain/models"
	"user-service/internal/repository/postgres"
)

type UserRepository interface {
	CreateUser(ctx context.Context, email string, passHash []byte) (uid int64, err error)
	GetUserByEmail(ctx context.Context, email string) (models.User, error)
	IsAdmin(ctx context.Context, userID int64) (bool, error)
}

type userRepository struct {
	db *sql.DB
}

func NewUserRepository(db *postgres.Database) UserRepository {
	return &userRepository{db: db.Db}
}

func (r *userRepository) CreateUser(ctx context.Context, email string, passHash []byte) (uid int64, err error) {
	const op = "repository.user_repository.CreateUser"
	query := "INSERT INTO users (email, pass_hash) VALUES ($1, $2) RETURNING id"

	// Выполняем запрос напрямую и получаем возвращаемый ID
	var id int64
	err = r.db.QueryRowContext(ctx, query, email, passHash).Scan(&id)
	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) && pgErr.Code == "23505" {
			return 0, fmt.Errorf("%s: %w", op, postgres.ErrUserExists)
		}
		return 0, fmt.Errorf("%s: %w", op, err)
	}

	return id, nil
}

func (r *userRepository) GetUserByEmail(ctx context.Context, email string) (models.User, error) {
	const op = "repository.user_repository.GetUserByEmail"
	query := "SELECT id, email, pass_hash FROM users WHERE email = $1"
	row := r.db.QueryRowContext(ctx, query, email)
	var user models.User
	err := row.Scan(&user.ID, &user.Email, &user.PassHash)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return models.User{}, fmt.Errorf("%s: %w", op, err)
		}
	}
	return user, nil
}

func (r *userRepository) IsAdmin(ctx context.Context, userID int64) (bool, error) {
	return true, nil
}
