package services

import (
	"context"
	"errors"
	"testing"

	"example.com/hello-service/models"
)

type fakeTaskRepository struct {
	// Create
	createCalled bool
	createTitle  string
	createErr    error

	// GetAll
	getAllCalled bool
	getAllTasks  []models.Task
	getAllErr    error
}

func (f *fakeTaskRepository) GetAll(ctx context.Context) ([]models.Task, error) {
	f.getAllCalled = true

	if f.getAllErr != nil {
		return []models.Task{}, f.createErr
	}

	return f.getAllTasks, nil
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

// GetAll Repository 錯誤
func TestTaskServiceGetTasksRepositoryError(t *testing.T) {
	// Arrange
	fakeRepo := fakeTaskRepository{
		getAllErr: ErrRepository,
	}
	service := NewTaskService(&fakeRepo)

	ctx := context.Background()

	// Act
	_, err := service.GetTasks(ctx)

	// Assert
	if err == nil {
		t.Fatal("預期回報錯誤")
	}

	if !errors.Is(err, ErrRepository) {
		t.Fatalf(
			"預期回報錯誤為: %v , 實際回報錯誤: %v",
			ErrRepository,
			err,
		)
	}

}

// GetAll 正常流程
func TestTaskServiceGetTasks(t *testing.T) {

	// Arrange

	// 測試資料
	testTasks := []models.Task{
		{
			ID:        1,
			Title:     "test title 1",
			Completed: true,
		},
		{
			ID:        100,
			Title:     "test title 100",
			Completed: false,
		},
	}

	fakeRepo := fakeTaskRepository{
		getAllTasks: testTasks,
	}
	service := NewTaskService(&fakeRepo)

	ctx := context.Background()

	// Act
	tasks, err := service.GetTasks(ctx)

	// Assert
	if err != nil {
		t.Fatalf(
			"Create Task 執行失敗: %v",
			err,
		)
	}

	if fakeRepo.getAllCalled == false {
		t.Fatalf(
			"預期 getAllCalled 為: %t , 實際為: %t",
			true,
			fakeRepo.getAllCalled,
		)
	}

	if fakeRepo.getAllErr != nil {
		t.Fatalf(
			"接收到錯誤 err: %v",
			err,
		)
	}

	if len(tasks) != len(testTasks) {
		t.Fatalf(
			"資料長度不一， 返回資料長度: %d , 預期資料長度: %d",
			len(tasks),
			len(testTasks),
		)
	}

	// 判斷資料是否一致
	for idx := range testTasks {

		if tasks[idx].ID != testTasks[idx].ID {
			t.Errorf(
				"返回資料的第 index: %d ,ID 和測試資料不符, 預期: %d , 實際: %d",
				idx,
				testTasks[idx].ID,
				tasks[idx].ID,
			)
		}

		if tasks[idx].Title != testTasks[idx].Title {
			t.Errorf(
				"返回資料的第 index: %d ,Title 和測試資料不符, 預期: %v , 實際: %v",
				idx,
				testTasks[idx].Title,
				tasks[idx].Title,
			)
		}

		if tasks[idx].Completed != testTasks[idx].Completed {
			t.Errorf(
				"返回資料的第 index: %d ,Completed 和測試資料不符, 預期: %t , 實際: %t",
				idx,
				testTasks[idx].Completed,
				tasks[idx].Completed,
			)
		}
	}
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

	if task.Title != "Learn Go" {
		t.Errorf(
			"預期創建 task title: %s , 實際 title: %s",
			"Learn Go",
			task.Title,
		)
	}

	if fakeRepo.createTitle != "Learn Go" {
		t.Errorf(
			"預期 Repository 收到 title %q，實際為 %q",
			"Learn Go",
			fakeRepo.createTitle,
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

	repositoryErr := errors.New("database failed")

	fakeRepo := fakeTaskRepository{
		createErr: repositoryErr,
	}
	service := NewTaskService(&fakeRepo)

	ctx := context.Background()

	// Act
	_, err := service.CreateTask(ctx, " internal error ")

	// Assert
	if err == nil {
		t.Fatal("預期回報錯誤")
	}

	if !errors.Is(err, repositoryErr) {
		t.Fatalf(
			"預期錯誤 %v，實際為 %v",
			repositoryErr,
			err,
		)
	}

	if !fakeRepo.createCalled {
		t.Fatal("Create 應該被執行")
	}

}
