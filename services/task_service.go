package services

import (
	"context"
	"strings"

	"example.com/hello-service/domain"
	"example.com/hello-service/models"
)

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
		return models.Task{}, domain.ErrTaskTitleRequired
	}

	return s.taskRepository.Create(ctx, normalizedTitle)

}

func (s *TaskService) UpdateTask(
	ctx context.Context,
	id int,
	title string,
	completed bool,
) (models.Task, error) {

	normalizedTitle := strings.TrimSpace(title)

	if normalizedTitle == "" {
		return models.Task{}, domain.ErrTaskTitleRequired
	}

	if id <= 0 {
		return models.Task{}, domain.ErrTaskPositiveInteger
	}

	return s.taskRepository.Update(ctx, id, normalizedTitle, completed)

}

func (s *TaskService) DeleteTask(
	ctx context.Context,
	id int,
) error {

	if id <= 0 {
		return domain.ErrTaskPositiveInteger
	}

	return s.taskRepository.Delete(ctx, id)
}
