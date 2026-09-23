package handlers

import (
	"context"

	"example.com/hello-service/models"
)

type TaskRepository interface {
	GetAll(
		ctx context.Context,
	) ([]models.Task, error)

	Create(
		ctx context.Context,
		title string,
	) (models.Task, error)

	Update(
		ctx context.Context,
		id int,
		title string,
		completed bool,
	) (models.Task, error)

	Delete(
		ctx context.Context,
		id int,
	) error
}

type TaskService interface {
	CreateTask(
		ctx context.Context,
		title string,
	) (models.Task, error)
}

type Handler struct {
	taskRepository TaskRepository
	taskService    TaskService
}

func New(
	taskRepository TaskRepository,
	taskService TaskService,
) *Handler {
	return &Handler{
		taskRepository: taskRepository,
		taskService:    taskService,
	}
}
