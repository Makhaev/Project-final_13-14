package api

import (
	"net/http"
	"time"

	"github.com/Makhaev/projectname/pkg/db"
)

func taskDoneHandler(w http.ResponseWriter, r *http.Request) {
	id := r.URL.Query().Get("id")
	if id == "" {
		writeJSON(w, map[string]string{"error": "Не указан идентификатор"})
		return
	}

	task, err := db.GetTask(id)
	if err != nil {
		writeJSON(w, map[string]string{"error": "Задача не найдена"})
		return
	}

	const layout = "20060102"

	if task.Repeat == "" {
		// одноразовая — удаляем
		err = db.DeleteTask(id)
		if err != nil {
			writeJSON(w, map[string]string{"error": err.Error()})
			return
		}
	} else {
		// преобразуем дату в time.Time
		t, err := time.Parse(layout, task.Date)
		if err != nil {
			writeJSON(w, map[string]string{"error": "Неверный формат даты"})
			return
		}

		// получаем следующую дату
		next, err := NextDate(t, task.Repeat, layout)
		if err != nil {
			writeJSON(w, map[string]string{"error": err.Error()})
			return
		}

		err = db.UpdateDate(next, id)
		if err != nil {
			writeJSON(w, map[string]string{"error": err.Error()})
			return
		}
	}

	writeJSON(w, map[string]string{})
}
