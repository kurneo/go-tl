package pkg

import (
	"github.com/google/wire"
	"github.com/kurneo/go-template/pkg/cache"
	"github.com/kurneo/go-template/pkg/database"
	"github.com/kurneo/go-template/pkg/jwt"
	logPkg "github.com/kurneo/go-template/pkg/log"
	"github.com/kurneo/go-template/pkg/middlewares"
	"github.com/labstack/echo/v4"
	"log"
)

var WireSet = wire.NewSet(
	ResolveCacheInstance,
	ResolveDatabaseInstance,
	ResolveLogInstance,
	ResolveTokenManager,
	ResolveJWTMiddlewareFunc,
)

func ResolveCacheInstance() cache.Contact {
	c, err := cache.New()
	if err != nil {
		log.Fatalf("init cache error: %s", err)
	}
	return c
}

func ResolveDatabaseInstance() database.Contract {
	d, err := database.New()
	if err != nil {
		log.Fatalf("init database error: %s", err)
	}
	return d
}

func ResolveLogInstance() logPkg.Contract {
	l, err := logPkg.New()
	if err != nil {
		log.Fatalf("init logger error: %s", err)
	}
	return l
}

func ResolveTokenManager(c cache.Contact) *jwt.TokenManager[int64] {
	return jwt.NewTokenManager[int64](c)
}

func ResolveJWTMiddlewareFunc(t *jwt.TokenManager[int64]) echo.MiddlewareFunc {
	return middlewares.JwtMiddleware(t)
}
