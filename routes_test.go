package main

import (
	"context"
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

func TestRouterGetDbTasks(t *testing.T) {
	// Arrange
	fakeRepo := &fakeTaskRepository{}
	handlers := handlers.New(fakeRepo)
	routers := newRouter(handlers)

	request := httptest.NewRequest(http.MethodPatch, "/db-tasexitk", nil)
	recorder := httptest.NewRecorder()

	// Act
	routers.ServeHTTP(
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

}
