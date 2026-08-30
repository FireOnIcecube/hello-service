package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	"strconv"

	"example.com/hello-service/models"

	"github.com/jackc/pgx/v5"
)

func (app *App) deleteDbTaskHandler(
	w http.ResponseWriter,
	r *http.Request,
) {

	idString := r.PathValue("id")

	id, err := strconv.Atoi(idString)

	if err != nil {
		http.Error(
			w,
			"傳遞 id 非數字",
			http.StatusBadRequest,
		)
		return
	}

	var deletedTaskID int

	err = app.dbPool.QueryRow(
		r.Context(),
		`
		DELETE FROM tasks
		WHERE id = $1
		RETURNING id;
		`,
		id,
	).Scan(
		&deletedTaskID,
	)

	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			http.Error(
				w,
				"刪除資料不存在",
				http.StatusNotFound,
			)
			return
		}

		log.Printf(
			"DELETE task 失敗: %v",
			err,
		)

		http.Error(
			w,
			"伺服器內部錯誤",
			http.StatusInternalServerError,
		)
		return
	}

	w.WriteHeader(
		http.StatusNoContent,
	)
}

func (app *App) updateDbTaskHandler(
	w http.ResponseWriter,
	r *http.Request,
) {

	idString := r.PathValue("id")
	id, err := strconv.Atoi(idString)

	if err != nil {
		http.Error(
			w,
			"傳遞 id 非數字",
			http.StatusBadRequest,
		)

		log.Println(err)

		return
	}

	var input UpdateTaskRequest

	err = json.NewDecoder(r.Body).Decode(&input)

	if err != nil {
		http.Error(
			w,
			"傳遞資料格式不正確",
			http.StatusBadRequest,
		)
		log.Println(err)

		return
	}

	var task models.Task
	err = app.dbPool.QueryRow(
		r.Context(),
		`
		UPDATE tasks
		SET
			title = $1,
			completed = $2
		WHERE id = $3
		RETURNING id, title, completed;
		`,

		input.Title,
		input.Completed,
		id,
	).Scan(
		&task.ID,
		&task.Title,
		&task.Completed,
	)

	if err != nil {

		if errors.Is(err, pgx.ErrNoRows) {
			http.Error(
				w,
				"修改資料不存在",
				http.StatusNotFound,
			)

			log.Printf("找不到對應 id: %v \n", err)
			return
		}

		http.Error(
			w,
			"伺服器內部錯誤",
			http.StatusInternalServerError,
		)
		log.Println(err)

		return
	}

	w.Header().Set(
		"Content-Type",
		"application/json",
	)

	w.WriteHeader(
		http.StatusOK,
	)

	err = json.NewEncoder(w).Encode(task)

	if err != nil {
		http.Error(
			w,
			"伺服器內部錯誤",
			http.StatusInternalServerError,
		)
		log.Printf("json 編碼錯誤:%v \n", err)
		return
	}

}

func (app *App) createDbTaskHandler(
	w http.ResponseWriter,
	r *http.Request,

) {

	var input CreateTaskRequest

	err := json.NewDecoder(r.Body).Decode(&input)

	if err != nil {
		http.Error(
			w,
			"傳遞參數格式不正確",
			http.StatusBadRequest,
		)
		return
	}

	var task models.Task

	err = app.dbPool.QueryRow(
		r.Context(),
		`
		INSERT INTO tasks (title)
		VALUES ($1) 
		RETURNING id , title , completed
		`,
		input.Title,
	).Scan(
		&task.ID,
		&task.Title,
		&task.Completed,
	)

	if err != nil {
		http.Error(
			w,
			fmt.Sprintf("SQL 查詢錯誤: %v", err),
			http.StatusInternalServerError,
		)
		return
	}

	w.Header().Set(
		"Content-Type",
		"application/json",
	)

	w.WriteHeader(
		http.StatusCreated,
	)

	err = json.NewEncoder(w).Encode(task)

}
