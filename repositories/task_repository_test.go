package repositories

import (
	"context"
	"os"
	"testing"

	"github.com/jackc/pgx/v5/pgxpool"
)

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
