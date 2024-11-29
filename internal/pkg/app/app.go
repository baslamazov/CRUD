package app

import (
	"EffectiveMobile/internal/config"
	"EffectiveMobile/internal/endpoint"
	"EffectiveMobile/internal/services"
	"EffectiveMobile/internal/storage/postgres"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"log"
	"log/slog"
	"net/http"
)

type App struct {
	handler *endpoint.Endpoint
	service *services.Service
	log     *slog.Logger
}

func New(conf *config.Config, log *slog.Logger) (*App, error) {
	app := &App{
		log: log,
	}
	// TODO: проверить показатели производительности, если каждому выделать отдельный сервис с пулом соединений
	storage := postgres.NewDBService(conf.Database)
	app.service = services.New(log, storage, storage, storage)
	app.handler = endpoint.New(app.service)
	return app, nil
}

// Старт апи
func (app *App) Start() error {
	r := chi.NewRouter()

	r.Use(middleware.Logger)

	r.Use(middleware.Recoverer)
	r.Route("/library", func(r chi.Router) {
		r.With(middleware.Heartbeat("/ping")).Get("/song", app.handler.GetSong)
		r.With(middleware.Heartbeat("/ping")).Delete("/song", app.handler.DeleteSong)
		r.With(middleware.Heartbeat("/ping")).Post("/song", func(writer http.ResponseWriter, request *http.Request) {

		})
		r.With(middleware.Heartbeat("/ping")).Get("/lyric", app.handler.GetLyric)

	})
	go func() {
		if err := http.ListenAndServe(":3000", r); err != nil {
			log.Fatal(err)
		}
	}()

	return nil
}
