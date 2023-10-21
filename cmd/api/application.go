package main

import (
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"net/http"
)

type Application struct {
	router *chi.Mux
}

func NewApplication() *Application {
	app := &Application{}
	app.router = chi.NewRouter()

	app.setupRouter()

	return app
}

func (a *Application) setupRouter() {
	a.router.Use(middleware.Logger)

	a.router.Get("/", helloWorldServer)
	a.router.Get("/ping", pingPongHandler)
}

func (a *Application) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	a.router.ServeHTTP(w, r)
}
