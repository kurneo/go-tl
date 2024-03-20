package v1

import (
	"github.com/kurneo/go-template/internal/auth/domain/usecase"
	"github.com/kurneo/go-template/pkg/database"
	jwtPkg "github.com/kurneo/go-template/pkg/jwt"
	"github.com/kurneo/go-template/pkg/log"
	"github.com/kurneo/go-template/pkg/support/http"
	"github.com/kurneo/go-template/pkg/support/validator"
	"github.com/labstack/echo/v4"
)

type Controller struct {
	l  log.Contract
	db database.Contract
	u  usecase.UserUseCaseContract
}

func (ctl Controller) Login(context echo.Context) error {
	body, err := http.ParseFormData[LoginFormData](context)
	if err != nil {
		ctl.l.Error(err)
		return http.ResponseBadRequest(context, err.Error())
	}

	if errValid := validator.ValidateStruct(body); len(errValid) > 0 {
		return http.ResponseUnprocessableEntity(context, errValid)
	}

	errTrans := ctl.db.Begin()
	if errTrans != nil {
		ctl.l.Error(errTrans)
		return http.ResponseError(context, errTrans.Error())
	}

	token, errLogin := ctl.u.Login(context.Request().Context(), body.Email, body.Password)

	if errLogin != nil {
		errTrans = ctl.db.Rollback()
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

	errTrans = ctl.db.Commit()
	if errTrans != nil {
		ctl.l.Error(errTrans)
		return http.ResponseError(context, errTrans.Error())
	}

	return context.JSON(200, token)
}

func (ctl Controller) Me(context echo.Context) error {
	auth := context.Get("auth").(*jwtPkg.AccessToken[int64])
	user, err := ctl.u.GetProfile(context.Request().Context(), auth.Sub)
	if err != nil {
		return http.ResponseError(context, err.GetMessage())
	}
	return http.ResponseOk(context, user)
}

func (ctl Controller) Logout(context echo.Context) error {
	auth := context.Get("auth").(*jwtPkg.AccessToken[int64])
	err := ctl.u.Logout(context.Request().Context(), auth)
	if err != nil {
		return http.ResponseError(context, err.GetMessage())
	}
	return http.ResponseOk(context, true)
}

func (ctl Controller) RegisterRoute(group *echo.Group, jwtMiddleware echo.MiddlewareFunc) {
	g := group.Group("/auth")
	g.POST("/login", ctl.Login)
	g.GET("/me", ctl.Me, jwtMiddleware)
	g.POST("/logout", ctl.Logout, jwtMiddleware)
}

func NewHtpV1Controller(
	u usecase.UserUseCaseContract,
	l log.Contract,
	db database.Contract,
) *Controller {
	return &Controller{l: l, u: u, db: db}
}
