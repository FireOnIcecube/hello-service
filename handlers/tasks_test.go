package handlers

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strconv"
	"strings"
	"testing"

	"example.com/hello-service/domain"
	"example.com/hello-service/models"
)

type fakeTaskService struct {

	// Create
	createTaskCalled bool
	createTaskTitle  string
	createTaskErr    error

	//  GetTasks
	getTasksCalled bool
	getTasks       []models.Task
	getTasksErr    error

	// Update
	updateTaskCalled bool
	updateID         int
	updateTitle      string
	updateCompleted  bool

	updateTaskResult models.Task
	updateTaskErr    error
}

func (f *fakeTaskService) UpdateTask(
	ctx context.Context,
	id int,
	title string,
	completed bool,
) (models.Task, error) {
	f.updateTaskCalled = true
	f.updateID = id
	f.updateTitle = title
	f.updateCompleted = completed

	if f.updateTaskErr != nil {
		return models.Task{}, f.updateTaskErr
	}

	return f.updateTaskResult, nil

}

func (f *fakeTaskService) CreateTask(
	ctx context.Context,
	title string,
) (models.Task, error) {
	f.createTaskCalled = true
	f.createTaskTitle = title

	if f.createTaskErr != nil {
		return models.Task{}, f.createTaskErr
	}

	return models.Task{
		ID:        10,
		Title:     title,
		Completed: false,
	}, nil

}

func (f *fakeTaskService) GetTasks(
	ctx context.Context,
) ([]models.Task, error) {

	f.getTasksCalled = true

	if f.getTasksErr != nil {
		return []models.Task{}, f.getTasksErr
	}

	return f.getTasks, nil

}

type fakeTaskRepository struct {
	getAllErr error

	createTitle  string
	createCalled bool
	createErr    error

	updateId        int
	updateTitle     string
	updateCompleted bool
	updateCalled    bool
	updateErr       error

	deleteErr    error
	deleteCalled bool
	deleteId     int
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

	f.createCalled = true
	f.createTitle = title

	if f.createErr != nil {
		return models.Task{}, f.createErr
	}

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

	f.updateCalled = true
	f.updateId = id
	f.updateTitle = title
	f.updateCompleted = completed

	if f.updateErr != nil {

		return models.Task{}, f.updateErr
	}

	return models.Task{
		ID:        id,
		Title:     title,
		Completed: completed,
	}, nil
}

func (f *fakeTaskRepository) Delete(
	ctx context.Context,
	id int,
) error {
	f.deleteCalled = true
	f.deleteId = id

	if f.deleteErr != nil {
		return f.deleteErr
	}

	return nil
}

func TestDeleteDbTaskError(t *testing.T) {
	// Arrange
	fakeRepository := fakeTaskRepository{
		deleteErr: errors.New(" delete task error "),
	}
	fakeService := fakeTaskService{}

	handler := New(
		&fakeRepository,
		&fakeService,
	)

	request := httptest.NewRequest(http.MethodDelete, "/db-task/10", nil)
	request.SetPathValue("id", "10")

	recorder := httptest.NewRecorder()

	// Act
	handler.DeleteDbTask(recorder, request)

	// Assert
	if recorder.Code != http.StatusInternalServerError {
		t.Fatalf("預期 status code 為: %d , 實際為: %d , \n response: %s ",
			http.StatusNotFound,
			recorder.Code,
			recorder.Body.String())
	}

	if !fakeRepository.deleteCalled {
		t.Fatalf("Repository.Delete 應該被呼叫, \n response: %s", recorder.Body.String())
	}
}

func TestDeleteDbTaskNotFound(t *testing.T) {
	// Arrange
	fakeRepository := fakeTaskRepository{
		deleteErr: domain.ErrTaskNotFound,
	}
	fakeService := fakeTaskService{}

	handler := New(&fakeRepository, &fakeService)

	request := httptest.NewRequest(http.MethodDelete, "/db-task/10", nil)
	request.SetPathValue("id", "10")

	recorder := httptest.NewRecorder()

	// Act
	handler.DeleteDbTask(recorder, request)

	// Assert
	if recorder.Code != http.StatusNotFound {
		t.Fatalf("預期 status code 為: %d , 實際為: %d , \n response: %s ",
			http.StatusNotFound,
			recorder.Code,
			recorder.Body.String())
	}

	if !fakeRepository.deleteCalled {
		t.Fatal("應該呼叫 Repository.Delete")
	}
}

func TestDeleteDbTaskInvalidId(t *testing.T) {
	// Arrange
	fakeRepository := fakeTaskRepository{}
	fakeService := fakeTaskService{}

	handler := New(&fakeRepository, &fakeService)

	request := httptest.NewRequest(http.MethodDelete, "/db-task/invalid", nil)
	request.SetPathValue("id", "invalid")

	recorder := httptest.NewRecorder()

	// Act
	handler.DeleteDbTask(recorder, request)

	// Assert
	if recorder.Code != http.StatusBadRequest {
		t.Fatalf("預期 status code 為: %d , 實際 status code: %d , \n response: %s",
			http.StatusBadRequest,
			recorder.Code,
			recorder.Body.String(),
		)
	}

	if fakeRepository.deleteCalled {
		t.Fatalf("Repository.Delete 不應被呼叫, \n response: %s", recorder.Body.String())
	}
}

func TestDeleteDbTask(t *testing.T) {
	tests := []struct {
		name       string
		pathID     string
		deleteErr  error
		wantStatus int
		wantCalled bool
	}{
		{
			name:       "success",
			pathID:     "10",
			deleteErr:  nil,
			wantStatus: http.StatusNoContent,
			wantCalled: true,
		},
		{
			name:       "invalid id",
			pathID:     "invalid",
			deleteErr:  nil,
			wantStatus: http.StatusBadRequest,
			wantCalled: false,
		}, {
			name:       "task not found",
			pathID:     "10",
			deleteErr:  domain.ErrTaskNotFound,
			wantStatus: http.StatusNotFound,
			wantCalled: true,
		}, {
			name:       "repository error",
			pathID:     "10",
			deleteErr:  errors.New("delete failed"),
			wantStatus: http.StatusInternalServerError,
			wantCalled: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			// Arrange
			fakeRepository := fakeTaskRepository{
				deleteErr: tt.deleteErr,
			}
			fakeService := fakeTaskService{}

			handler := New(&fakeRepository, &fakeService)

			request := httptest.NewRequest(http.MethodDelete, "/db-task/"+tt.pathID, nil)
			request.SetPathValue("id", tt.pathID)

			recorder := httptest.NewRecorder()

			// Act
			handler.DeleteDbTask(recorder, request)

			// Assert
			if recorder.Code != tt.wantStatus {
				t.Fatalf(
					"預期 status %d，實際得到 %d",
					tt.wantStatus,
					recorder.Code,
				)
			}

			if fakeRepository.deleteCalled != tt.wantCalled {
				t.Errorf(
					"預期 deleteCalled=%t，實際為 %t",
					tt.wantCalled,
					fakeRepository.deleteCalled,
				)
			}

			if tt.wantCalled {

				if tt.pathID != strconv.Itoa(fakeRepository.deleteId) {
					t.Errorf(
						"預期操作 id 為 %v ， 實際操作 id 為 %v",
						tt.pathID,
						fakeRepository.deleteId,
					)
				}
			}

		})

	}
}

// func TestDeleteDbTask(t *testing.T) {
// 	// Arrage
// 	fakeRepository := fakeTaskRepository{}
// 	handler := New(&fakeRepository)

// 	request := httptest.NewRequest(http.MethodDelete, "/db-task/10", nil)
// 	request.SetPathValue("id", "10")

// 	recorder := httptest.NewRecorder()

// 	// Act
// 	handler.DeleteDbTask(recorder, request)

// 	// Assert
// 	if recorder.Code != http.StatusNoContent {
// 		t.Fatalf("預期 status code 為: %d , 實際收到: %d , response: %s",
// 			http.StatusNoContent,
// 			recorder.Code,
// 			recorder.Body.String())
// 	}

// 	if !fakeRepository.deleteCalled {
// 		t.Fatal("應該呼叫 Repository.Delete")
// 	}

// 	if fakeRepository.deleteId != 10 {
// 		t.Errorf(
// 			"預期 Repository.Delete 收到 id %d，實際得到 %d",
// 			10,
// 			fakeRepository.deleteId,
// 		)
// 	}

// 	if recorder.Body.Len() != 0 {
// 		t.Errorf("不應有任何回傳值 , response: %s \n", recorder.Body.String())
// 	}
// }

func TestUpdateDbTaskNotFound(t *testing.T) {
	// Arrange
	spyTask := models.Task{
		ID:        10,
		Title:     "not found",
		Completed: false,
	}

	fakeRepo := fakeTaskRepository{}
	fakeService := fakeTaskService{
		updateTaskErr: domain.ErrTaskNotFound,
	}

	handler := New(&fakeRepo, &fakeService)

	request := httptest.NewRequest(http.MethodPut,
		fmt.Sprintf(`/db-task/%d`, spyTask.ID),
		strings.NewReader(
			fmt.Sprintf(`{"title": %q , "completed": %t}`,
				spyTask.Title, spyTask.Completed),
		),
	)
	request.SetPathValue(
		"id",
		strconv.Itoa(spyTask.ID),
	)

	recorder := httptest.NewRecorder()

	// Act
	handler.UpdateDbTask(recorder, request)

	// Assert
	if recorder.Code != http.StatusNotFound {
		t.Fatalf("預期 status code 為 %d ， 實際得到 %d ， response: %s",
			http.StatusNotFound,
			recorder.Code,
			recorder.Body.String())
	}

	if fakeService.updateTaskCalled == false {
		t.Fatal("Service 應該被呼叫")
	}
}

func TestUpdateDbTaskPositiveInteger(t *testing.T) {

	// Arrange
	spyTask := models.Task{
		ID:        10,
		Title:     "positive",
		Completed: false,
	}

	fakeRepo := fakeTaskRepository{}
	fakeService := fakeTaskService{
		updateTaskErr: domain.ErrTaskPositiveInteger,
	}

	handler := New(&fakeRepo, &fakeService)

	request := httptest.NewRequest(
		http.MethodPut,
		fmt.Sprintf("/db-task/%d", spyTask.ID),
		strings.NewReader(
			fmt.Sprintf(`{"title": %q , "completed": %t}`,
				spyTask.Title,
				spyTask.Completed,
			),
		),
	)
	request.SetPathValue("id", strconv.Itoa(spyTask.ID))

	recorder := httptest.NewRecorder()

	// Act
	handler.UpdateDbTask(recorder, request)

	// Assert
	if recorder.Code != http.StatusBadRequest {
		t.Fatalf("預期 status code 為 %d ， 實際得到 %d ， response: %s",
			http.StatusBadRequest,
			recorder.Code,
			recorder.Body.String())
	}

	if fakeService.updateTaskCalled == false {
		t.Fatal("Service 應該被呼叫")
	}
}

func TestUpdateDbTaskTitleRequired(t *testing.T) {

	// Arrange
	spyTask := models.Task{
		ID:        10,
		Title:     "anything",
		Completed: false,
	}

	fakeRepo := fakeTaskRepository{}
	fakeService := fakeTaskService{
		updateTaskErr: domain.ErrTaskTitleRequired,
	}

	handler := New(&fakeRepo, &fakeService)

	request := httptest.NewRequest(
		http.MethodPut,
		fmt.Sprintf("/db-task/%d", spyTask.ID),
		strings.NewReader(
			fmt.Sprintf(`{"title": %q , "completed": %t }`,
				spyTask.Title, spyTask.Completed),
		),
	)
	request.SetPathValue("id",
		strconv.Itoa(spyTask.ID))

	recorder := httptest.NewRecorder()

	// Act
	handler.UpdateDbTask(recorder, request)

	// Assert
	if recorder.Code != http.StatusBadRequest {
		t.Fatalf("預期 status code 為 %d ， 實際得到 %d ， response: %s",
			http.StatusBadRequest,
			recorder.Code,
			recorder.Body.String())
	}

	if fakeService.updateTaskCalled == false {
		t.Fatal("Service 應該被呼叫")
	}

}

func TestUpdateDbTaskInvalidJson(t *testing.T) {
	// Arrange
	fakeRepository := fakeTaskRepository{}
	fakeService := fakeTaskService{}

	handler := New(&fakeRepository, &fakeService)

	request := httptest.NewRequest(http.MethodPut,
		"/db-test/10",
		strings.NewReader(`{"title":"invaljson}`),
	)

	request.SetPathValue("id", "10")
	recorder := httptest.NewRecorder()

	// Act
	handler.UpdateDbTask(recorder, request)

	// Assert
	if recorder.Code != http.StatusBadRequest {
		t.Fatalf(
			"預期 status code 為 %d，實際得到 %d，response: %s",
			http.StatusBadRequest,
			recorder.Code,
			recorder.Body.String(),
		)

	}

	if fakeService.updateTaskCalled {
		t.Fatal("不應呼叫 Service")
	}

}

func TestUpdateDbTaskInvalidId(t *testing.T) {
	// Arrange
	fakeRepository := fakeTaskRepository{}
	fakeService := fakeTaskService{}

	handler := New(&fakeRepository, &fakeService)

	var invalidId = "bad"

	request := httptest.NewRequest(
		http.MethodPut,
		fmt.Sprintf("/db-task/%v", invalidId),
		strings.NewReader(
			`{"title":"invalid id request"}`,
		))

	request.SetPathValue(
		"id",
		invalidId,
	)

	recorder := httptest.NewRecorder()

	// Act
	handler.UpdateDbTask(recorder, request)

	// Assert

	if recorder.Code != http.StatusBadRequest {
		t.Fatalf(
			"預期 status code 為 %d，實際得到 %d，response: %s",
			http.StatusBadRequest,
			recorder.Code,
			recorder.Body.String(),
		)
	}

	if fakeService.updateTaskCalled == true {
		t.Fatal(
			"請求 ID 格式錯誤時，不應呼叫 Service",
		)
	}
}

func TestUpdateDbTaskError(t *testing.T) {

	// Arrange
	fakeRepository := fakeTaskRepository{}

	serviceErr := errors.New("service failed")
	fakeService := fakeTaskService{
		updateTaskErr: serviceErr,
	}

	handler := New(&fakeRepository, &fakeService)

	request := httptest.NewRequest(http.MethodPut, "/db-task/10", strings.NewReader(
		`{"title":"internalError"}`,
	))
	request.SetPathValue("id", "10")

	recorder := httptest.NewRecorder()

	// Act
	handler.UpdateDbTask(recorder, request)

	// Assert
	if recorder.Code != http.StatusInternalServerError {
		t.Fatalf("預期 status code 為: %d , 實際拿到 %d , response: %s ",
			http.StatusInternalServerError,
			recorder.Code,
			recorder.Body.String(),
		)
	}

	if !fakeService.updateTaskCalled {
		t.Fatal("應該呼叫 Service")
	}

}

// 正常 update 測試
func TestUpdateDbTask(t *testing.T) {

	// 測試資料
	spyTask := models.Task{
		ID:        10,
		Title:     " test buy milk ",
		Completed: true,
	}

	stubUpdateTaskResult := models.Task{
		ID:        10,
		Title:     "test buy milk",
		Completed: true,
	}

	// Arrange
	fakeRepository := fakeTaskRepository{}
	fakeService := fakeTaskService{
		updateTaskResult: stubUpdateTaskResult,
	}

	handler := New(&fakeRepository, &fakeService)

	// var targetId = 10
	// var targetTitle = "test"
	// var targetCompleted = false

	request := httptest.NewRequest(
		http.MethodPut,
		fmt.Sprintf("/db-task/%d", spyTask.ID),
		strings.NewReader(
			fmt.Sprintf(`{"title": %q , "completed": %t }`,
				spyTask.Title, spyTask.Completed),
		),
	)

	request.SetPathValue(
		"id",
		strconv.Itoa(spyTask.ID),
	)

	recorder := httptest.NewRecorder()

	// Act
	handler.UpdateDbTask(recorder, request)

	// Assert

	// 確認 spy
	if !fakeService.updateTaskCalled {
		t.Fatal("Service.UpdateTask 應該被呼叫")
	}
	if fakeService.updateID != spyTask.ID {
		t.Errorf("預期 service 接收 updateID: %d , 實際接收 updateID: %d",
			spyTask.ID,
			fakeService.updateID,
		)
	}
	if fakeService.updateTitle != spyTask.Title {
		t.Errorf("預期 service 接收 updateTitle: %s , 實際接收 updateTitle: %s",
			spyTask.Title,
			fakeService.updateTitle,
		)
	}
	if fakeService.updateCompleted != spyTask.Completed {
		t.Errorf("預期 service 接收 updateCompleted: %t , 實際接收 updateCompleted: %t",
			spyTask.Completed,
			fakeService.updateCompleted,
		)
	}

	// 確認 stub
	if recorder.Code != http.StatusOK {
		t.Fatalf(
			"預期 status code 為 %d，實際得到 %d，response: %s",
			http.StatusOK,
			recorder.Code,
			recorder.Body.String(),
		)
	}

	var task models.Task

	err := json.NewDecoder(
		recorder.Body,
	).Decode(&task)

	if err != nil {
		t.Fatalf(
			"response JSON 解碼失敗: %v",
			err,
		)
	}

	if task.ID != 10 {
		t.Errorf(
			"預期操作 id 為 %d ， 實際操作 id 為 %d",
			10,
			task.ID,
		)
	}

	if task.Title != "test buy milk" {
		t.Errorf(
			"預期輸出 title 為 %v ， 實際輸出 title 為 %v",
			"test buy milk",
			task.Title,
		)
	}

	if task.Completed != true {
		t.Errorf(
			"預期輸出 completed 為 %t ， 實際輸出 completed 為 %t",
			true,
			task.Completed,
		)
	}

	if task != stubUpdateTaskResult {
		t.Errorf(
			"預期 回傳 result 為: %+v, 實際為: %+v",
			stubUpdateTaskResult,
			task,
		)
	}

}

func TestCreateDbTaskRepositoryError(t *testing.T) {

	// Arrange
	fakeRepository := fakeTaskRepository{}
	fakeService := fakeTaskService{
		createTaskErr: errors.New(
			"create failed",
		),
	}

	handler := New(&fakeRepository, &fakeService)

	request := httptest.NewRequest(
		http.MethodPost,
		"/db-task",
		strings.NewReader(
			`{"title":"Learn Testing"}`,
		),
	)

	recorder := httptest.NewRecorder()

	// Act
	handler.CreateDbTask(
		recorder,
		request,
	)

	// Assert
	if recorder.Code != http.StatusInternalServerError {
		t.Fatalf(
			"預期 status code 為: %d , 實際得到: %d",
			http.StatusInternalServerError,
			recorder.Code,
		)
	}

	if !fakeService.createTaskCalled {
		t.Error(
			"預期 Repository.Create 被呼叫",
		)
	}

	if fakeService.createTaskTitle != "Learn Testing" {
		t.Errorf(
			"預期 Repository 收到 title %q，實際得到 %q",
			"Learn Testing",
			fakeService.createTaskTitle,
		)
	}

}

func TestCreateDbTaskInvalidJSON(t *testing.T) {
	// Arrange

	fakeRepository := fakeTaskRepository{
		createErr: errors.New(
			"create failed",
		),
	}
	fakeService := fakeTaskService{}

	handler := New(&fakeRepository, &fakeService)

	request := httptest.NewRequest(
		http.MethodPost,
		"/db-task",
		strings.NewReader(
			`{"title":`,
		),
	)

	recorder := httptest.NewRecorder()

	// Act

	handler.CreateDbTask(
		recorder,
		request,
	)

	// Assert

	if fakeService.createTaskCalled {
		t.Error(
			"JSON 無效時，不應呼叫 Repository.Create",
		)
	}

	if recorder.Code != http.StatusBadRequest {
		t.Fatalf(
			"預期 status code 為 %d，實際得到 %d",
			http.StatusBadRequest,
			recorder.Code,
		)
	}

}

func TestCreateTitleRequried(t *testing.T) {
	// Arrange

	fakeRepository := fakeTaskRepository{}
	fakeService := fakeTaskService{
		createTaskErr: domain.ErrTaskTitleRequired,
	}

	handler := New(&fakeRepository, &fakeService)

	request := httptest.NewRequest(
		http.MethodPost,
		"/db-task",
		strings.NewReader(
			`{"title":"whatever"}`,
		),
	)
	recorder := httptest.NewRecorder()

	// Act
	handler.CreateDbTask(
		recorder,
		request,
	)

	// Assert
	if recorder.Code != http.StatusBadRequest {
		t.Fatalf(
			"預期 status %d，實際得到 %d",
			http.StatusBadRequest,
			recorder.Code,
		)
	}

}

func TestCreateDbTask(t *testing.T) {

	// Arrange
	fakeRepository := fakeTaskRepository{}
	fakeService := fakeTaskService{}

	handler := New(&fakeRepository, &fakeService)

	request := httptest.NewRequest(
		http.MethodPost,
		"/db-task",
		strings.NewReader(
			`{"title": "Learn Testing"}`,
		),
	)

	recorder := httptest.NewRecorder()

	// Act
	handler.CreateDbTask(
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

	if fakeService.createTaskTitle != "Learn Testing" {
		t.Errorf(
			"預期 Repository 收到 title %q，實際得到 %q",
			"Learn Testing",
			fakeService.createTaskTitle,
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
	fakeRepository := fakeTaskRepository{
		getAllErr: errors.New("repositories failed"),
	}
	fakeService := fakeTaskService{}

	handler := New(&fakeRepository, &fakeService)

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
	fakeService := fakeTaskService{}

	handler := New(fakeRepository, &fakeService)

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
