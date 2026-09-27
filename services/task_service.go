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
	GetAll(
		ctx context.Context,
	) ([]models.Task, error)

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

func (s *TaskService) GetTasks(
	ctx context.Context,
) ([]models.Task, error) {

	return s.taskRepository.GetAll(ctx)
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
