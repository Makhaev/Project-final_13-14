package api

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/Makhaev/projectname/pkg/db"
)

func addTaskHandler(w http.ResponseWriter, r *http.Request) {
	var task db.Task

	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&task); err != nil {
		writeJSON(w, map[string]string{"error": "Ошибка чтения JSON"})
		return
	}

	// Проверка обязательного поля title
	if task.Title == "" {
		writeJSON(w, map[string]string{"error": "Не указан заголовок задачи"})
		return
	}

	// Проверка и исправление даты
	if err := checkDate(&task); err != nil {
		writeJSON(w, map[string]string{"error": err.Error()})
		return
	}

	id, err := db.AddTask(&task)
	if err != nil {
		writeJSON(w, map[string]string{"error": err.Error()})
		return
	}

	writeJSON(w, map[string]interface{}{"id": id})
}

func checkDate(task *db.Task) error {
	now := time.Now()

	if task.Date == "" {
		task.Date = now.Format(DateFormat)
	}

	t, err := time.Parse(DateFormat, task.Date)
	if err != nil {
		return err
	}

	// Если правило указано — проверяем его и пересчитываем дату при необходимости
	if task.Repeat != "" {
		next, err := NextDate(now, task.Date, task.Repeat)
		if err != nil {
			return err
		}
		if afterNow(now, t) {
			task.Date = next
		}
	} else {
		if afterNow(now, t) {
			task.Date = now.Format(DateFormat)
		}
	}

	return nil
}

func afterNow(now, t time.Time) bool {
	return now.After(t)
}

func writeJSON(w http.ResponseWriter, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(data)
}
