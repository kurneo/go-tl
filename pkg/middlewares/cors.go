package middlewares

import (
	"github.com/kurneo/go-template/config"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"strings"
)

func CorsMiddleware() echo.MiddlewareFunc {
	cfg, _ := config.NewConfig()
	return middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOrigins: strings.Split(cfg.HTTP.AllowOrigin, ","),
	})
}
