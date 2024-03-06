package middlewares

import (
	"context"
	"errors"
	"fmt"
	"github.com/golang-jwt/jwt/v4"
	"github.com/kurneo/go-template/config"
	"github.com/kurneo/go-template/internal/admin/entities"
	errPkg "github.com/kurneo/go-template/pkg/error"
	httpPkg "github.com/kurneo/go-template/pkg/support/http"
	echoJwt "github.com/labstack/echo-jwt/v4"
	"github.com/labstack/echo/v4"
)

type UseCase interface {
	CheckToken(ctx context.Context, token string) (*entities.AdminAccessToken, errPkg.Contract)
}

func JwtMiddleware(cfg config.JWT, useCase UseCase) echo.MiddlewareFunc {
	signingKey := []byte(cfg.Secret)

	keyFunc := func(token *jwt.Token) (interface{}, error) {
		if token.Method.Alg() != echoJwt.AlgorithmHS256 {
			return nil, fmt.Errorf("unexpected jwt signing method=%v", token.Header["alg"])
		}
		return signingKey, nil
	}

	return echoJwt.WithConfig(echoJwt.Config{
		SigningKey: signingKey,
		Skipper: func(c echo.Context) bool {
			return false
		},
		TokenLookup: "header:Authorization:Bearer ,query:access_token",
		ParseTokenFunc: func(c echo.Context, auth string) (interface{}, error) {
			token, err := jwt.ParseWithClaims(auth, jwt.MapClaims{}, keyFunc)

			if err != nil {
				if err.Error() == "token contains an invalid number of segments" {
					return nil, jwt.ErrTokenMalformed
				}
				return nil, err
			}

			if !token.Valid {
				return nil, jwt.ErrTokenMalformed
			}

			t, errCheck := useCase.CheckToken(c.Request().Context(), auth)

			if errCheck != nil {
				return nil, errCheck.GetError()
			}

			c.Set("auth", t.Admin)

			return token, nil
		},
		ErrorHandler: func(c echo.Context, err error) error {
			if err.Error() == "missing value in request header" || err.Error() == "missing value in the query string" {
				return httpPkg.ResponseUnauthorized(c)
			}
			switch true {
			case errors.Is(err, jwt.ErrTokenNotValidYet):
				return httpPkg.ResponseBadRequest(c, "token is invalid")
			case errors.Is(err, jwt.ErrTokenExpired):
				return httpPkg.ResponseBadRequest(c, "token is expired")
			case errors.Is(err, jwt.ErrTokenMalformed):
				return httpPkg.ResponseBadRequest(c, "token is malformed")
			case errors.Is(err, jwt.ErrTokenSignatureInvalid):
				return httpPkg.ResponseBadRequest(c, "signature is invalid")
			default:
				return httpPkg.ResponseError(c, err.Error())
			}
		},
	})
}
