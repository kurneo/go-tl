package v1

import (
	"github.com/kurneo/go-template/cmd/app"
	"github.com/kurneo/go-template/internal/auth/usecase"
	"github.com/kurneo/go-template/pkg/logger"
	"github.com/kurneo/go-template/pkg/support/http"
	jwtPkg "github.com/kurneo/go-template/pkg/support/jwt"
	"github.com/kurneo/go-template/pkg/support/validator"
	"github.com/labstack/echo/v4"
)

type controller[T jwtPkg.SubType] struct {
	a app.Contract
	l logger.Contract
	u usecase.UserUseCaseContract[T]
}

func (ctl controller[T]) Login(context echo.Context) error {
	body, err := http.ParseFormData[LoginFormData](context)
	if err != nil {
		ctl.l.Error(err)
		return http.ResponseBadRequest(context, err.Error())
	}

	if errValid := validator.ValidateStruct(body); len(errValid) > 0 {
		return http.ResponseUnprocessableEntity(context, errValid)
	}

	errTrans := ctl.a.GetDB().Begin()
	if errTrans != nil {
		ctl.l.Error(errTrans)
		return http.ResponseError(context, errTrans.Error())
	}

	token, errLogin := ctl.u.Login(context.Request().Context(), body.Email, body.Password)

	if errLogin != nil {
		errTrans = ctl.a.GetDB().Rollback()
		if errTrans != nil {
			ctl.l.Error(errTrans)
			return http.ResponseError(context, errTrans.Error())
		}
		if errLogin.IsDomainError() {
			return http.ResponseBadRequest(context, errLogin.GetMessage())
		} else {
			return http.ResponseError(context, errLogin.GetMessage())
		}
	}

	errTrans = ctl.a.GetDB().Commit()
	if errTrans != nil {
		ctl.l.Error(errTrans)
		return http.ResponseError(context, errTrans.Error())
	}

	return context.JSON(200, token)
}

func (ctl controller[T]) Me(context echo.Context) error {
	auth := context.Get("auth").(*jwtPkg.AccessToken[T])
	user, err := ctl.u.GetProfile(context.Request().Context(), auth.Sub)
	if err != nil {
		return http.ResponseError(context, err.GetMessage())
	}
	return http.ResponseOk(context, user)
}

func (ctl controller[T]) Logout(context echo.Context) error {
	auth := context.Get("auth").(*jwtPkg.AccessToken[T])
	err := ctl.u.Logout(context.Request().Context(), auth)
	if err != nil {
		return http.ResponseError(context, err.GetMessage())
	}
	return http.ResponseOk(context, true)
}

func New[T jwtPkg.SubType](a app.Contract, u usecase.UserUseCaseContract[T]) {
	c := &controller[T]{l: a.GetLogger(), u: u, a: a}
	a.RegisterAdminV1Route(func(group *echo.Group, jwtMiddleware echo.MiddlewareFunc) {
		g := group.Group("/auth")
		g.POST("/login", c.Login)
		g.GET("/me", c.Me, jwtMiddleware)
		g.POST("/logout", c.Logout, jwtMiddleware)
	})
}
