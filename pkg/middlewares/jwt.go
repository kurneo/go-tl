package middlewares

import (
	"errors"
	"github.com/golang-jwt/jwt/v4"
	jwtPkg "github.com/kurneo/go-template/pkg/jwt"
	httpPkg "github.com/kurneo/go-template/pkg/support/http"
	echoJwt "github.com/labstack/echo-jwt/v4"
	"github.com/labstack/echo/v4"
)

func JwtMiddleware(t *jwtPkg.TokenManager[int64]) echo.MiddlewareFunc {
	secret, _ := t.GetSecret()
	return echoJwt.WithConfig(echoJwt.Config{
		SigningKey: secret,
		Skipper: func(c echo.Context) bool {
			return false
		},
		TokenLookup: "header:Authorization:Bearer ,query:access_token",
		ParseTokenFunc: func(c echo.Context, tokenString string) (interface{}, error) {
			token, err := t.CheckToken(c.Request().Context(), tokenString)

			if err != nil {
				return nil, err
			}

			c.Set("auth", token)

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
