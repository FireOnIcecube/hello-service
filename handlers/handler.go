package handlers

import (
	"example.com/hello-service/repositories"
)

type Handler struct {
	taskRepository *repositories.TaskRepository
}

func New(
	taskRepository *repositories.TaskRepository,
) *Handler {
	return &Handler{
		taskRepository: taskRepository,
	}
}
