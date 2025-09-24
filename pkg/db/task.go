package db

import (
	"fmt"
)

type Task struct {
	ID      string `json:"id"`
	Date    string `json:"date"`
	Title   string `json:"title"`
	Comment string `json:"comment"`
	Repeat  string `json:"repeat"`
}

func Tasks(limit int) ([]*Task, error) {
	query := "SELECT id, date, title, comment, repeat FROM scheduler ORDER BY date ASC LIMIT ?"

	rows, err := db.Query(query, limit)
	if err != nil {
		return nil, fmt.Errorf("ошибка при получении задач из базы данных: %w", err)
	}
	defer rows.Close()

	var tasks []*Task
	for rows.Next() {
		task := &Task{}
		err := rows.Scan(&task.ID, &task.Date, &task.Title, &task.Comment, &task.Repeat)
		if err != nil {
			return nil, fmt.Errorf("ошибка при сканировании данных: %w", err)
		}
		tasks = append(tasks, task)
	}

	// Гарантируем, что не вернём nil
	if tasks == nil {
		tasks = []*Task{}
	}

	return tasks, nil
}

func UpdateTask(task *Task) error {
	// параметры пропущены, не забудьте указать WHERE
	query := `
	UPDATE scheduler
	SET date = ?, title = ?, comment = ?, repeat = ?
	WHERE id = ?
	`
	res, err := db.Exec(query, task.Date, task.Title, task.Comment, task.Repeat, task.ID)
	if err != nil {
		return err
	}
	// метод RowsAffected() возвращает количество записей к которым
	// был применена SQL команда
	count, err := res.RowsAffected()
	if err != nil {
		return err
	}
	if count == 0 {
		return fmt.Errorf(`incorrect id for updating task`)
	}
	return nil
}

func DeleteTask(id string) error {
	_, err := db.Exec("DELETE FROM scheduler WHERE id = ?", id)
	return err
}

func UpdateDate(nextDate string, id string) error {
	_, err := db.Exec("UPDATE scheduler SET date = ? WHERE id = ?", nextDate, id)
	return err
}
