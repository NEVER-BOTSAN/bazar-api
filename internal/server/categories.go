package server

import (
	"encoding/json"
	"fmt"
	"net/http"
)

func (s *Server) HandlerCategories(w http.ResponseWriter, r *http.Request) {

	categories, err := s.App.GetCategories()
	if err != nil {
		fmt.Println(err.Error())
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
	fmt.Println(categories)

	respjson, err := json.Marshal(categories)
	if err != nil {
		fmt.Println(err.Error())
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}

	w.WriteHeader(http.StatusOK)
	w.Write(respjson)

}
