package db

import (
	"database/sql"
	"fmt"
	"os"

	_ "modernc.org/sqlite"
)

var db *sql.DB

func InitDB(dbFile string) error {
	_, err := os.Stat(dbFile)

	var install bool
	if err != nil {
		install = true
	}

	dbLocal, err := sql.Open("sqlite", dbFile)
	if err != nil {
		return fmt.Errorf("ошибка открытия базы: %w", err)
	}

	if err := dbLocal.Ping(); err != nil {
		return fmt.Errorf("ошибка подключения к базе: %w", err)
	}

	db = dbLocal

	if install {
		creating := `
		CREATE TABLE scheduler (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			date TEXT NOT NULL,
			title TEXT NOT NULL,
			comment TEXT,
			repeat TEXT CHECK(length(repeat) <= 128)
		);
		CREATE INDEX idx_scheduler_date ON scheduler(date);
		`
		_, err = db.Exec(creating)
		if err != nil {
			return fmt.Errorf("ошибка создания таблицы: %w", err)
		}
	}

	return nil
}

func AddTask(task *Task) (int64, error) {
	query := `
		INSERT INTO scheduler (date, title, comment, repeat)
		VALUES (?, ?, ?, ?)
	`
	res, err := db.Exec(query, task.Date, task.Title, task.Comment, task.Repeat)
	if err != nil {
		return 0, err
	}
	return res.LastInsertId()
}

func GetTask(id string) (*Task, error) {
	query := "SELECT id, date, title, comment, repeat FROM scheduler WHERE id = ?"
	var task Task

	err := db.QueryRow(query, id).Scan(&task.ID, &task.Date, &task.Title, &task.Comment, &task.Repeat)
	if err != nil {
		return nil, fmt.Errorf("задача не найдена")
	}

	return &task, nil
}
