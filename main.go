package main

import (
	"encoding/json"
	"log"
	"net/http"
	"strconv"
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

	err := http.ListenAndServe(":8080", mux)
	if err != nil {
		log.Fatal(err)
	}

}
