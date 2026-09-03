package repositories

import (
	"context"

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

			return nil, err
		}

		tasks = append(tasks, task)
	}

	if err := rows.Err(); err != nil {

		return nil, err
	}

	return tasks, nil

}
