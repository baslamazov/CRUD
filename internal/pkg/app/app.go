package app

import (
	"EffectiveMobile/internal/config"
	"EffectiveMobile/internal/endpoint"
	"EffectiveMobile/internal/storage/postgres"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"log/slog"
	"net/http"
)

type App struct {
	handler   *endpoint.Endpoint
	dbService *postgres.DBService
	log       *slog.Logger
}

func New(conf *config.Config, log *slog.Logger) (*App, error) {
	app := &App{
		log: log,
	}
	app.dbService = postgres.NewDBService(conf.Database)
	app.handler = endpoint.New(app.dbService)
	return app, nil
}

// Старт апи
func (app *App) Start() error {
	r := chi.NewRouter()

	r.Use(middleware.Logger)

	r.Use(middleware.Recoverer)
	r.Route("/info", func(r chi.Router) {

		// TODO: Как будто бы должна быть пагинация
		r.With(middleware.Heartbeat("/ping")).Get("/", app.handler.Info)
	})
	http.ListenAndServe(":3000", r)

	return nil
}
