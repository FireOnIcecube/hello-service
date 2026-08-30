package handlers

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"example.com/hello-service/models"
)

func (h *Handler) GetDbTasks(
	w http.ResponseWriter,
	r *http.Request,
) {

	rows, err := h.dbPool.Query(
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
