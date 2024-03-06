package app

import (
	"context"
	"github.com/kurneo/go-template/config"
	"github.com/kurneo/go-template/pkg/database"
	"github.com/kurneo/go-template/pkg/logger"
	"github.com/kurneo/go-template/pkg/middlewares"
	"github.com/kurneo/go-template/pkg/redis"
	"github.com/kurneo/go-template/pkg/support/jwt"
	"github.com/labstack/echo/v4"
)

type Contract interface {
	Start() error
	Shutdown(ctx context.Context) error
	GetConfig() *config.Config
	GetLogger() logger.Contract
	GetDB() database.Contract
	GetRedis() redis.Contact
	GetHttpHandler() *echo.Echo
	RegisterAdminV1Route(f func(group *echo.Group, jwtMiddleware echo.MiddlewareFunc))
}

type application struct {
	echo          *echo.Echo
	cfg           *config.Config
	lg            logger.Contract
	db            database.Contract
	r             redis.Contact
	jetMiddleware echo.MiddlewareFunc
}

func (app *application) Start() error {
	return app.echo.Start(":" + app.getHttpPort())
}

func (app *application) Shutdown(ctx context.Context) error {
	return app.echo.Shutdown(ctx)
}

func (app *application) GetConfig() *config.Config {
	return app.cfg
}

func (app *application) GetLogger() logger.Contract {
	return app.lg
}

func (app *application) GetDB() database.Contract {
	return app.db
}

func (app *application) GetRedis() redis.Contact {
	return app.r
}

func (app *application) GetHttpHandler() *echo.Echo {
	return app.echo
}

func (app *application) getHttpPort() string {
	return app.cfg.HTTP.Port
}

func (app *application) RegisterAdminV1Route(f func(group *echo.Group, jwtMiddleware echo.MiddlewareFunc)) {
	g := app.echo.Group("/api/admin/v1")
	f(g, app.jetMiddleware)
}

func NewApplication(
	cfg *config.Config,
	lg logger.Contract,
	db database.Contract,
	r redis.Contact,
) (Contract, error) {
	echoApp := echo.New()
	// configure global middleware here
	echoApp.Use(
		middlewares.CorsMiddleware(),
		middlewares.RateLimiterMiddleware(),
		middlewares.GzipMiddleware(),
		middlewares.AddExposeHeaderMiddleware(cfg.HTTP),
	)
	echoApp.HideBanner = true
	// configure global middleware here
	return &application{
		echo:          echoApp,
		cfg:           cfg,
		lg:            lg,
		db:            db,
		r:             r,
		jetMiddleware: middlewares.JwtMiddleware(jwt.NewTokenManager[int64](cfg.JWT, r)),
	}, nil
}
