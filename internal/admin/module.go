package admin

import (
	"github.com/kurneo/go-template/cmd/app"
	"github.com/kurneo/go-template/internal/admin/datasource"
	"github.com/kurneo/go-template/internal/admin/transport/http/v1"
	"github.com/kurneo/go-template/internal/admin/usecase"
)

func RegisterModule(app app.Contract) error {
	lg := app.GetLogger()
	cfg := app.GetConfig()
	db := app.GetDB()

	tokenRepo := datasource.NewTokenRepo(lg, db)
	userRepo := datasource.NewUserRepo(lg, db)
	tokenManager := usecase.NewTokenManager(cfg.JWT, lg, tokenRepo)
	v1.New(app, usecase.New(userRepo, lg, tokenManager, usecase.NewPasswordChecker()))
	return nil
}
