package api

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/Makhaev/projectname/pkg/auth"
	"github.com/Makhaev/projectname/pkg/db"
)

// Функция инициализации API
func Init() {
	http.HandleFunc("/api/nextdate", nextDateHandler)
	http.HandleFunc("/api/task", auth.Auth(taskHandler))
	http.HandleFunc("/api/tasks", auth.Auth(tasksHandler))
	http.HandleFunc("/api/task/done", auth.Auth(taskDoneHandler))
	http.HandleFunc("/api/signin", signinHandler)
}

func taskHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	// обработка других методов будет добавлена на следующих шагах
	case http.MethodPost:
		addTaskHandler(w, r)
	case http.MethodGet:
		handleGetTask(w, r)
	case http.MethodPut:
		handleUpdateTask(w, r)
	case http.MethodDelete:
		id := r.URL.Query().Get("id")
		if id == "" {
			writeJSON(w, map[string]string{"error": "Не указан идентификатор"})
			return
		}
		err := db.DeleteTask(id)
		if err != nil {
			writeJSON(w, map[string]string{"error": err.Error()})
			return
		}
		writeJSON(w, map[string]string{})

	}
}

func nextDateHandler(w http.ResponseWriter, r *http.Request) {
	nowStr := r.FormValue("now")
	dateStr := r.FormValue("date")
	repeat := r.FormValue("repeat")

	var now time.Time
	var err error
	if nowStr == "" {
		now = time.Now()
	} else {
		now, err = time.Parse(DateFormat, nowStr)
		if err != nil {
			http.Error(w, "неверный формат now", http.StatusBadRequest)
			return
		}
	}

	result, err := NextDate(now, dateStr, repeat)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// JSON-ответ
	w.Header().Set("Content-Type", "text/plain")
	w.Write([]byte(result))
}

func handleGetTask(w http.ResponseWriter, r *http.Request) {
	id := r.URL.Query().Get("id")
	if id == "" {
		writeJSON(w, map[string]string{
			"error": "Не указан идентификатор",
		})
		return
	}

	task, err := db.GetTask(id)
	if err != nil {
		writeJSON(w, map[string]string{
			"error": "Задача не найдена",
		})
		return
	}

	writeJSON(w, task)
}

func handleUpdateTask(w http.ResponseWriter, r *http.Request) {
	var task db.Task
	err := json.NewDecoder(r.Body).Decode(&task)
	if err != nil {
		writeJSON(w, map[string]string{
			"error": "Невозможно прочитать тело запроса",
		})
		return
	}

	// Простая проверка
	if task.ID == "" || task.Date == "" || task.Title == "" {
		writeJSON(w, map[string]string{
			"error": "Неверные параметры задачи",
		})
		return
	}

	err = db.UpdateTask(&task)
	if err != nil {
		writeJSON(w, map[string]string{
			"error": err.Error(),
		})
		return
	}

	writeJSON(w, map[string]string{})
}
