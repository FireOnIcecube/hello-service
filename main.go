package main

import (
	"context"
	"encoding/json"
	"log"
	"net/http"
	"os"
	"strconv"

	"example.com/hello-service/handlers"
	"example.com/hello-service/models"
	"github.com/jackc/pgx/v5/pgxpool"
)

var tasks = []models.Task{
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

	task := models.Task{
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

	databaseURL, ok :=
		os.LookupEnv("DATABASE_URL")

	if !ok || databaseURL == "" {
		log.Fatal(
			"DATABASE_URL 環境變數未設定",
		)
	}

	var err error

	dbPool, err := pgxpool.New(
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

	handler := handlers.New(dbPool)

	defer dbPool.Close()

	http.HandleFunc(
		"GET /db-task",
		handler.GetDbTasks,
	)

	http.HandleFunc(
		"POST /db-task",
		handler.CreateDbTask,
	)

	http.HandleFunc(
		"PUT /db-task/{id}",
		handler.UpdateDbTask,
	)

	http.HandleFunc(
		"DELETE /db-task/{id}",
		handler.DeleteDbTask,
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
		log.Printf(
			"HTTP Server 發生錯誤: %v",
			err,
		)
		return
	}

}
