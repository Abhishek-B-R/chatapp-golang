package routes

import (
	"github.com/Abhishek-B-R/chat-app-golang/internals/app"
	"github.com/go-chi/chi"
)

func SetupRoutes(app *app.Application) *chi.Mux{
	r := chi.NewRouter()

	r.Get("/health",app.HealthCheck)
	return r
}