package main

import (
	"log"

	"github.com/Makhaev/projectname/pkg/db"
	"github.com/Makhaev/projectname/pkg/server"
)

func main() {

	dbFile := "scheduler.db"
	err := db.InitDB(dbFile)
	if err != nil {
		log.Fatal("Ошибка инициализации базы:", err)
	}

	// Запуск сервера
	err = server.Run()
	if err != nil {
		log.Fatal("Ошибка запуска сервера:", err)
	}
}
