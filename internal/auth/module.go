package auth

import (
	"github.com/kurneo/go-template/cmd/app"
	"github.com/kurneo/go-template/internal/auth/datasource"
	"github.com/kurneo/go-template/internal/auth/transport/http/v1"
	"github.com/kurneo/go-template/internal/auth/usecase"
)

func RegisterModule(app app.Contract) error {
	lg := app.GetLogger()
	cfg := app.GetConfig()
	db := app.GetDB()

	userRepo := datasource.NewUserRepo[int64](lg, db)
	tokenManager := usecase.NewTokenManager[int64](cfg.JWT, lg)
	v1.New(app, usecase.New[int64](userRepo, lg, tokenManager, usecase.NewPasswordChecker()))
	return nil
}
