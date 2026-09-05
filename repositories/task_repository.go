package repositories

import (
	"context"
	"errors"

	"example.com/hello-service/models"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type TaskRepository struct {
	dbPool *pgxpool.Pool
}

var ErrTaskNotFound = errors.New("task not found")

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

func (r *TaskRepository) Create(
	ctx context.Context,
	title string,
) (models.Task, error) {

	var task models.Task

	err := r.dbPool.QueryRow(
		ctx,
		`
		INSERT INTO tasks (title)
		VALUES ($1)
		RETURNING id , title , completed
		`,
		title,
	).Scan(
		&task.ID,
		&task.Title,
		&task.Completed,
	)

	if err != nil {
		return models.Task{}, err
	}

	return task, nil

}

func (r *TaskRepository) Update(ctx context.Context,
	id int, title string, completed bool) (models.Task, error) {

	var task models.Task

	err := r.dbPool.QueryRow(
		ctx,
		`
		UPDATE tasks
		SET
			title = $1,
			completed = $2
		WHERE id = $3
		RETURNING id, title, completed;
		`,
		title,
		completed,
		id,
	).Scan(
		&task.ID,
		&task.Title,
		&task.Completed,
	)

	if errors.Is(err, pgx.ErrNoRows) {
		return models.Task{}, ErrTaskNotFound
	}

	if err != nil {
		return models.Task{}, err
	}

	return task, nil

}
