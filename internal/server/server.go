package server

import (
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/httplog/v2"
	ut "github.com/go-playground/universal-translator"
	"github.com/go-playground/validator/v10"
	"github.com/mrityunjaygr8/guzei/store"
	"net/http"
)

type Server struct {
	router     *chi.Mux
	store      store.GuzeiStore
	logger     *httplog.Logger
	validator  *validator.Validate
	translator ut.Translator
}

func (a *Server) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	a.router.ServeHTTP(w, r)
}

func NewServer(store store.GuzeiStore, logger *httplog.Logger) (*Server, error) {
	app := &Server{store: store, logger: logger}
	app.router = chi.NewRouter()
	app.setupRouter()
	err := app.setupValidation()
	if err != nil {
		return nil, err
	}

	return app, nil
}

func (a *Server) helloWorldServer(w http.ResponseWriter, _ *http.Request) {
	w.Write([]byte("Hello World"))
}

func (a *Server) pingPongHandler(w http.ResponseWriter, _ *http.Request) {
	w.Write([]byte("PONG"))
}
