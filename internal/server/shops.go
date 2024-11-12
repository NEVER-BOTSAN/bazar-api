package server

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strconv"
	"test-server/internal/app"
)

type ShopRequest struct {
	Shop        app.Shop `json:"shop"`
	CategoryIDs []int    `json:"categories"`
}

func (s *Server) HandlerShops(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		s.GetHandlerShops(w, r)
	case http.MethodPost:
		s.PostHandlerShops(w, r)
	case http.MethodPut:
		s.PutHandlerShops(w, r)
	case http.MethodDelete:
		s.DeleteHandlerShops(w, r)
	case http.MethodPatch:
		s.PatchHandlerShops(w, r)
	default:
		http.Error(w, "Метод не доступен", http.StatusMethodNotAllowed)
	}

}

func (s *Server) GetHandlerShops(w http.ResponseWriter, r *http.Request) {
	// Получаем параметры page, limit и category_id из запроса
	pageStr := r.URL.Query().Get("page")
	limitStr := r.URL.Query().Get("limit")
	categoryID := r.URL.Query().Get("category_id")

	// Устанавливаем значения по умолчанию
	page := 1
	limit := 10

	// Преобразуем page и limit в целые числа

	p, err := strconv.Atoi(pageStr)
	if err == nil && p > 0 {
		page = p
	}

	l, err := strconv.Atoi(limitStr)
	if err == nil && l > 0 {
		limit = l
	}

	// Рассчитываем offset для SQL-запроса
	offset := (page - 1) * limit

	var shops []app.Shop
	var shopWithCategories []app.ShopWithCategories

	// Если categoryID не пустой, используем фильтрацию по категории
	if categoryID != "" {
		shops, err = s.App.GetShopsByCategoryID(categoryID, limit, offset)
		if err != nil {
			fmt.Println(err.Error())
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		// Отправляем результат в формате JSON
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(shops)
	} else {
		// Если categoryID пустой, получаем все магазины
		shopWithCategories, err = s.App.GetShops(limit, offset)
		if err != nil {
			fmt.Println(err.Error())
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		// Отправляем результат в формате JSON
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(shopWithCategories)
	}

}

func (s *Server) PostHandlerShops(w http.ResponseWriter, r *http.Request) {
	// Структура для обработки данных запроса

	var reqBody ShopRequest
	var shopData app.Shop

	// Чтение и сохранение тела запроса в переменной для многократного использования
	bodyBytes, err := io.ReadAll(r.Body)
	if err != nil {
		fmt.Println("Ошибка чтения тела запроса:", err.Error())
		http.Error(w, "Ошибка чтения тела запроса", http.StatusBadRequest)
		return
	}
	// Восстанавливаем тело запроса для повторного использования
	r.Body.Close()
	r.Body = io.NopCloser(bytes.NewReader(bodyBytes))

	// Попытка декодировать JSON как ShopRequest
	err = json.Unmarshal(bodyBytes, &reqBody)
	if err != nil || reqBody.Shop.Name == "" {
		// Если ошибка декодирования или имя магазина пустое, пробуем декодировать как app.Shop
		if err := json.Unmarshal(bodyBytes, &shopData); err != nil {
			fmt.Println("Ошибка декодирования JSON:", err.Error())
			http.Error(w, "Ошибка декодирования JSON", http.StatusBadRequest)
			return
		}
		reqBody.Shop = shopData // Присваиваем декодированные данные в reqBody.Shop
	}

	// Сохраняем магазин в базу данных
	shopID, err := s.App.CreateNewShop(reqBody.Shop)
	if err != nil {
		fmt.Println("Ошибка при добавлении магазина:", err.Error())
		http.Error(w, "Ошибка при добавлении магазина", http.StatusInternalServerError)
		return
	}

	// Добавляем связи между магазином и категориями
	if len(reqBody.CategoryIDs) > 0 {
		if err := s.App.AddShopCategories(shopID, reqBody.CategoryIDs); err != nil {
			fmt.Println("Ошибка при добавлении категорий:", err.Error())
			http.Error(w, "Ошибка при добавлении категорий", http.StatusInternalServerError)
			return
		}
	}

	fmt.Fprintln(w, "Магазин успешно добавлен!")
}

func (s *Server) DeleteHandlerShops(w http.ResponseWriter, r *http.Request) {
	// Получаем параметр id из URL
	query := r.URL.Query()
	id := query.Get("id")

	if id == "" {
		http.Error(w, "ID не указан", http.StatusBadRequest)
		return
	}

	// Вызываем метод для удаления магазина
	err := s.App.DeleteShopByID(id)
	if err != nil {
		http.Error(w, fmt.Sprintf("Ошибка при удалении магазина: %v", err), http.StatusInternalServerError)
		return
	}

	// Если удаление прошло успешно
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Магазин успешно удален"))

}

func (s *Server) PutHandlerShops(w http.ResponseWriter, r *http.Request) {
	// Получаем параметр id из URL
	query := r.URL.Query()
	id := query.Get("id")

	if id == "" {
		http.Error(w, "ID не указан", http.StatusBadRequest)
		return
	}

	// Декодируем тело запроса в структуру ShopUpdateRequest
	var request ShopRequest
	err := json.NewDecoder(r.Body).Decode(&request)
	if err != nil {
		http.Error(w, fmt.Sprintf("Ошибка декодирования данных: %v", err), http.StatusBadRequest)
		return
	}

	// Обновляем магазин в базе данных
	err = s.App.UpdateShopByID(id, request.Shop)
	if err != nil {
		http.Error(w, fmt.Sprintf("Ошибка при обновлении магазина: %v", err), http.StatusInternalServerError)
		return
	}

	// Обновляем категории магазина, если они указаны

	err = s.App.UpdateShopCategories(id, request.CategoryIDs)
	if err != nil {
		http.Error(w, fmt.Sprintf("Ошибка при обновлении категорий: %v", err), http.StatusInternalServerError)
		return
	}

	// Если обновление прошло успешно
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Магазин и категории успешно обновлены"))
}
func (s *Server) PatchHandlerShops(w http.ResponseWriter, r *http.Request) {
	// Получаем `id` магазина из URL-параметров
	query := r.URL.Query()
	id := query.Get("id")
	if id == "" {
		http.Error(w, "ID магазина не указан", http.StatusBadRequest)
		return
	}

	// Декодируем тело запроса
	var reqBody map[string]interface{}
	err := json.NewDecoder(r.Body).Decode(&reqBody)
	if err != nil {
		http.Error(w, "Ошибка декодирования данных: "+err.Error(), http.StatusBadRequest)
		return
	}

	// Обновляем поля магазина, если они присутствуют
	if shopFields, ok := reqBody["shop"].(map[string]interface{}); ok {
		err = s.App.UpdateShopFields(id, shopFields)
		if err != nil {
			http.Error(w, "Ошибка обновления данных магазина: "+err.Error(), http.StatusInternalServerError)
			return
		}
	}

	// Обновляем категории, если они присутствуют
	if categoryIDs, ok := reqBody["categories"].([]interface{}); ok {
		// Преобразуем интерфейсы в слайс целых чисел
		var categories []int
		for _, categoryID := range categoryIDs {
			if idFloat, ok := categoryID.(float64); ok {
				categories = append(categories, int(idFloat)) // Преобразуем float64 в int
			}
		}

		err = s.App.UpdateShopCategories(id, categories)
		if err != nil {
			http.Error(w, "Ошибка обновления категорий магазина: "+err.Error(), http.StatusInternalServerError)
			return
		}
	}

	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, "Магазин с ID %s успешно обновлен", id)
}
