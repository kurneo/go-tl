package v1

import (
	"github.com/kurneo/go-template/cmd/app"
	"github.com/kurneo/go-template/internal/category/entities"
	"github.com/kurneo/go-template/internal/category/usecase"
	"github.com/kurneo/go-template/pkg/logger"
	"github.com/kurneo/go-template/pkg/support/http"
	"github.com/kurneo/go-template/pkg/support/slices"
	"github.com/kurneo/go-template/pkg/support/validator"
	"github.com/labstack/echo/v4"
	"strconv"
)

type controller struct {
	a app.Contract
	l logger.Contract
	u usecase.CategoryUseCaseContract
}

func (c controller) List(context echo.Context) error {
	filters := http.GetFilterParams(context, []string{"name", "created_at", "status"})
	page, limit, errPaginate := http.GetPaginateParams(context)
	sorts := http.GetSortParams(context)

	errValidate := http.MergeErrorValidate(errPaginate, validator.ValidateStruct(struct {
		Status string `validate:"omitempty,oneof=10 11"`
	}{Status: filters["status"]}))

	if len(errValidate) > 0 {
		return http.ResponseUnprocessableEntity(context, errValidate)
	}

	list, pg, err := c.u.List(context.Request().Context(), filters, sorts, page, limit)

	if err != nil {
		return http.ResponseError(context, err.GetMessage())
	}

	http.SetHeaderCountAndTotal(context, strconv.Itoa(int(pg.Total)), strconv.Itoa(pg.TotalPages))

	if len(list) == 0 {
		return http.ResponseEmptyList(context)
	}

	return context.JSON(
		200,
		slices.Map[entities.Category, map[string]interface{}](list, func(category entities.Category) map[string]interface{} {
			return category.ToMap()
		}),
	)
}

func (c controller) Get(context echo.Context) error {
	id, errGetId := http.GetIDRouteParam(context)

	if errGetId != nil {
		return http.ResponseUnprocessableEntity(context, errGetId)
	}

	category, err := c.u.Get(context.Request().Context(), int64(id))

	if err != nil {
		return http.ResponseError(context, err.GetMessage())
	}

	if category == nil {
		return http.ResponseNotFound(context)
	}
	return http.ResponseOk(context, category.ToMap())
}

func (c controller) Store(context echo.Context) error {
	body, err := http.ParseFormData[CategoryFormData](context)
	if err != nil {
		c.l.Error(err)
		return http.ResponseBadRequest(context, err.Error())
	}

	if errVald := validator.ValidateStruct(body); len(errVald) > 0 {
		return http.ResponseUnprocessableEntity(context, errVald)
	}

	errTrans := c.a.GetDB().Begin()

	if errTrans != nil {
		c.l.Error(errTrans)
		return http.ResponseError(context, errTrans.Error())
	}
	cat, errCrt := c.u.Store(context.Request().Context(), body)

	if errCrt != nil {
		errTrans = c.a.GetDB().Rollback()
		if errTrans != nil {
			c.l.Error(errTrans)
			return http.ResponseError(context, errTrans.Error())
		}
		if errCrt.IsDomainError() {
			return http.ResponseBadRequest(context, errCrt.GetMessage())
		} else {
			return http.ResponseError(context, errCrt.GetMessage())
		}
	}

	errTrans = c.a.GetDB().Commit()
	if errTrans != nil {
		c.l.Error(errTrans)
		return http.ResponseError(context, errTrans.Error())
	}

	return http.ResponseOk(context, (*cat).ToMap())
}

func (c controller) Update(context echo.Context) error {
	id, errGetId := http.GetIDRouteParam(context)

	if errGetId != nil {
		return http.ResponseUnprocessableEntity(context, errGetId)
	}

	body, errParse := http.ParseFormData[CategoryFormData](context)
	if errParse != nil {
		c.l.Error(errParse)
		return http.ResponseBadRequest(context)
	}

	if errorsValidate := validator.ValidateStruct(body); len(errorsValidate) > 0 {
		return http.ResponseUnprocessableEntity(context, errorsValidate)
	}

	category, errGet := c.u.Get(context.Request().Context(), int64(id))
	if errGet != nil {
		if errGet.IsDomainError() {
			return http.ResponseBadRequest(context, errGet.GetMessage())
		} else {
			return http.ResponseError(context, errGet.GetMessage())
		}
	}

	if category == nil {
		return http.ResponseNotFound(context)
	}

	errTrans := c.a.GetDB().Begin()
	if errTrans != nil {
		c.l.Error(errTrans)
		return http.ResponseError(context, errTrans.Error())
	}
	errUpdate := c.u.Update(context.Request().Context(), category, body)

	if errUpdate != nil {
		errTrans = c.a.GetDB().Rollback()
		if errTrans != nil {
			c.l.Error(errTrans)
			return http.ResponseError(context, errTrans.Error())
		}
		if errUpdate.IsDomainError() {
			return http.ResponseBadRequest(context, errUpdate.GetMessage())
		} else {
			return http.ResponseError(context, errUpdate.GetMessage())
		}
	}
	errTrans = c.a.GetDB().Commit()
	if errTrans != nil {
		c.l.Error(errTrans)
		return http.ResponseError(context, errTrans.Error())
	}

	return http.ResponseOk(context, category.ToMap())

}

func (c controller) Delete(context echo.Context) error {
	id, errGetId := http.GetIDRouteParam(context)

	if errGetId != nil {
		return http.ResponseUnprocessableEntity(context, errGetId)
	}

	category, errGet := c.u.Get(context.Request().Context(), int64(id))
	if errGet != nil {
		if errGet.IsDomainError() {
			return http.ResponseBadRequest(context, errGet.GetMessage())
		} else {
			return http.ResponseError(context, errGet.GetMessage())
		}
	}

	if category == nil {
		return http.ResponseNotFound(context)
	}

	errTrans := c.a.GetDB().Begin()
	if errTrans != nil {
		c.l.Error(errTrans)
		return http.ResponseError(context, errTrans.Error())
	}

	errDel := c.u.Delete(context.Request().Context(), category)

	if errDel != nil {
		errTrans = c.a.GetDB().Rollback()
		if errTrans != nil {
			c.l.Error(errTrans)
			return http.ResponseError(context, errTrans.Error())
		}
		if errDel.IsDomainError() {
			return http.ResponseBadRequest(context, errDel.GetMessage())
		} else {
			return http.ResponseError(context, errDel.GetMessage())
		}
	}
	errTrans = c.a.GetDB().Commit()
	if errTrans != nil {
		c.l.Error(errTrans)
		return http.ResponseError(context, errTrans.Error())
	}

	return http.ResponseOk(context, true)
}

func New(a app.Contract, u usecase.CategoryUseCaseContract) {
	c := &controller{l: a.GetLogger(), u: u, a: a}
	a.RegisterAdminV1Route(func(group *echo.Group, jwtMiddleware echo.MiddlewareFunc) {
		g := group.Group("/categories", jwtMiddleware)
		g.GET("", c.List)
		g.POST("", c.Store)
		g.GET("/:id", c.Get)
		g.PUT("/:id", c.Update)
		g.DELETE("/:id", c.Delete)
	})
}
