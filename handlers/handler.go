package handlers

import (
	"context"

	"example.com/hello-service/models"
)

type TaskService interface {
	CreateTask(
		ctx context.Context,
		title string,
	) (models.Task, error)

	GetTasks(ctx context.Context,
	) ([]models.Task, error)

	UpdateTask(ctx context.Context,
		id int,
		title string,
		completed bool) (models.Task, error)

	DeleteTask(ctx context.Context,
		id int) error
}

type Handler struct {
	taskService TaskService
}

func New(
	taskService TaskService,
) *Handler {
	return &Handler{
		taskService: taskService,
	}
}
