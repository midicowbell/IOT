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
	router.HandleFunc("/api/shelves", s.handlers.HandleGetStatus).Methods("GET")
	router.HandleFunc("/api/shelves", s.handlers.HandleAddShelf).Methods("POST")
	router.HandleFunc("/api/shelves/{id}", s.handlers.DeleteShelf).Methods("DELETE")
	router.HandleFunc("/api/shelves/weight", s.handlers.HandleUpdateWeight).Methods("PATCH", "PUT")
	router.HandleFunc("/api/products", s.handlers.HandleAddProduct).Methods("POST")
	router.HandleFunc("/api/products/{id}", s.handlers.GetProduct).Methods("GET")
	router.HandleFunc("/api/products/{id}", s.handlers.DeleteProduct).Methods("DELETE")

	if err := http.ListenAndServe(":8080", router); err != nil {
		if errors.Is(err, http.ErrServerClosed) {
			return nil
		}
		return err
	}
	return nil
}
