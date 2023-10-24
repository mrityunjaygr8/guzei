package server

import (
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/mrityunjaygr8/guzei/store"
	"net/http"
)

type Server struct {
	router *chi.Mux
	store  store.GuzeiStore
}

func NewServer(store store.GuzeiStore) *Server {
	app := &Server{store: store}
	app.router = chi.NewRouter()

	app.setupRouter()

	return app
}

func (a *Server) setupRouter() {
	a.router.Use(middleware.Logger)

	a.router.Get("/", a.helloWorldServer)
	a.router.Get("/ping", a.pingPongHandler)
	a.router.Post("/users", a.UserInsert)
	a.router.Get("/users", a.UserList)
}

func (a *Server) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	a.router.ServeHTTP(w, r)
}
func (a *Server) helloWorldServer(w http.ResponseWriter, _ *http.Request) {
	w.Write([]byte("Hello World"))
}

func (a *Server) pingPongHandler(w http.ResponseWriter, _ *http.Request) {
	w.Write([]byte("PONG"))
}
