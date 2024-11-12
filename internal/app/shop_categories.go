package app

import (
	"fmt"
	"log"
)

type ShopCategory struct {
	ShopID     int `json:"shop_id"`
	CategoryID int `json:"category_id"`
}

func (app *App) CreateShopCategoryTable() {
	query := `
	CREATE TABLE IF NOT EXISTS shop_categories (
		shop_id INT NOT NULL,
		category_id INT NOT NULL,
		PRIMARY KEY (shop_id, category_id),
		FOREIGN KEY (shop_id) REFERENCES shops(id) ON DELETE CASCADE,
		FOREIGN KEY (category_id) REFERENCES categories(id) ON DELETE CASCADE
	);`

	_, err := app.db.Exec(query)
	if err != nil {
		log.Fatal("Ошибка при создании таблицы shop_categories:", err)
	}

	fmt.Println("Таблица shop_categories успешно создана (если её не было)!")
}

func (app *App) InsertShopCategory(shopID, categoryID int) {
	query := `INSERT INTO shop_categories (shop_id, category_id) VALUES ($1, $2)`
	_, err := app.db.Exec(query, shopID, categoryID)
	if err != nil {
		log.Fatal("Ошибка при вставке данных в таблицу shop_categories:", err)
	}

	fmt.Println("Связь между магазином и категорией успешно добавлена!")
}

func (app *App) GetShopCategories() ([]ShopCategory, error) {
	query := `SELECT shop_id, category_id FROM shop_categories;`

	rows, err := app.db.Query(query)
	if err != nil {
		return nil, fmt.Errorf("ошибка при получении данных из таблицы shop_categories: %v", err)
	}
	defer rows.Close()

	var shopCategories []ShopCategory
	for rows.Next() {
		var shopCategory ShopCategory
		err := rows.Scan(&shopCategory.ShopID, &shopCategory.CategoryID)
		if err != nil {
			return nil, fmt.Errorf("ошибка при сканировании данных из таблицы shop_categories: %v", err)
		}
		shopCategories = append(shopCategories, shopCategory)
	}

	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("ошибка во время обработки строк: %v", err)
	}

	return shopCategories, nil
}
