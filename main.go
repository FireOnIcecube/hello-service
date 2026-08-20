package main

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	"strconv"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type Task struct {
	ID        int    `json:"id"`
	Title     string `json:"title"`
	Completed bool   `json:"completed"`
}

type CreateTaskRequest struct {
	Title string `json:"title"`
}

type UpdateTaskRequest struct {
	Title     string `json:"title"`
	Completed bool   `json:"completed"`
}

type DeleteTaskRequest struct {
	ID int `json:"id"`
}

var tasks = []Task{
	{
		ID:        1,
		Title:     "學習 Go JSON",
		Completed: false,
	},
	{
		ID:        2,
		Title:     "建立第一個 API",
		Completed: true,
	},
}

var nextID = 3
var dbPool *pgxpool.Pool

func deleteDbTaskHandler(
	w http.ResponseWriter,
	r *http.Request,
) {
	databaseURL :=
		"postgres://task_user:task_password@127.0.0.1:5433/task_db?sslmode=disable"

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

	conn, err := pgx.Connect(
		r.Context(),
		databaseURL,
	)

	if err != nil {
		log.Printf(
			"資料庫連線失敗: %v",
			err,
		)

		http.Error(
			w,
			"伺服器內部錯誤",
			http.StatusInternalServerError,
		)
		return
	}

	defer conn.Close(context.Background())

	var deletedTaskID int

	err = conn.QueryRow(
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

func updateDbTaskHandler(
	w http.ResponseWriter,
	r *http.Request,
) {
	databaseURL := "postgres://task_user:task_password@127.0.0.1:5433/task_db?sslmode=disable"

	idString := r.PathValue("id")
	id, err := strconv.Atoi(idString)

	if err != nil {
		http.Error(
			w,
			"傳遞 id 非數字",
			http.StatusBadRequest,
		)

		log.Fatalln(err)

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
		log.Fatalln(err)

		return
	}

	conn, err := pgx.Connect(r.Context(), databaseURL)

	if err != nil {
		http.Error(
			w,
			"伺服器內部錯誤",
			http.StatusInternalServerError,
		)

		log.Fatalln(err)

		return
	}

	defer conn.Close(context.Background())

	var task Task
	err = conn.QueryRow(
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

			log.Fatalf("找不到對應 id: %v \n", err)
			return
		}

		http.Error(
			w,
			"伺服器內部錯誤",
			http.StatusInternalServerError,
		)
		log.Fatalln(err)

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
		log.Fatalf("json 編碼錯誤:%v \n", err)
		return
	}

}

func createDbTaskHandler(
	w http.ResponseWriter,
	r *http.Request,

) {
	databaseURL := "postgres://task_user:task_password@127.0.0.1:5433/task_db?sslmode=disable"

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

	conn, err := pgx.Connect(r.Context(), databaseURL)

	if err != nil {
		http.Error(
			w,
			fmt.Sprintf("資料庫連線失敗: %v", err),
			http.StatusInternalServerError,
		)
		return
	}

	defer conn.Close(context.Background())

	var task Task

	err = conn.QueryRow(
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

func getDbTasksHandler(
	w http.ResponseWriter,
	r *http.Request,
) {
	databaseURL :=
		"postgres://task_user:task_password@127.0.0.1:5433/task_db?sslmode=disable"

	conn, err := pgx.Connect(r.Context(), databaseURL)

	if err != nil {
		http.Error(
			w,
			fmt.Sprintf("資料庫連線失敗: %v", err),
			http.StatusInternalServerError,
		)
		return
	}

	defer conn.Close(context.Background())

	rows, err := conn.Query(
		r.Context(),
		"SELECT id,title , completed FROM tasks",
	)

	if err != nil {
		http.Error(
			w,
			fmt.Sprintf("SQL 執行失敗: %v", err),
			http.StatusInternalServerError,
		)

		return
	}

	tasks := make([]Task, 0)
	for rows.Next() {
		var task Task
		err := rows.Scan(
			&task.ID,
			&task.Title,
			&task.Completed,
		)

		if err != nil {
			http.Error(
				w,
				"讀取資料失敗",
				http.StatusInternalServerError,
			)
			return
		}

		tasks = append(tasks, task)
	}

	if err := rows.Err(); err != nil {
		http.Error(
			w,
			"遍歷資料失敗",
			http.StatusInternalServerError,
		)

		return
	}

	w.Header().Set(
		"Content-Type",
		"application/json",
	)

	err = json.NewEncoder(w).Encode(tasks)

	if err != nil {
		log.Printf(
			"JSON 編碼失敗: %v", err,
		)
	}

}

func deleteTaskHandler(w http.ResponseWriter, r *http.Request) {
	idString := r.PathValue("id")

	id, err := strconv.Atoi(idString)
	if err != nil {
		http.Error(
			w,
			"輸入 id 非數字",
			http.StatusBadRequest,
		)
		return
	}

	for i := range tasks {
		if tasks[i].ID == id {
			tasks = append(
				tasks[:i],
				tasks[i+1:]...,
			)

			w.WriteHeader(
				http.StatusNoContent,
			)

			return
		}
	}

	http.Error(
		w,
		"找不到刪除目標",
		http.StatusBadRequest,
	)

}

func updateTaskHandler(
	w http.ResponseWriter,
	r *http.Request,
) {
	var idString = r.PathValue("id")

	id, err := strconv.Atoi(idString)
	if err != nil {
		http.Error(
			w,
			"輸入 id 不為數字",
			http.StatusBadRequest,
		)
		return
	}

	var input UpdateTaskRequest

	err = json.NewDecoder(r.Body).Decode(&input)
	if err != nil {
		http.Error(
			w,
			"輸入 json 格式錯誤",
			http.StatusBadRequest,
		)
		return
	}

	for i := range tasks {
		if tasks[i].ID == id {
			tasks[i].Title = input.Title
			tasks[i].Completed = input.Completed

			w.Header().Set(
				"Content-Type",
				"application/json",
			)

			err := json.NewEncoder(w).Encode(tasks[i])
			if err != nil {
				log.Printf("無法正確寫入 json %v", err)
			}
			return
		}
	}

	http.Error(
		w,
		"找不到目標 task",
		http.StatusBadRequest,
	)

}

func createTaskHandler(
	w http.ResponseWriter,
	r *http.Request,
) {
	var input CreateTaskRequest

	err := json.NewDecoder(r.Body).Decode(&input)

	if err != nil {
		http.Error(w,
			"無效 json",
			http.StatusBadRequest)

		return
	}

	task := Task{
		ID:        nextID,
		Title:     input.Title,
		Completed: false,
	}

	nextID++

	tasks = append(tasks, task)

	w.Header().Set(
		"Content-Type",
		"application/json",
	)

	w.WriteHeader(
		http.StatusCreated,
	)

	err = json.NewEncoder(w).Encode(task)
	if err != nil {
		log.Printf("JSON 編碼失敗: %v", err)
	}

}

func getTasksHandler(
	w http.ResponseWriter,
	r *http.Request,
) {

	w.Header().Set(
		"Content-Type",
		"application/json",
	)

	err := json.NewEncoder(w).Encode(tasks)
	if err != nil {
		log.Printf("JSON 編碼失敗: %v", err)
	}
}

func main() {

	databaseURL :=
		"postgres://task_user:task_password@127.0.0.1:5433/task_db?sslmode=disable"

	var err error

	dbPool, err = pgxpool.New(
		context.Background(),
		databaseURL,
	)

	if err != nil {
		log.Fatal(
			"建立資料庫 connection pool 失敗: ",
			err,
		)
	}

	err = dbPool.Ping(
		context.Background(),
	)

	if err != nil {
		dbPool.Close()

		log.Fatal(
			"PostgreSQL 連線失敗: ",
			err,
		)
	}

	defer dbPool.Close()

	http.HandleFunc(
		"GET /db-task",
		getDbTasksHandler,
	)

	http.HandleFunc(
		"POST /db-task",
		createDbTaskHandler,
	)

	http.HandleFunc(
		"PUT /db-task/{id}",
		updateDbTaskHandler,
	)

	http.HandleFunc(
		"DELETE /db-task/{id}",
		deleteDbTaskHandler,
	)

	http.HandleFunc(
		"GET /tasks",
		getTasksHandler,
	)

	http.HandleFunc(
		"POST /tasks",
		createTaskHandler,
	)

	http.HandleFunc(
		"PUT /tasks/{id}",
		updateTaskHandler,
	)
	http.HandleFunc(
		"DELETE /tasks/{id}",
		deleteTaskHandler,
	)

	var mux = http.DefaultServeMux

	log.Println("localhost:8080 上已啟動")

	err = http.ListenAndServe(":8080", mux)
	if err != nil {
		log.Fatal(err)
	}

}
