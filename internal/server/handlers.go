package server

import (
	"github.com/getsentry/sentry-go"
	sentrynegroni "github.com/getsentry/sentry-go/negroni"
	"github.com/gorilla/mux"
	"github.com/phyber/negroni-gzip/gzip"
	"github.com/rs/cors"
	"github.com/unrolled/secure"
	"github.com/urfave/negroni"

	cityHttp "weather4you/internal/city/delivery/http"
	cityRepository "weather4you/internal/city/repository"
	cityUseCase "weather4you/internal/city/usecase"
	apiMiddlewares "weather4you/internal/middleware"
	"weather4you/pkg/metric"

	httpSwagger "github.com/swaggo/http-swagger"
	"github.com/swaggo/swag/example/basic/docs"
)

// Map Server Handlers
func (s *Server) MapHandlers(router *mux.Router) error {
	metrics, err := metric.CreateMetrics(s.cfg.Metrics.URL, s.cfg.Metrics.ServiceName)
	if err != nil {
		s.logger.Errorf("CreateMetrics Error: %s", err)
	}
	s.logger.Info(
		"Metrics available URL: %s, ServiceName: %s",
		s.cfg.Metrics.URL,
		s.cfg.Metrics.ServiceName,
	)

	// Init repositories
	cRepo := cityRepository.NewNewsRepository(s.db)

	// Init useCases
	cityUC := cityUseCase.NewCityUseCase(s.cfg, cRepo, s.logger)

	// Init handlers
	cityHandlers := cityHttp.NewCityHandlers(s.cfg, cityUC, s.logger)

	mw := apiMiddlewares.NewMiddlewareManager(s.cfg, s.logger)

	n := negroni.New()
	if err := sentry.Init(sentry.ClientOptions{
		Dsn:              "https://6bbac42b0d4279903503ca03f60cc979@o4507918794620928.ingest.de.sentry.io/4507918798291024",
		TracesSampleRate: 1.0,
	}); err != nil {
		s.logger.Infof("Sentry initialization failed: %v\n", err)
	}
	n.Use(sentrynegroni.New(sentrynegroni.Options{}))
	recovery := negroni.NewRecovery()
	recovery.PanicHandlerFunc = sentrynegroni.PanicHandlerFunc

	n.Use(recovery)
	n.UseHandler(router)

	n.Use(mw.RequestLoggerMiddleware())

	docs.SwaggerInfo.Title = "Go weather REST API"
	router.HandleFunc("/swagger/*", httpSwagger.WrapHandler)

	c := cors.New(cors.Options{
		AllowedOrigins: []string{"*"},
	})
	n.Use(c)
	n.Use(mw.MetricsMiddleware(metrics))

	n.Use(gzip.Gzip(gzip.DefaultCompression))
	secureMiddleware := secure.New(secure.Options{
		FrameDeny:          true,
		ContentTypeNosniff: true,
		BrowserXssFilter:   true,
	})
	n.Use(negroni.HandlerFunc(secureMiddleware.HandlerFuncWithNext))
	if s.cfg.Server.Debug {
		n.Use(mw.DebugMiddleware())
	}

	v1 := router.PathPrefix("/api/v1").Subrouter()

	cityGroup := v1.PathPrefix("/city").Subrouter()

	cityHttp.MapCityRoutes(cityGroup, cityHandlers, mw)

	return nil
}
