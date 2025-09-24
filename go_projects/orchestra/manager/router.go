package manager

import (
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

type ManagerApi struct {
	Address string
	Port int
	Router *chi.Mux
	Manger *Manager
}

func router() (*chi.Mux){
	router := chi.NewRouter()
	
	//set up the necessary middlewares
	router.Use(middleware.RequestID)
	router.Use(middleware.RealIP)
	router.Use(middleware.Logger)
	router.Use(middleware.Recoverer)
	router.Use(middleware.Timeout(60 * time.Second))
	

	return router
}


