package repositories

import (
	"context"
	"fmt"
	"os"
	"testing"

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
		defer dbPool.Close()
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
	// 目標: 測試 List 功能是否正常運作

	// Arrange
	// 連接資料庫
	databaseURL, ok := os.LookupEnv(
		"TEST_DATABASE_URL",
	)

	if !ok || databaseURL == "" {
		t.Skip(
			"TEST_DATABASE_URL 未設定，跳過 integration test",
		)
	}

	// 建立 db Pool
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

	defer dbPool.Close()

	// 清空資料
	_, err = dbPool.Exec(
		ctx,
		"TRUNCATE TABLE tasks RESTART IDENTITY",
	)

	if err != nil {
		t.Fatalf(
			"清理 test tasks 失敗: %v",
			err,
		)
	}

	// 放入測試資料

	repositories := NewTaskRepository(dbPool)

	for i := 0; i < 2; i++ {
		_, err := repositories.Create(
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
	tasks, err := repositories.GetAll(ctx)

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

// Create
func TestTaskRepositoryCreate(t *testing.T) {
	databaseURL, ok := os.LookupEnv(
		"TEST_DATABASE_URL",
	)

	if !ok || databaseURL == "" {
		t.Skip(
			"TEST_DATABASE_URL 未設定，跳過 integration test",
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

	defer dbPool.Close()

	// 清空資料
	_, err = dbPool.Exec(
		ctx,
		"DELETE FROM tasks",
	)

	if err != nil {
		t.Fatalf(
			"清理 test tasks 失敗: %v",
			err,
		)
	}

	repositories := NewTaskRepository(dbPool)

	task, err := repositories.Create(
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
