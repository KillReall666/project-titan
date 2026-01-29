package storage

import (
	"context"
	"fmt"

	"titan/internal/auth-service/model"

	"github.com/jackc/pgx/v5/pgxpool"
)

type AuthRepository interface {
	SetUser(ctx context.Context, user model.User) error
	GetUser(ctx context.Context, email string) (model.User, error)
}
type Database struct {
	db *pgxpool.Pool
}

const createPublicationTableQuery = `
      CREATE TABLE IF NOT EXISTS users (
	id UUID PRIMARY KEY,
    email VARCHAR NOT NULL,
    pass_hash VARCHAR NOT NULL,
    created_at TIMESTAMP DEFAULT NOW()
);`

func New(ctx context.Context, connString string) (*Database, error) {
	conn, err := pgxpool.New(ctx, connString)
	if err != nil {
		return nil, fmt.Errorf("unable to connect to database: %w", err)
	}

	_, err = conn.Exec(ctx, createPublicationTableQuery)
	if err != nil {
		return nil, fmt.Errorf("unable to connect to database: %w", err)
	}

	return &Database{db: conn}, nil
}

// SetUser Пока один метод - можно оставить так, но как только начнёт разрастаться, разбить на репозитории.
func (d *Database) SetUser(ctx context.Context, user model.User) error {
	createQuery := `INSERT INTO users (id, email, pass_hash, created_at) VALUES ($1, $2, $3, $4)`

	_, err := d.db.Exec(ctx, createQuery, user.ID, user.Email, user.PasswordHash, user.CreatedAt)
	if err != nil {
		return err
	}

	return nil
}

func (d *Database) GetUser(ctx context.Context, email string) (model.User, error) {
	var user model.User

	getQuery := `SELECT id, email, pass_hash FROM users WHERE email = $1`

	err := d.db.QueryRow(ctx, getQuery, email).Scan(&user.ID, &user.Email, &user.PasswordHash)
	if err != nil {
		return model.User{}, fmt.Errorf("err when get user: %v", err)
	}

	return user, nil
}
