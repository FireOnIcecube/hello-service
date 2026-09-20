package main

import (
	"net/http"

	"example.com/hello-service/handlers"
)

func newRouter(
	handler *handlers.Handler,
) *http.ServeMux {

	var mux = http.NewServeMux()

	mux.HandleFunc(
		"GET /db-task",
		handler.GetDbTasks,
	)

	mux.HandleFunc(
		"POST /db-task",
		handler.CreateDbTask,
	)

	mux.HandleFunc(
		"PUT /db-task/{id}",
		handler.UpdateDbTask,
	)

	mux.HandleFunc(
		"DELETE /db-task/{id}",
		handler.DeleteDbTask,
	)

	mux.HandleFunc(
		"GET /tasks",
		getTasksHandler,
	)

	mux.HandleFunc(
		"POST /tasks",
		createTaskHandler,
	)

	mux.HandleFunc(
		"PUT /tasks/{id}",
		updateTaskHandler,
	)
	mux.HandleFunc(
		"DELETE /tasks/{id}",
		deleteTaskHandler,
	)

	return mux

}
