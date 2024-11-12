package app

import (
	"database/sql"
	"fmt"
	"log"
)

type Shop struct {
	ID          int    `json:"id"`
	Name        string `json:"name"`
	Image       string `json:"image"`
	Price       int    `json:"price"`
	Description string `json:"description"`
}
type ShopWithCategories struct {
	Shop        Shop     `json:"shop"`
	CategoryIDs []string `json:"categories"`
}

func (app *App) GetShops(limit, offset int) ([]ShopWithCategories, error) {
	// SQL-запрос для получения данных о магазинах и категориях
	query := `
        SELECT s.id, s.name, s.image, s.price, s.description, c.name AS category_name
        FROM shops s
        LEFT JOIN shop_categories sc ON s.id = sc.shop_id
        LEFT JOIN categories c ON sc.category_id = c.id
        LIMIT $1 OFFSET $2`

	rows, err := app.db.Query(query, limit, offset)
	if err != nil {
		fmt.Println("Ошибка запроса к базе данных:", err)
		return nil, err
	}
	defer rows.Close()

	var result []ShopWithCategories

	// Проходим по всем строкам результата запроса
	for rows.Next() {
		var shop Shop
		var categoryName sql.NullString

		// Сканируем данные о магазине
		if err := rows.Scan(&shop.ID, &shop.Name, &shop.Image, &shop.Price, &shop.Description, &categoryName); err != nil {
			fmt.Println("Ошибка сканирования данных:", err)
			return nil, err
		}

		// Формируем структуру для магазина с категориями
		shopWithCategories := ShopWithCategories{
			Shop:        shop,       // Заполняем информацию о магазине
			CategoryIDs: []string{}, // Инициализируем список категорий
		}

		// Если категория присутствует, добавляем её в список категорий
		if categoryName.Valid {
			shopWithCategories.CategoryIDs = append(shopWithCategories.CategoryIDs, categoryName.String)
		}

		// Проверяем, если магазин уже добавлен, добавляем категорию
		found := false
		for i, s := range result {
			if s.Shop.ID == shopWithCategories.Shop.ID {
				result[i].CategoryIDs = append(result[i].CategoryIDs, categoryName.String)
				found = true
				break
			}
		}

		// Если магазин не найден, добавляем его в список
		if !found {
			result = append(result, shopWithCategories)
		}
	}

	return result, nil
}

func (app *App) CreateTable() {
	query := `
	CREATE TABLE IF NOT EXISTS shops (
		id SERIAL PRIMARY KEY,
		name TEXT NOT NULL,
		image TEXT NOT NULL,
		price INTEGER NOT NULL,
		description TEXT NOT NULL
	);`

	_, err := app.db.Exec(query)
	if err != nil {
		log.Fatal("Ошибка при создании таблицы:", err)
	}

	fmt.Println("Таблица shops успешно создана (если её не было)!")
}

// Функция для добавления одной записи в таблицу shops
func (app *App) InsertShop(name, image string, price int, description string) {
	query := `INSERT INTO shops (name, image, price, description) VALUES ($1, $2, $3, $4)`
	_, err := app.db.Exec(query, name, image, price, description)
	if err != nil {
		log.Fatal("Ошибка при вставке данных:", err)
	}

	fmt.Println("Данные успешно добавлены!")
}

// Функция для добавления нескольких записей
func (app *App) CreateTableWithShops() {
	app.InsertShop("Магазин 1", "картинка1.jpeg", 100, "Описание 1")
	app.InsertShop("Магазин 2", "картинка2.jpeg", 200, "Описание 2")
	app.InsertShop("Магазин 3", "картинка3.jpeg", 300, "Описание 3")
	app.InsertShop("Магазин 4", "картинка4.jpeg", 400, "Описание 4")
	app.InsertShop("Магазин 5", "картинка5.jpeg", 500, "Описание 5")
	app.InsertShop("Магазин 6", "картинка6.jpeg", 600, "Описание 6")
}
func (app *App) CreateNewShop(shop Shop) (int, error) {
	query := `INSERT INTO shops (name, image, price, description) VALUES ($1, $2, $3, $4) RETURNING id`
	var shopID int
	err := app.db.QueryRow(query, shop.Name, shop.Image, shop.Price, shop.Description).Scan(&shopID)
	if err != nil {
		return 0, fmt.Errorf("ошибка при добавлении нового магазина: %v", err)
	}
	fmt.Println("Новый магазин успешно добавлен с ID:", shopID)
	return shopID, nil
}

func (app *App) DeleteShopByID(id string) error {
	query := `DELETE FROM shops WHERE id = $1`

	_, err := app.db.Exec(query, id)
	if err != nil {
		return fmt.Errorf("ошибка при удалении магазина: %v", err)
	}
	return nil
}

func (app *App) UpdateShopByID(id string, updatedShop Shop) error {
	query := `
		UPDATE shops 
		SET name = $1, image = $2, price = $3, description = $4 
		WHERE id = $5`

	_, err := app.db.Exec(query, updatedShop.Name, updatedShop.Image, updatedShop.Price, updatedShop.Description, id)
	if err != nil {
		return fmt.Errorf("ошибка при обновлении магазина: %v", err)
	}
	return nil
}
func (app *App) UpdateShopFields(id string, fields map[string]interface{}) error {
	query := "UPDATE shops SET "
	args := []interface{}{}
	i := 1

	// Динамически создаем часть запроса для каждого переданного поля
	for field, value := range fields {
		query += fmt.Sprintf("%s = $%d, ", field, i)
		args = append(args, value)
		i++
	}
	query = query[:len(query)-2] // Удаляем последнюю запятую и пробел
	query += " WHERE id = $" + fmt.Sprintf("%d", i)
	args = append(args, id)

	// Выполняем запрос
	_, err := app.db.Exec(query, args...)
	return err
}
func (app *App) GetShopsByCategoryID(categoryID string, limit, offset int) ([]Shop, error) {
	// Формируем SQL-запрос для получения магазинов по категории с LIMIT и OFFSET
	query := `
	SELECT s.id, s.name, s.image, s.price, s.description
	FROM shops s
	JOIN shop_categories sc ON s.id = sc.shop_id
	WHERE sc.category_id = $1
	LIMIT $2 OFFSET $3`

	rows, err := app.db.Query(query, categoryID, limit, offset)
	if err != nil {
		return nil, fmt.Errorf("ошибка при выполнении запроса: %v", err)
	}
	defer rows.Close()

	// Собираем данные о магазинах
	var shops []Shop
	for rows.Next() {
		var shop Shop
		if err := rows.Scan(&shop.ID, &shop.Name, &shop.Image, &shop.Price, &shop.Description); err != nil {
			return nil, fmt.Errorf("ошибка сканирования данных: %v", err)
		}
		shops = append(shops, shop)
	}

	return shops, nil
}
func (app *App) AddShopCategories(shopID int, categoryIDs []int) error {
	query := `INSERT INTO shop_categories (shop_id, category_id) VALUES ($1, $2)`

	// Открываем транзакцию
	tx, err := app.db.Begin()
	if err != nil {
		return fmt.Errorf("ошибка при начале транзакции: %v", err)
	}
	defer tx.Rollback()

	// Добавляем связи между магазином и категориями
	for _, categoryID := range categoryIDs {
		_, err := tx.Exec(query, shopID, categoryID)
		if err != nil {
			return fmt.Errorf("ошибка при добавлении категории %d для магазина %d: %v", categoryID, shopID, err)
		}
	}

	// Подтверждаем транзакцию
	if err := tx.Commit(); err != nil {
		return fmt.Errorf("ошибка при подтверждении транзакции: %v", err)
	}

	fmt.Println("Категории успешно добавлены для магазина с ID:", shopID)
	return nil
}

func (app *App) UpdateShopCategories(shopID string, categoryIDs []int) error {
	// Удаляем старые категории для данного магазина
	_, err := app.db.Exec("DELETE FROM shop_categories WHERE shop_id = $1", shopID)
	if err != nil {
		return fmt.Errorf("не удалось удалить старые категории: %v", err)
	}

	// Добавляем новые категории
	for _, categoryID := range categoryIDs {
		_, err := app.db.Exec("INSERT INTO shop_categories (shop_id, category_id) VALUES ($1, $2)", shopID, categoryID)
		if err != nil {
			return fmt.Errorf("не удалось добавить категорию с ID %d: %v", categoryID, err)
		}
	}

	return nil
}
