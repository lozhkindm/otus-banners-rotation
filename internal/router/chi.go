package router

import (
	"net/http"

	"github.com/go-chi/chi/v5" //nolint:typecheck
	"github.com/go-chi/chi/v5/middleware"
)

type Router interface {
	ServeHTTP(w http.ResponseWriter, r *http.Request)
}

type Application interface {
	GetName() string
}

type Handlers interface {
	AddBanner(w http.ResponseWriter, r *http.Request)
	RemoveBanner(w http.ResponseWriter, r *http.Request)
	ClickBanner(w http.ResponseWriter, r *http.Request)
	PickBanner(w http.ResponseWriter, r *http.Request)
}

func NewChiRouter(handlers Handlers) Router {
	mux := chi.NewRouter() //nolint:typecheck
	mux.Use(middleware.Logger)

	mux.Post("/banner/rotation", handlers.AddBanner)
	mux.Delete("/banner/rotation", handlers.RemoveBanner)
	mux.Post("/banner/click", handlers.ClickBanner)
	mux.Get("/banner/pick", handlers.PickBanner)
	return mux
}
