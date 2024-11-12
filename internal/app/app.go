package app

import (
	"database/sql"
	"fmt"
	"log"

	_ "github.com/lib/pq"
)

type App struct {
	db *sql.DB
}

func NewApp(connStr string) *App {

	// Открываем соединение с базой данных
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		log.Fatal("Ошибка открытия соединения:", err)
	}

	// Проверяем подключение к базе данных
	err = db.Ping()
	if err != nil {
		log.Fatal("Не удалось подключиться к базе данных:", err)
	}
	fmt.Println("Успешное подключение к базе данных PostgreSQL!")

	app := App{

		db: db,
	}
	return &app

}
