package main

import (
	"context"
	"github.com/kurneo/go-template/cmd/app"
	"github.com/kurneo/go-template/config"
	"github.com/kurneo/go-template/internal/auth"
	"github.com/kurneo/go-template/internal/category"
	"github.com/kurneo/go-template/internal/finder"
	"github.com/kurneo/go-template/pkg/database"
	"github.com/kurneo/go-template/pkg/logger"
	"github.com/kurneo/go-template/pkg/redis"
	"log"
	"net/http"
	"os/signal"
	"syscall"
	"time"
)

func main() {
	// Init config
	cfg, errCfg := config.NewConfig()

	if errCfg != nil {
		log.Fatalf("Config error: %s", errCfg)
	}

	// Init logger
	lg, errLg := logger.New(cfg.Log)
	if errLg != nil {
		log.Fatalf("Init logger error: %s", errLg)
	}

	// Init DB
	db, errDB := database.New(cfg.DB)
	defer func() {
		errDB = db.Close()
		if errDB != nil {
			log.Fatalf("DB error: close connection failed: %s", errDB)
		}
	}()
	if errDB != nil {
		log.Fatalf("Init DB error: %s", errDB)
	}
	// Init redis
	r, errR := redis.NewRedisClient(cfg.Redis)
	if errDB != nil {
		log.Fatalf("Init Redis client error: %s", errR)
	}

	// Create new application
	application, errApp := app.NewApplication(cfg, lg, db, r)
	if errApp != nil {
		log.Fatalf("Create application module failed: %s", errApp)
	}

	errApp = registerModules(application)
	if errApp != nil {
		log.Fatalf("Register application module failed: %s", errApp)
	}

	//Graceful Shutdown
	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()
	go func() {
		if err := application.Start(); err != nil && err != http.ErrServerClosed {
			log.Fatal("Shutting down the server")
		}
	}()
	<-ctx.Done()
	log.Println("Shutting down")
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	log.Println("Stop http server")
	if err := application.Shutdown(ctx); err != nil {
		log.Fatal(err)
	}
	log.Println("Close database connection")
	err := db.Close()
	if err != nil {
		log.Fatal(err)
	}
}

func registerModules(app app.Contract) error {
	if err := auth.RegisterModule(app); err != nil {
		return err
	}

	if err := category.RegisterModule(app); err != nil {
		return err
	}

	if err := finder.RegisterModule(app); err != nil {
		return err
	}

	return nil
}
