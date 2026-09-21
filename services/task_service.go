package services

import (
	"context"
	"errors"
	"strings"

	"example.com/hello-service/models"
)

var ErrTaskTitleRequired = errors.New("task title is required")

type TaskService struct {
	taskRepository TaskRepository
}

type TaskRepository interface {
	Create(
		ctx context.Context,
		title string,
	) (models.Task, error)
}

func NewTaskService(
	taskRepository TaskRepository,
) *TaskService {
	return &TaskService{
		taskRepository: taskRepository,
	}
}

func (s *TaskService) CreateTask(
	ctx context.Context,
	title string,
) (models.Task, error) {

	normalizedTitle := strings.TrimSpace(title)

	if normalizedTitle == "" {
		return models.Task{}, ErrTaskTitleRequired
	}

	return s.taskRepository.Create(ctx, normalizedTitle)

}
