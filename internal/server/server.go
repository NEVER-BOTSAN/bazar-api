package server

import (
	"fmt"
	"net/http"
	"test-server/internal/app"
)

type Server struct {
	http.Server
	App *app.App
}

func New(serviceApp *app.App) *Server {
	var srv Server
	srv.Addr = ":8080"
	srv.App = serviceApp
	return &srv
}

func (s *Server) Run() {
	fmt.Println("Сервер запущен")
	s.Handler = s.InitRoutes()
	if err := s.ListenAndServe(); err != nil {
		fmt.Println(err.Error())
	}

}

func (s *Server) InitRoutes() http.Handler {
	mux := http.NewServeMux()
	mux.HandleFunc("/api/v1/shops", s.HandlerShops)
	mux.HandleFunc("/api/v1/categories", s.HandlerCategories)
	mux.HandleFunc("/api/v1/shop_categories", s.HandlerShopCategories)
	return mux
}
