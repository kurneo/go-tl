package main

import (
	"github.com/kurneo/go-template/cmd/app"
	"github.com/kurneo/go-template/config"
	"github.com/kurneo/go-template/internal/admin"
	"github.com/kurneo/go-template/internal/category"
	"github.com/kurneo/go-template/internal/finder"
	"github.com/kurneo/go-template/pkg/database"
	"github.com/kurneo/go-template/pkg/logger"
	"log"
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

	// Create new application
	application, errApp := app.NewApplication(cfg, lg, db)
	if errApp != nil {
		log.Fatalf("Create application module failed: %s", errApp)
	}

	errApp = registerModules(application)
	if errApp != nil {
		log.Fatalf("Register application module failed: %s", errApp)
	}

	errApp = application.Start()
	if errApp != nil {
		log.Fatalf("Start application failed: %s", errApp)
	}
}

func registerModules(app app.Contract) error {
	if err := admin.RegisterModule(app); err != nil {
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
