package services

import (
	"context"
	"errors"
	"testing"

	"example.com/hello-service/models"
)

type fakeTaskRepository struct {
	createCalled bool
	createTitle  string
	createErr    error
}

func (f *fakeTaskRepository) Create(ctx context.Context,
	title string) (models.Task, error) {
	f.createCalled = true
	f.createTitle = title

	if f.createErr != nil {
		return models.Task{}, f.createErr
	}

	return models.Task{
		ID:        2,
		Title:     title,
		Completed: false,
	}, nil
}

// Create 正常流程
func TestTaskServiceCreateTask(t *testing.T) {

	// Arrange
	fakeRepo := fakeTaskRepository{}
	service := NewTaskService(&fakeRepo)

	ctx := context.Background()

	// Act
	task, err := service.CreateTask(ctx, "  Learn Go  ")

	// Assert
	if err != nil {
		t.Fatalf(
			"Create Task 執行失敗: %v",
			err,
		)
	}

	if !fakeRepo.createCalled {
		t.Fatal(
			"預期 create task 應該被呼叫",
		)
	}

	if task.Title != fakeRepo.createTitle {
		t.Errorf(
			"預期創建 task title: %s , 實際 title: %s",
			fakeRepo.createTitle,
			task.Title,
		)
	}
}

// Create title 為空
func TestTaskServiceCreateTaskEmptyTitle(t *testing.T) {

	// Arrange
	fakeRepo := fakeTaskRepository{}
	service := NewTaskService(&fakeRepo)

	ctx := context.Background()

	// Act
	_, err := service.CreateTask(ctx, "    ")

	// Assert
	if err == nil {
		t.Fatal(
			"預期回報錯誤",
		)
	}

	if !errors.Is(err, ErrTaskTitleRequired) {
		t.Fatalf(
			"預期回報錯誤為: %v , 實際回報錯誤: %v",
			ErrTaskTitleRequired,
			err,
		)
	}

	if fakeRepo.createCalled {
		t.Fatal(
			"create 不該被呼叫",
		)
	}

}

// Create internal service error
func TestTaskServiceCreateTaskRepositoryError(t *testing.T) {
	// Arrange
	fakeRepo := fakeTaskRepository{
		createErr: errors.New("database failed"),
	}
	service := NewTaskService(&fakeRepo)

	ctx := context.Background()

	// Act
	_, err := service.CreateTask(ctx, " internal error ")

	// Assert
	if err == nil {
		t.Fatal("預期回報錯誤")
	}

	if !fakeRepo.createCalled {
		t.Fatal("Create 應該被執行")
	}

}
