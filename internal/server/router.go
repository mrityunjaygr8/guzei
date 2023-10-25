package server

import (
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/httplog/v2"
	"time"
)

func (a *Server) setupRouter() {
	// A good base middleware stack
	a.router.Use(middleware.RequestID)
	a.router.Use(middleware.RealIP)
	a.router.Use(middleware.Recoverer)
	a.router.Use(httplog.RequestLogger(a.logger))
	a.router.Use(middleware.Heartbeat("/api/v1/heartbeat"))

	// Set a timeout value on the request context (ctx), that will signal
	// through ctx.Done() that the request has timed out and further
	// processing should be stopped.
	a.router.Use(middleware.Timeout(60 * time.Second))

	v1Router := chi.NewRouter()

	//v1Router.Get("/", a.helloWorldServer)
	v1Router.Get("/ping", a.pingPongHandler)
	v1Router.Post("/users", a.UserInsert)
	v1Router.Get("/users", a.UserList)

	a.router.Mount("/api/v1", v1Router)
}
