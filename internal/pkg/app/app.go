package app

import (
	_ "EffectiveMobile/docs"
	"EffectiveMobile/internal/config"
	"EffectiveMobile/internal/endpoint"
	"EffectiveMobile/internal/services"
	"EffectiveMobile/internal/storage/postgres"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/swaggo/http-swagger/v2"
	"log"
	"log/slog"
	"net/http"
)

type App struct {
	Handler *endpoint.Endpoint
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
	app.Handler = endpoint.New(app.service)
	return app, nil
}

// @title Swagger Example API
// @version 1.0
// @description This is a sample server Petstore server.
// @termsOfService http://swagger.io/terms/

// @contact.name API Support
// @contact.url http://www.swagger.io/support
// @contact.email support@swagger.io

// @license.name Apache 2.0
// @license.url http://www.apache.org/licenses/LICENSE-2.0.html

// @host petstore.swagger.io
// @BasePath /v2
func (app *App) Start() error {
	r := chi.NewRouter()

	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)

	// Добавляем Swagger UI
	r.Get("/swagger/*", httpSwagger.Handler(
		httpSwagger.URL("http://localhost:3000/swagger/doc.json"), // The URL pointing to API definition
	))

	r.Route("/library", func(r chi.Router) {
		r.With(middleware.Heartbeat("/ping")).Get("/song", app.Handler.GetSong)
		r.With(middleware.Heartbeat("/ping")).Delete("/song", app.Handler.DeleteSong)
		r.With(middleware.Heartbeat("/ping")).Post("/song", app.Handler.AddSong)
		r.With(middleware.Heartbeat("/ping")).Put("/song", app.Handler.UpdateSong)

		r.With(middleware.Heartbeat("/ping")).Get("/lyric", app.Handler.GetLyric)

	})
	go func() {
		if err := http.ListenAndServe(":3000", r); err != nil {
			log.Fatal(err)
		}
	}()

	return nil
}
