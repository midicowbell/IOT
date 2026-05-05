package server

import (
	"errors"
	"iot/handlers"
	"net/http"

	"github.com/gorilla/mux"
)

type Server struct {
	handlers *handlers.HTTPHandlers
}

func NewServer(handlers *handlers.HTTPHandlers) *Server {
	return &Server{
		handlers: handlers,
	}
}

func (s *Server) StartServer() error {
	router := mux.NewRouter()
	router.HandleFunc("/shelves", s.handlers.HandleGetStatus).Methods("GET")
	router.HandleFunc("/shelves/weight", s.handlers.HandleUpdateWeight).Methods("PATCH")
	router.HandleFunc("/shelves", s.handlers.HandleAddShelf).Methods("POST")
	router.HandleFunc("/shelves/{id}", s.handlers.DeleteShelf).Methods("DELETE")
	router.HandleFunc("/shelves/product", s.handlers.HandleAddProduct).Methods("POST")
	router.HandleFunc("/shelves/{id}/product", s.handlers.DeleteProduct).Methods("DELETE")
	router.HandleFunc("/shelves/{id}/product", s.handlers.GetProduct).Methods("GET")

	if err := http.ListenAndServe(":8080", router); err != nil {
		if errors.Is(err, http.ErrServerClosed) {
			return nil
		}
		return err
	}
	return nil
}
