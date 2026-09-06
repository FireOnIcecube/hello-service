package handlers

import (
	"context"

	"example.com/hello-service/models"
)

type FakeTaskRepository struct {
}

func (f *FakeTaskRepository) GetAll(
	ctx context.Context,
) ([]models.Task, error) {
	return []models.Task{
		{
			ID:        100,
			Title:     "Fake Task A",
			Completed: false,
		},
		{
			ID:        200,
			Title:     "Fake Task B",
			Completed: true,
		},
	}, nil
}

func (f *FakeTaskRepository) Create(
	ctx context.Context,
	title string,
) (models.Task, error) {
	return models.Task{}, nil
}

func (f *FakeTaskRepository) Update(
	ctx context.Context,
	id int,
	title string,
	completed bool,
) (models.Task, error) {
	return models.Task{}, nil
}

func (f *FakeTaskRepository) Delete(
	ctx context.Context,
	id int,
) error {
	return nil
}
