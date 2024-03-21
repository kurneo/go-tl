package internal

import (
	"context"
	authv1 "github.com/kurneo/go-template/internal/auth/transport/http/v1"
	catv1 "github.com/kurneo/go-template/internal/category/transport/http/v1"
	"github.com/kurneo/go-template/pkg/cache"
	"github.com/kurneo/go-template/pkg/database"
	"github.com/kurneo/go-template/pkg/hashing"
	logPkg "github.com/kurneo/go-template/pkg/log"
	"github.com/labstack/echo/v4"
	"log"
	"net/http"
	"os/signal"
	"strconv"
	"sync"
	"syscall"
	"time"
)

var (
	appInstance App
	appOnce     sync.Once
)

type App interface {
	Start(p int)
	GetLogger() logPkg.Contract
	GetCache() cache.Contact
	GetDB() database.Contract
	GetHashing() hashing.Contact
	GetHttpHandler() *echo.Echo
}

type application struct {
	e  *echo.Echo
	lg logPkg.Contract
	db database.Contract
	c  cache.Contact
	s  hashing.Contact
}

// Start server with gracefully shutdown.
func (app *application) Start(p int) {
	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()
	go func() {
		if err := app.e.Start(":" + strconv.Itoa(p)); err != nil && err != http.ErrServerClosed {
			log.Fatal("Shutting down the server")
		}
	}()
	<-ctx.Done()
	log.Println("Shutting down")
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	log.Println("Stop http server")
	if err := app.e.Shutdown(ctx); err != nil {
		log.Fatal(err)
	}
	log.Println("Close database connection")
	err := app.GetDB().Close()
	if err != nil {
		log.Fatal(err)
	}
}

// GetLogger used by application
func (app *application) GetLogger() logPkg.Contract {
	return app.lg
}

// GetDB used by application
func (app *application) GetDB() database.Contract {
	return app.db
}

// GetCache used by application
func (app *application) GetCache() cache.Contact {
	return app.c
}

// GetHashing used by application
func (app *application) GetHashing() hashing.Contact {
	return app.s
}

// GetHttpHandler that create server
func (app *application) GetHttpHandler() *echo.Echo {
	return app.e
}

// NewApplication make new application with can start/stop
func NewApplication(
	e *echo.Echo,
	lg logPkg.Contract,
	db database.Contract,
	c cache.Contact,
	s hashing.Contact,
	authV1 *authv1.Controller,
	catV1 *catv1.Controller,
) App {
	appOnce.Do(func() {
		appInstance = &application{
			e:  e,
			lg: lg,
			db: db,
			c:  c,
			s:  s,
		}
		g := e.Group("/api/admin/v1")
		authV1.RegisterRoute(g)
		catV1.RegisterRoute(g)
	})

	return appInstance
}

func GetApplication() App {
	return appInstance
}
