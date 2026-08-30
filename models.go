package main

import "github.com/jackc/pgx/v5/pgxpool"

type CreateTaskRequest struct {
	Title string `json:"title"`
}

type UpdateTaskRequest struct {
	Title     string `json:"title"`
	Completed bool   `json:"completed"`
}

type DeleteTaskRequest struct {
	ID int `json:"id"`
}

type App struct {
	dbPool *pgxpool.Pool
}
