package storage

import (
	"context"
	"fmt"

	"titan/internal/auth-service/model"

	"github.com/jackc/pgx/v5/pgxpool"
)

type AuthRepository interface {
	SetUser(ctx context.Context, user model.User) error
}
type Database struct {
	db *pgxpool.Pool
}

const createPublicationTableQuery = `
      CREATE TABLE IF NOT EXISTS users (
	id UUID PRIMARY KEY,
    user_name VARCHAR NOT NULL,
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

func (d *Database) SetUser(ctx context.Context, user model.User) error {
	createQuery := `INSERT INTO users (id, user_name, pass_hash, created_at) VALUES ($1, $2, $3, $4)`

	_, err := d.db.Exec(ctx, createQuery, user.ID, user.UserName, user.PasswordHash, user.CreatedAt)
	if err != nil {
		return err
	}

	return nil
}
