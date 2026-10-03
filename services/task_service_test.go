package services

import (
	"context"
	"errors"
	"testing"

	"example.com/hello-service/domain"
	"example.com/hello-service/models"
)

var repositoryErr = errors.New("get tasks failed")

type fakeTaskRepository struct {
	// Create
	createCalled bool
	createTitle  string
	createErr    error

	// GetAll
	getAllCalled bool
	getAllTasks  []models.Task
	getAllErr    error

	// Update
	updateCalled    bool
	updateID        int
	updateTitle     string
	updateCompleted bool

	updateResult models.Task
	updateErr    error

	// Delete
	deleteCalled bool
	deleteID     int
	deleteErr    error
}

func (f *fakeTaskRepository) GetAll(ctx context.Context) ([]models.Task, error) {
	f.getAllCalled = true

	if f.getAllErr != nil {
		return []models.Task{}, f.getAllErr
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

func (f *fakeTaskRepository) Update(ctx context.Context,
	id int,
	title string,
	completed bool) (models.Task, error) {

	f.updateCalled = true
	f.updateID = id
	f.updateTitle = title
	f.updateCompleted = completed

	if f.updateErr != nil {
		return models.Task{}, f.updateErr
	}

	return f.updateResult, nil
}

func (f *fakeTaskRepository) Delete(ctx context.Context, id int) error {
	f.deleteCalled = true
	f.deleteID = id

	if f.deleteErr != nil {

		return f.deleteErr
	}

	return nil
}

// GetAll Repository 錯誤
func TestTaskServiceGetTasksRepositoryError(t *testing.T) {
	// Arrange
	fakeRepo := fakeTaskRepository{
		getAllErr: repositoryErr,
	}
	service := NewTaskService(&fakeRepo)

	ctx := context.Background()

	// Act
	_, err := service.GetTasks(ctx)

	// Assert
	if err == nil {
		t.Fatal("預期回報錯誤")
	}

	if !errors.Is(err, repositoryErr) {
		t.Fatalf(
			"預期回報錯誤為: %v , 實際回報錯誤: %v",
			repositoryErr,
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

	if !errors.Is(err, domain.ErrTaskTitleRequired) {
		t.Fatalf(
			"預期回報錯誤為: %v , 實際回報錯誤: %v",
			domain.ErrTaskTitleRequired,
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

// Update 正常流程
func TestTaskServiceUpdateTask(t *testing.T) {

	// update 資料
	inputTask := models.Task{
		ID:        1,
		Title:     " 買牛奶 ",
		Completed: true,
	}

	expectedResult := models.Task{
		ID:        1,
		Title:     "買牛奶",
		Completed: true,
	}

	// Arrange
	fakeRepo := fakeTaskRepository{
		updateResult: expectedResult,
	}
	service := NewTaskService(&fakeRepo)

	ctx := context.Background()

	// Act
	task, err := service.UpdateTask(ctx, inputTask.ID, inputTask.Title, inputTask.Completed)

	// Assert
	if err != nil {
		t.Fatalf("接收到錯誤，操作失敗: %v", err)
	}

	if !fakeRepo.updateCalled {
		t.Fatal("Repository.Update 應該被呼叫")
	}

	// 驗證資料
	if fakeRepo.updateID != 1 {
		t.Errorf("預期 ID 為: %d, 實際為: %d",
			1,
			fakeRepo.updateID,
		)
	}

	if fakeRepo.updateTitle != "買牛奶" {
		t.Errorf("預期 Title 為: %v, 實際為: %v",
			"買牛奶",
			fakeRepo.updateTitle,
		)
	}

	if fakeRepo.updateCompleted != true {
		t.Errorf("預期 Completed 為: %t, 實際為: %t",
			true,
			fakeRepo.updateCompleted,
		)
	}

	// Repository → Service

	if task != expectedResult {
		t.Errorf("預期 回傳 result 為: %+v, 實際為: %+v",
			expectedResult,
			task,
		)
	}

}

// Update NotFound
func TestTaskServiceUpdateTaskNotFound(t *testing.T) {

	// 測試資料
	temTask := models.Task{
		ID:        999,
		Title:     "NotFound",
		Completed: false,
	}

	// Arrange
	fakeRepo := fakeTaskRepository{updateErr: domain.ErrTaskNotFound}
	service := NewTaskService(&fakeRepo)

	ctx := context.Background()

	// Act
	_, err := service.UpdateTask(ctx, temTask.ID, temTask.Title, temTask.Completed)

	// Assert
	if err == nil {
		t.Fatal("預期應該收到錯誤")
	}

	if !errors.Is(err, domain.ErrTaskNotFound) {
		t.Fatalf("預期應該收到錯誤為: %v , 實際收到: %v",
			domain.ErrTaskNotFound,
			err,
		)
	}

	if !fakeRepo.updateCalled {
		t.Fatal("Update 應該被執行")
	}
}

// Update positive integer error
func TestTaskServiceUpdateTaskPositiveIntegerErr(t *testing.T) {

	// 測試資料
	temTask := models.Task{
		ID:        -1,
		Title:     "integer error",
		Completed: false,
	}

	// Arrange
	fakeRepo := fakeTaskRepository{updateErr: domain.ErrTaskPositiveInteger}
	service := NewTaskService(&fakeRepo)

	ctx := context.Background()

	// Act
	_, err := service.UpdateTask(ctx, temTask.ID, temTask.Title, temTask.Completed)

	// Assert
	if err == nil {
		t.Fatal("預期應該收到錯誤")
	}

	if !errors.Is(err, domain.ErrTaskPositiveInteger) {
		t.Fatalf("預期應該收到錯誤為: %v , 實際收到: %v",
			domain.ErrTaskPositiveInteger,
			err,
		)
	}

	if fakeRepo.updateCalled {
		t.Fatal("Update 不應該被執行")
	}
}

// Delete 正常流程
func TestTaskServiceDeleteTask(t *testing.T) {

	// 測試資料
	spyTaskID := 1

	// Arrange
	fakeRepository := fakeTaskRepository{}
	service := NewTaskService(&fakeRepository)

	ctx := context.Background()

	// Act
	err := service.DeleteTask(ctx, spyTaskID)

	// Assert
	if err != nil {
		t.Fatalf("接收到錯誤，操作失敗: %v", err)
	}

	if fakeRepository.deleteCalled == false {
		t.Fatal("Repository.Delete 應該被呼叫")
	}

	// 核對 spy 資料
	if fakeRepository.deleteID != 1 {
		t.Errorf("預期 repo 接收 id 為: %d , 實際接收 id: %d",
			1,
			fakeRepository.deleteID)
	}

}

func TestTaskServiceDeleteTaskNotFound(t *testing.T) {

	// 測試資料
	spyTaskID := 999

	// Arrange
	fakeRepo := fakeTaskRepository{
		deleteErr: domain.ErrTaskNotFound,
	}
	service := NewTaskService(&fakeRepo)

	ctx := context.Background()

	// Act
	err := service.DeleteTask(ctx, spyTaskID)

	// Assert
	if err == nil {
		t.Fatal("應該回報錯誤")
	}

	if !errors.Is(err, domain.ErrTaskNotFound) {
		t.Fatalf("預期回報錯誤為: %v , 實際回報錯誤為: %v",
			domain.ErrTaskNotFound, err)
	}

	if fakeRepo.deleteCalled == false {
		t.Fatalf(
			"預期 repo deletedCalled 為: %t , 實際為: %t",
			true,
			fakeRepo.deleteCalled,
		)
	}

	// 核對 spy
	if fakeRepo.deleteID != spyTaskID {
		t.Errorf("預期 repo 接收 deleteID : %d, repo 實際接收 deleteID: %d",
			spyTaskID,
			fakeRepo.deleteID,
		)
	}

}

func TestTaskServiceDeleteTaskPositiveInteger(t *testing.T) {
	// 測試資料
	spyTaskID := -1

	// Arrange
	fakeRepo := fakeTaskRepository{}
	service := NewTaskService(&fakeRepo)

	ctx := context.Background()

	// Act
	err := service.DeleteTask(ctx, spyTaskID)

	// Assert
	if err == nil {
		t.Fatal("應該回報錯誤")
	}

	if !errors.Is(err, domain.ErrTaskPositiveInteger) {
		t.Fatalf("預期回報錯誤為: %v , 實際回報錯誤為: %v",
			domain.ErrTaskPositiveInteger, err)
	}

	if fakeRepo.deleteCalled == true {
		t.Fatalf(
			"預期 repo deletedCalled 為: %t , 實際為: %t",
			false,
			fakeRepo.deleteCalled,
		)
	}

}

func TestTaskServiceDeleteTaskRepositoryError(t *testing.T) {
	// 測試資料
	spyTaskID := 1
	spyErr := errors.New("delete repo error")

	// Arrange
	fakeRepo := fakeTaskRepository{
		deleteErr: spyErr,
	}
	service := NewTaskService(&fakeRepo)

	ctx := context.Background()

	// Act
	err := service.DeleteTask(ctx, spyTaskID)

	// Assert
	if err == nil {
		t.Fatal("應該回報錯誤")
	}

	if !errors.Is(err, spyErr) {
		t.Fatalf("預期回報錯誤為: %v , 實際回報錯誤為: %v",
			spyErr, err)
	}

	if fakeRepo.deleteCalled == false {
		t.Fatalf(
			"預期 repo deletedCalled 為: %t , 實際為: %t",
			true,
			fakeRepo.deleteCalled,
		)
	}

}
