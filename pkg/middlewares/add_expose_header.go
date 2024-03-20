package middlewares

import (
	"github.com/labstack/echo/v4"
)

func AddExposeHeaderMiddleware() echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(context echo.Context) error {
			context.Response().Header().Set("Access-Control-Expose-Headers", "")
			return next(context)
		}
	}
}
