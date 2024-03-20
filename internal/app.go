package internal

import (
	"context"
	authv1 "github.com/kurneo/go-template/internal/auth/transport/http/v1"
	catv1 "github.com/kurneo/go-template/internal/category/transport/http/v1"
	"github.com/kurneo/go-template/pkg/cache"
	"github.com/kurneo/go-template/pkg/database"
	logPkg "github.com/kurneo/go-template/pkg/log"
	"github.com/kurneo/go-template/pkg/middlewares"
	"github.com/labstack/echo/v4"
	"github.com/spf13/viper"
	"log"
	"net/http"
	"os/signal"
	"sync"
	"syscall"
	"time"
)

var (
	appInstance App
	appOnce     sync.Once
)

type App interface {
	Start()
	GetLogger() logPkg.Contract
	GetCache() cache.Contact
	GetDB() database.Contract
	GetHttpHandler() *echo.Echo
}

type application struct {
	echo *echo.Echo
	lg   logPkg.Contract
	db   database.Contract
	c    cache.Contact
}

func (app *application) Start() {
	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()
	go func() {
		if err := app.echo.Start(":" + app.getHttpPort()); err != nil && err != http.ErrServerClosed {
			log.Fatal("Shutting down the server")
		}
	}()
	<-ctx.Done()
	log.Println("Shutting down")
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	log.Println("Stop http server")
	if err := app.echo.Shutdown(ctx); err != nil {
		log.Fatal(err)
	}
	log.Println("Close database connection")
	err := app.GetDB().Close()
	if err != nil {
		log.Fatal(err)
	}
}

func (app *application) GetLogger() logPkg.Contract {
	return app.lg
}

func (app *application) GetDB() database.Contract {
	return app.db
}

func (app *application) GetCache() cache.Contact {
	return app.c
}

func (app *application) GetHttpHandler() *echo.Echo {
	return app.echo
}

func (app *application) getHttpPort() string {
	return viper.GetString("APP_HTTP_PORT")
}

func NewApplication(
	lg logPkg.Contract,
	db database.Contract,
	c cache.Contact,
	jwtMiddleware echo.MiddlewareFunc,
	authV1 *authv1.Controller,
	catV1 *catv1.Controller,
) App {
	appOnce.Do(func() {
		echoApp := echo.New()
		// configure global middleware here
		echoApp.Use(
			middlewares.CorsMiddleware(),
			middlewares.RateLimiterMiddleware(),
			middlewares.GzipMiddleware(),
			middlewares.AddExposeHeaderMiddleware(),
		)
		echoApp.HideBanner = true

		g := echoApp.Group("/api/admin/v1")

		appInstance = &application{
			echo: echoApp,
			lg:   lg,
			db:   db,
			c:    c,
		}

		authV1.RegisterRoute(g, jwtMiddleware)
		catV1.RegisterRoute(g, jwtMiddleware)
	})

	return appInstance
}

func GetApplication() App {
	return appInstance
}
