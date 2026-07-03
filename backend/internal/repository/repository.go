package repository

import (
	"context"

	"github.com/jackc/pgx/v5/pgxpool"
)

type Repository struct {
	db *pgxpool.Pool
}

func New(db *pgxpool.Pool) *Repository {
	return &Repository{db: db}
}

func (r *Repository) GetRootMessage() string {
	return "reservation system backend is running"
}

func (r *Repository) GetHealthStatus() string {
	return "ok"
}

func (r *Repository) CheckHealth(ctx context.Context) error {
	return r.db.Ping(ctx)
}
