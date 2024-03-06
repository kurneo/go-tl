package v1

import (
	"github.com/golang-jwt/jwt/v4"
	"github.com/kurneo/go-template/cmd/app"
	"github.com/kurneo/go-template/internal/admin/entities"
	"github.com/kurneo/go-template/internal/admin/usecase"
	"github.com/kurneo/go-template/pkg/logger"
	"github.com/kurneo/go-template/pkg/support/http"
	"github.com/kurneo/go-template/pkg/support/validator"
	"github.com/labstack/echo/v4"
)

type controller struct {
	a app.Contract
	l logger.Contract
	u usecase.AdminUseCaseContract
}

func (ctl controller) Login(context echo.Context) error {
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

func (ctl controller) Me(context echo.Context) error {
	user := context.Get("auth").(*entities.Admin)
	return http.ResponseOk(context, user.ToMap())
}

func (ctl controller) RefreshToken(context echo.Context) error {
	token := (context.Get("user").(*jwt.Token)).Raw

	errTrans := ctl.a.GetDB().Begin()
	if errTrans != nil {
		ctl.l.Error(errTrans)
		return http.ResponseError(context, errTrans.Error())
	}

	authToken, err := ctl.u.RefreshToken(context.Request().Context(), token)
	if err != nil {
		errTrans = ctl.a.GetDB().Rollback()
		if errTrans != nil {
			ctl.l.Error(errTrans)
			return http.ResponseError(context, errTrans.Error())
		}
		return http.ResponseError(context, err.GetMessage())
	}

	errTrans = ctl.a.GetDB().Commit()
	if errTrans != nil {
		ctl.l.Error(errTrans)
		return http.ResponseError(context, errTrans.Error())
	}
	return http.ResponseOk(context, authToken)
}

func (ctl controller) Logout(context echo.Context) error {
	token := (context.Get("user").(*jwt.Token)).Raw
	err := ctl.u.Logout(context.Request().Context(), token)
	if err != nil {
		return http.ResponseError(context, err.GetMessage())
	}
	return http.ResponseOk(context, true)
}

func New(a app.Contract, u usecase.AdminUseCaseContract) {
	c := &controller{l: a.GetLogger(), u: u, a: a}
	a.RegisterAdminV1Route(func(group *echo.Group, jwtMiddleware echo.MiddlewareFunc) {
		g := group.Group("/auth")
		g.POST("/login", c.Login)
		g.GET("/me", c.Me, jwtMiddleware)
		g.POST("/refresh-token", c.RefreshToken, jwtMiddleware)
		g.POST("/logout", c.Logout, jwtMiddleware)
	})
}
