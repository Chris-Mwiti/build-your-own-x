package worker

import (
	"context"
	"log"
	"net/http"
	"time"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

type WorkerApi struct {
	Address string
	Port string
	Worker *Worker
	Router *chi.Mux
}

func router() (*chi.Mux) {

	//create a new instance of the router
	router := chi.NewRouter()

	//setting up a good middleware stack to receive call
	router.Use(middleware.RequestID)	
	router.Use(middleware.RealIP)
	router.Use(middleware.Logger)
	//will automatically return an internal server error whenever the program panics
	router.Use(middleware.Recoverer)
	router.Use(middleware.Timeout(60 * time.Second))

	return router
}



