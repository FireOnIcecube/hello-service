package main

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"example.com/hello-service/handlers"
	"example.com/hello-service/models"
)

type fakeTaskRepository struct {
}

func (f *fakeTaskRepository) GetAll(ctx context.Context,
) ([]models.Task, error) {

	return []models.Task{
		{
			ID:        1,
			Title:     "title no.1",
			Completed: false,
		},
		{
			ID:        2,
			Title:     "title no.2",
			Completed: true,
		},
	}, nil
}

func (f *fakeTaskRepository) Create(ctx context.Context,
	title string) (models.Task, error) {

	return models.Task{}, nil
}

func (f *fakeTaskRepository) Update(ctx context.Context,
	id int, title string, completed bool,
) (models.Task, error) {
	return models.Task{}, nil
}

func (f *fakeTaskRepository) Delete(ctx context.Context,
	id int) error {
	return nil
}

// NotFound
func TestRouterNotFound(t *testing.T) {
	// Arrange
	fakeRepo := &fakeTaskRepository{}
	handler := handlers.New(fakeRepo)
	router := newRouter(handler)

	request := httptest.NewRequest(http.MethodGet, "/not-exist", nil)
	recorder := httptest.NewRecorder()

	// Act
	router.ServeHTTP(
		recorder,
		request,
	)

	// Assert
	if recorder.Code != http.StatusNotFound {
		t.Fatalf(
			"預期 status code 為: %d , 實際為: %d , \n response: %s ",
			http.StatusNotFound,
			recorder.Code,
			recorder.Body.String(),
		)
	}
}

// NotAllowed
func TestRouterMethodNotAllowed(t *testing.T) {
	// Arrange
	fakeRepo := &fakeTaskRepository{}
	handler := handlers.New(fakeRepo)
	router := newRouter(handler)

	request := httptest.NewRequest(http.MethodPatch, "/db-task", nil)
	recorder := httptest.NewRecorder()

	// Act
	router.ServeHTTP(
		recorder,
		request,
	)

	// Assert
	if recorder.Code != http.StatusMethodNotAllowed {
		t.Fatalf(
			"預期 status code 為: %d , 實際為: %d , \n response: %s ",
			http.StatusMethodNotAllowed,
			recorder.Code,
			recorder.Body.String(),
		)
	}

}

// Get
func TestRouterGetDbTasks(t *testing.T) {
	// Arrange
	fakeRepo := &fakeTaskRepository{}
	handler := handlers.New(fakeRepo)
	router := newRouter(handler)

	request := httptest.NewRequest(http.MethodGet, "/db-task", nil)
	recorder := httptest.NewRecorder()

	// Act
	router.ServeHTTP(
		recorder,
		request,
	)

	// Assert
	if recorder.Code != http.StatusOK {
		t.Fatalf("預期 status code 為: %d , 實際為: %d , \n response: %s ",
			http.StatusOK,
			recorder.Code,
			recorder.Body.String())
	}

	var tasks []models.Task

	err := json.NewDecoder(
		recorder.Body,
	).Decode(&tasks)

	if err != nil {
		t.Fatalf(
			"解析 response body 失敗: %v",
			err,
		)
	}

	if len(tasks) != 2 {
		t.Errorf(
			"預期資料筆數為: %d , 實際為: %d , response: %s",
			2,
			len(tasks),
			recorder.Body.String(),
		)
	}

	if tasks[0].Title != "title no.1" {
		t.Errorf(
			"預期 tasks[0]  title 為: %s , 實際為: %s , response: %s",
			"title no.1",
			tasks[0].Title,
			recorder.Body.String(),
		)
	}

}
