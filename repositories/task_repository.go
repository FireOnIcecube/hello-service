package repositories

import (
	"context"
	"log"

	"example.com/hello-service/models"
	"github.com/jackc/pgx/v5/pgxpool"
)

type TaskRepository struct {
	dbPool *pgxpool.Pool
}

func NewTaskRepository(
	dbPool *pgxpool.Pool,
) *TaskRepository {
	return &TaskRepository{
		dbPool: dbPool,
	}
}

func (r *TaskRepository) GetAll(
	ctx context.Context,
) ([]models.Task, error) {

	rows, err := r.dbPool.Query(
		ctx,
		"SELECT id,title,completed FROM tasks",
	)

	if err != nil {
		log.Printf(
			"查詢 tasks 失敗: %v",
			err,
		)

		return nil, err
	}

	defer rows.Close()

	tasks := make([]models.Task, 0)
	for rows.Next() {
		var task models.Task
		err := rows.Scan(
			&task.ID,
			&task.Title,
			&task.Completed,
		)

		if err != nil {
			log.Printf(
				"id: %v 的資料寫入時失敗 : %v",
				task.ID,
				err,
			)

			return nil, err
		}

		tasks = append(tasks, task)
	}

	if err := rows.Err(); err != nil {
		log.Printf("遍歷資料失敗: %v", err)

		return nil, err
	}

	return tasks, nil

}
