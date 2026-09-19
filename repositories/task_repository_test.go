package repositories

import (
	"context"
	"fmt"
	"os"
	"testing"

	"example.com/hello-service/models"
	"github.com/jackc/pgx/v5/pgxpool"
)

// setup function
func setupTestRepository(t *testing.T) (context.Context,
	*pgxpool.Pool,
	*TaskRepository) {

	t.Helper()

	databaseURL, ok := os.LookupEnv("TEST_DATABASE_URL")

	if !ok || databaseURL == "" {
		t.Skip(
			"TEST_DATABASE_URL 未設定，跳過 integration",
		)
	}

	ctx := context.Background()

	dbPool, err := pgxpool.New(
		ctx,
		databaseURL,
	)

	if err != nil {
		t.Fatalf(
			"建立 test db pool 失敗: %v",
			err,
		)
	}

	t.Cleanup(func() {
		dbPool.Close()
	})

	// 清空資料
	_, err = dbPool.Exec(ctx,
		"TRUNCATE TABLE tasks RESTART IDENTITY",
	)

	if err != nil {
		t.Fatalf(
			"清理 test tasks 失敗: %v",
			err,
		)
	}

	taskRepo := NewTaskRepository(dbPool)

	return ctx, dbPool, taskRepo

}

// List
func TestTaskRepositoryGetAll(t *testing.T) {

	// Arrange
	ctx, _, taskRepo := setupTestRepository(t)

	for i := 0; i < 2; i++ {
		_, err := taskRepo.Create(
			ctx,
			fmt.Sprintf("Test Task for GetAll: %v", i),
		)

		if err != nil {
			t.Fatalf(
				"放入 index 為 %d 的測試資料時發生錯誤: %v",
				i,
				err,
			)
		}
	}

	// Act
	tasks, err := taskRepo.GetAll(ctx)

	// Assert

	// 檢查是否有錯誤

	if err != nil {
		t.Fatalf(
			"GetAll() 執行失敗: %v",
			err,
		)
	}
	// 檢查資料筆數

	if len(tasks) != 2 {
		t.Fatalf(
			"預期取得資料筆數: %d, 實際取得資料筆數: %d",
			2,
			len(tasks),
		)
	}

	// 核對資料

	for index, task := range tasks {
		if task.ID != index+1 {
			t.Errorf(
				"預期 task id 為: %d , 實際為 %d",
				index+1,
				task.ID,
			)

		}

		if task.Title != fmt.Sprintf("Test Task for GetAll: %v", index) {
			t.Errorf(
				"預期 task title 為: %s , 實際為 %s",
				fmt.Sprintf("Test Task for GetAll: %v", index),
				task.Title,
			)
		}

		if task.Completed != false {
			t.Errorf(
				"預期 task completed 為: %t , 實際為 %t",
				false,
				task.Completed,
			)
		}

	}

}

// Update
func TestTaskRepositoryUpdate(t *testing.T) {
	// Arrange
	ctx, dbPool, taskRepo := setupTestRepository(t)

	// 建立初始資料
	task, err := taskRepo.Create(ctx, "task for update test ")

	// 預期更改後的資料
	expectedTask := models.Task{
		ID:        task.ID,
		Title:     "After Update",
		Completed: true,
	}

	if err != nil {
		t.Fatalf(
			"建立初始資料失敗: %v",
			err,
		)
	}

	// Act
	updatedTask, err := taskRepo.Update(ctx, expectedTask.ID, expectedTask.Title, expectedTask.Completed)

	// Assert

	if err != nil {
		t.Fatalf(
			"修改資料時發生錯誤: %v",
			err,
		)
	}

	// 檢查 update 回傳的資料

	if updatedTask.ID != expectedTask.ID {
		t.Errorf("預期回傳 id 為: %d , 實際為: %d",
			expectedTask.ID,
			updatedTask.ID)
	}

	if updatedTask.Title != expectedTask.Title {
		t.Errorf("預期回傳 title 為: %s , 實際為: %s",
			expectedTask.Title,
			updatedTask.Title,
		)
	}

	if updatedTask.Completed != expectedTask.Completed {
		t.Errorf("預期回傳 Completed 為: %t , 實際為: %t",
			expectedTask.Completed,
			updatedTask.Completed,
		)
	}

	// 從 postgreSQL 獲取資料檢查

	var dbTask models.Task

	err = dbPool.QueryRow(
		ctx,
		`SELECT id, title , completed FROM tasks WHERE id = $1`,
		task.ID,
	).Scan(
		&dbTask.ID,
		&dbTask.Title,
		&dbTask.Completed,
	)

	if err != nil {
		t.Fatalf(
			"調用資料庫失敗: %v",
			err,
		)
	}

	if dbTask.ID != expectedTask.ID {
		t.Errorf("預期資料庫內 task id 為: %d , 實際為: %d",
			expectedTask.ID,
			dbTask.ID,
		)
	}

	if dbTask.Title != expectedTask.Title {
		t.Errorf("預期資料庫內 task title 為: %s , 實際為: %s",
			expectedTask.Title,
			dbTask.Title,
		)
	}

	if dbTask.Completed != expectedTask.Completed {
		t.Errorf("預期資料庫內 task Completed 為: %t , 實際為: %t",
			expectedTask.Completed,
			dbTask.Completed,
		)
	}

}

// Create
func TestTaskRepositoryCreate(t *testing.T) {

	ctx, dbPool, taskRepo := setupTestRepository(t)

	task, err := taskRepo.Create(
		ctx,
		"Integration Test Task",
	)

	if err != nil {
		t.Fatalf(
			"Create() 發生錯誤: %v",
			err,
		)
	}

	if task.ID <= 0 {
		t.Errorf(
			"預期 ID > 0，實際為 %d",
			task.ID,
		)
	}

	if task.Title != "Integration Test Task" {
		t.Errorf("預期操作 task title 為: %v , 實際為: %v", "Integration Test Task", task.Title)
	}

	if task.Completed {
		t.Errorf("預期操作 task compelted 為: %t , 實際為: %t", false, task.Completed)
	}

	// 確認資料是否存在

	var savedTitle string

	err = dbPool.QueryRow(
		ctx,
		`
	SELECT title
	FROM tasks
	WHERE id = $1
	`,
		task.ID,
	).Scan(
		&savedTitle,
	)

	if err != nil {
		t.Fatalf(
			"查詢建立後的 task 失敗: %v",
			err,
		)
	}

	if savedTitle != "Integration Test Task" {
		t.Errorf(
			"預期 DB title %q，實際為 %q",
			"Integration Test Task",
			savedTitle,
		)
	}

}
