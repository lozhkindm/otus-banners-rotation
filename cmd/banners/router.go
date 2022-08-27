package main

import (
	"errors"

	"github.com/lozhkindm/otus-banners-rotation/internal/app"
	"github.com/lozhkindm/otus-banners-rotation/internal/handlers"
	internalrouter "github.com/lozhkindm/otus-banners-rotation/internal/router"
)

const routerChi = "chi"

var undefinedRouterType = errors.New("undefined router type")

func NewRouter(handlers *handlers.Handlers, config Config) (app.Router, error) {
	var router app.Router

	switch config.App.RouterType {
	case routerChi:
		router = internalrouter.NewChiRouter(handlers)
	default:
		return nil, undefinedRouterType
	}

	return router, nil
}
