package routes

import (
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/ituoga/toolbox/hotreload"
	"net/http"
)

func NewDsRouter() chi.Router {
	dsRouter := chi.NewRouter()
	dsRouter.Get("/hotreload", hotreload.Handler)
	return dsRouter
}

func NewRouter(dsRouter chi.Router) chi.Router {
	router := chi.NewMux()
	router.Use(middleware.Logger)
	router.Use(middleware.Recoverer)

	router.Use(middleware.RouteHeaders().Route("Datastar-Request", "true", middleware.New(dsRouter)).Handler)

	router.Get("/static/*", http.StripPrefix("/static/", http.FileServer(http.Dir("static"))).ServeHTTP)

	return router
}
