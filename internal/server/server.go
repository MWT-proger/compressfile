package server

import (
	"log"
	"net/http"
	"time"

	"github.com/MWT-proger/compressfile/configs"
	"github.com/go-chi/chi"
)

type Server struct {
	Router *chi.Mux
}

// Запускает Http server
// при сбое возвращает ошибку
func (s *Server) Run() error {

	conf := configs.GetConfig()

	log.Println("Running server on", conf.HostServer)

	server := &http.Server{
		Addr:    conf.HostServer,
		Handler: s.Router,

		ReadTimeout:  10 * time.Second,
		WriteTimeout: 30 * time.Second,
		IdleTimeout:  120 * time.Second,
	}

	return server.ListenAndServe()
}
