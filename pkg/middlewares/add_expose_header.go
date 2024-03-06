package middlewares

import (
	"github.com/kurneo/go-template/config"
	"github.com/labstack/echo/v4"
)

func AddExposeHeaderMiddleware(cfg config.HTTP) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(context echo.Context) error {
			context.Response().Header().Set("Access-Control-Expose-Headers", cfg.ExposeHeader)
			return next(context)
		}
	}
}
