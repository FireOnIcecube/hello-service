package handlers

import (
	"encoding/json"
	"errors"
	"log"
	"net/http"
	"strconv"

	"example.com/hello-service/repositories"
)

func (h *Handler) GetDbTasks(
	w http.ResponseWriter,
	r *http.Request,
) {

	tasks, err := h.taskRepository.GetAll(r.Context())

	if err != nil {

		log.Printf(
			"repositories 調用失敗: %v",
			err,
		)

		http.Error(
			w,
			"伺服器內部錯誤",
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

func (h *Handler) CreateDbTask(
	w http.ResponseWriter,
	r *http.Request,

) {

	var input createTaskRequest

	err := json.NewDecoder(r.Body).Decode(&input)

	if err != nil {
		http.Error(
			w,
			"傳遞參數格式不正確",
			http.StatusBadRequest,
		)
		return
	}

	task, err := h.taskRepository.Create(r.Context(), input.Title)

	if err != nil {

		log.Printf(
			"建立 task 失敗: %v",
			err,
		)

		http.Error(
			w,
			"伺服器內部錯誤",
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

	if err != nil {
		log.Printf("json 編碼錯誤:%v \n", err)
	}

}

func (h *Handler) UpdateDbTask(
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

	var input updateTaskRequest

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

	task, err := h.taskRepository.Update(r.Context(), id, input.Title, input.Completed)

	if err != nil {

		if errors.Is(err, repositories.ErrTaskNotFound) {
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

	err = json.NewEncoder(w).Encode(task)

	if err != nil {
		log.Printf("json 編碼錯誤:%v \n", err)
	}

}

func (h *Handler) DeleteDbTask(
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

	err = h.taskRepository.Delete(r.Context(), id)

	if err != nil {
		if errors.Is(err, repositories.ErrTaskNotFound) {
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
