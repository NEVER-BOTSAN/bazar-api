package server

import (
	"encoding/json"
	"fmt"
	"net/http"
)

func (s *Server) HandlerShopCategories(w http.ResponseWriter, r *http.Request) {
	// Получаем связи между магазинами и категориями
	shopCategories, err := s.App.GetShopCategories()
	if err != nil {
		fmt.Println(err.Error())
		http.Error(w, "Ошибка при получении категорий магазинов: "+err.Error(), http.StatusInternalServerError)
		return // Не забываем выходить из функции при ошибке
	}

	// Преобразуем shopCategories в JSON
	respjson, err := json.Marshal(shopCategories)
	if err != nil {
		fmt.Println(err.Error())
		http.Error(w, "Ошибка при преобразовании данных в JSON: "+err.Error(), http.StatusInternalServerError)
		return // Не забываем выходить из функции при ошибке
	}

	// Устанавливаем заголовок Content-Type и возвращаем статус и данные
	w.Header().Set("Content-Type", "application/json") // Устанавливаем заголовок, чтобы указать, что это JSON
	w.WriteHeader(http.StatusOK)                       // Устанавливаем статус 200 OK
	w.Write(respjson)                                  // Отправляем данные клиенту
}
