package app

import (
	"fmt"
	"log"
)

type Category struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
}

// Метод для создания таблицы categories
func (app *App) CreateCategoryTable() {
	query := `
	CREATE TABLE IF NOT EXISTS categories (
		id SERIAL PRIMARY KEY,
		name TEXT NOT NULL
	);`

	_, err := app.db.Exec(query)
	if err != nil {
		log.Fatal("Ошибка при создании таблицы categories:", err)
	}

	fmt.Println("Таблица categories успешно создана (если её не было)!")
}

// Метод для добавления одной записи в таблицу categories
func (app *App) InsertCategory(name string) {
	query := `INSERT INTO categories (name) VALUES ($1)`
	_, err := app.db.Exec(query, name)
	if err != nil {
		log.Fatal("Ошибка при вставке данных в таблицу categories:", err)
	}

	fmt.Println("Категория успешно добавлена!")
}

// Метод для добавления нескольких категорий
func (app *App) InsertSampleCategories() {
	app.InsertCategory("Категория 1")
	app.InsertCategory("Категория 2")
	app.InsertCategory("Категория 3")
	app.InsertCategory("Категория 4")
	app.InsertCategory("Категория 5")
	app.InsertCategory("Категория 6")
}

// Метод для получения всех категорий из таблицы
func (app *App) GetCategories() ([]Category, error) {
	query := `SELECT id, name FROM categories;`

	rows, err := app.db.Query(query)
	if err != nil {
		return nil, fmt.Errorf("ошибка при получении данных из таблицы categories: %v", err)
	}
	defer rows.Close()

	var categories []Category
	for rows.Next() {
		var category Category
		err := rows.Scan(&category.ID, &category.Name)
		if err != nil {
			return nil, fmt.Errorf("ошибка при сканировании данных из таблицы categories: %v", err)
		}
		categories = append(categories, category)
	}

	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("ошибка во время обработки строк: %v", err)
	}

	return categories, nil
}
