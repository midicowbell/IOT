package server

import (
	"iot/handlers"

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
	router.HandleFunc("/shelves/{id}/weight", s.handlers.HandleUpdateWeight).Methods("PATCH")
	router.HandleFunc("/shelves", s.handlers.HandleAddShelf).Methods("POST")
	router.HandleFunc("/shelves/{id}", s.handlers.DeleteShelf).Methods("DELETE")

}
