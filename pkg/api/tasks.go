package api

import (
	"net/http"

	"github.com/Makhaev/projectname/pkg/db"
)

type TasksResp struct {
	Tasks []*db.Task `json:"tasks"`
}

func tasksHandler(w http.ResponseWriter, r *http.Request) {
	tasks, err := db.Tasks(50)
	if err != nil {
		writeJSON(w, map[string]string{
			"error": "не удалось получить задачи",
		})
		return
	}

	writeJSON(w, TasksResp{
		Tasks: tasks,
	})
}
