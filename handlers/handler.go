package handlers

import "github.com/jackc/pgx/v5/pgxpool"

type Handler struct {
	dbPool *pgxpool.Pool
}

func New(
	dbPool *pgxpool.Pool,
) *Handler {
	return &Handler{
		dbPool: dbPool,
	}
}
