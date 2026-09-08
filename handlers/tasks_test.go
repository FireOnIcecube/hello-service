package handlers

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"example.com/hello-service/models"
)

type fakeTaskRepository struct {
	getAllErr   error
	createTitle string
}

func (f *fakeTaskRepository) GetAll(
	ctx context.Context,
) ([]models.Task, error) {

	if f.getAllErr != nil {
		return nil, f.getAllErr
	}

	return []models.Task{
		{
			ID:        1,
			Title:     "Task A",
			Completed: false,
		},
		{
			ID:        2,
			Title:     "Task B",
			Completed: true,
		},
	}, nil
}

func (f *fakeTaskRepository) Create(
	ctx context.Context,
	title string,
) (models.Task, error) {

	f.createTitle = title

	return models.Task{
		ID:        10,
		Title:     title,
		Completed: false,
	}, nil

}

func (f *fakeTaskRepository) Update(
	ctx context.Context,
	id int,
	title string,
	completed bool,
) (models.Task, error) {
	return models.Task{}, nil
}

func (f *fakeTaskRepository) Delete(
	ctx context.Context,
	id int,
) error {
	return nil
}

func TestCreateDbTask(t *testing.T) {

	// Arrange
	fakeRepository := fakeTaskRepository{}

	hander := New(&fakeRepository)

	request := httptest.NewRequest(
		http.MethodPost,
		"/db-task",
		strings.NewReader(
			`{"title": "Learn Testing"}`,
		),
	)

	recorder := httptest.NewRecorder()

	// Act
	hander.CreateDbTask(
		recorder,
		request,
	)

	// Assert
	if recorder.Code != http.StatusCreated {
		t.Fatalf(
			"預期 status code 為: %d , 實際得到: %d",
			http.StatusCreated,
			recorder.Code,
		)
	}

	if fakeRepository.createTitle != "Learn Testing" {
		t.Errorf(
			"預期 Repository 收到 title %q，實際得到 %q",
			"Learn Testing",
			fakeRepository.createTitle,
		)
	}

	var task models.Task

	err := json.NewDecoder(
		recorder.Body,
	).Decode(&task)

	if err != nil {
		t.Fatalf(
			"response JSON 解碼失敗: %v \n",
			err,
		)
	}

	if task.ID != 10 {
		t.Errorf("期望創建的 task id 為 %d ， 實際為: %d", 10, task.ID)
	}

	if task.Title != "Learn Testing" {
		t.Errorf("期望創建的 task title 為 %s ， 實際為: %s", "Learn Testing", task.Title)
	}

	if task.Completed != false {
		t.Errorf("期望創建的 task completed 為 %t，實際為: %t", false, task.Completed)
	}

}

func TestGetDbTasksRepositoryError(t *testing.T) {

	// Arrange
	fakeTaskRepository := fakeTaskRepository{
		getAllErr: errors.New("repositories failed"),
	}

	handler := New(&fakeTaskRepository)

	request := httptest.NewRequest(
		http.MethodGet,
		"/db-task",
		nil,
	)

	recorder := httptest.NewRecorder()

	// Act
	handler.GetDbTasks(
		recorder,
		request,
	)

	// Assert
	if recorder.Code != http.StatusInternalServerError {
		t.Fatalf(
			"預期 status %d , 實際得到 %d",
			http.StatusInternalServerError,
			recorder.Code,
		)
	}

}

func TestGetDbTasks(t *testing.T) {
	fakeRepository := &fakeTaskRepository{}

	handler := New(fakeRepository)

	request := httptest.NewRequest(
		http.MethodGet,
		"/db-task",
		nil,
	)

	recorder := httptest.NewRecorder()

	handler.GetDbTasks(
		recorder,
		request,
	)

	if recorder.Code != http.StatusOK {
		t.Fatalf(
			"預期 status %d，實際得到 %d",
			http.StatusOK,
			recorder.Code,
		)
	}

	var tasks []models.Task

	err := json.NewDecoder(
		recorder.Body,
	).Decode(&tasks)

	if err != nil {
		t.Fatalf(
			"response JSON 解碼失敗: %v",
			err,
		)
	}

	if len(tasks) != 2 {
		t.Fatalf(
			"預期 2 個 tasks ， 實際得到 %d",
			len(tasks),
		)
	}

	if tasks[0].Title != "Task A" {
		t.Errorf(
			"預期第一個 title 為 Task A，實際為 %s",
			tasks[0].Title,
		)
	}
}
