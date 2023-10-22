package main

import (
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"net/http"
)

type Application struct {
	router *chi.Mux
	store  GuzeiStore
}

func NewApplication(store GuzeiStore) *Application {
	app := &Application{store: store}
	app.router = chi.NewRouter()

	app.setupRouter()

	return app
}

func (a *Application) setupRouter() {
	a.router.Use(middleware.Logger)

	a.router.Get("/", a.helloWorldServer)
	a.router.Get("/ping", a.pingPongHandler)
}

func (a *Application) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	a.router.ServeHTTP(w, r)
}
func (a *Application) helloWorldServer(w http.ResponseWriter, _ *http.Request) {
	w.Write([]byte("Hello World"))
}

func (a *Application) pingPongHandler(w http.ResponseWriter, _ *http.Request) {
	w.Write([]byte("PONG"))
}
