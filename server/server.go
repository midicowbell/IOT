package server

import (
	"iot/handlers"
)

type Server struct {
	handlers *handlers.HTTPHandlers
}

func NewServer(handlers *handlers.HTTPHandlers) *Server {
	return &Server{
		handlers: handlers,
	}
}

// func (s *Server) StartServer() error {
// 	router := mux.NewRouter()

// }
