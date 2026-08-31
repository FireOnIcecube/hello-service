package handlers

type createTaskRequest struct {
	Title string `json:"title"`
}

type updateTaskRequest struct {
	Title     string `json:"title"`
	Completed bool   `json:"completed"`
}
