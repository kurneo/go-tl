package v1

import (
	"github.com/kurneo/go-template/internal/category/domain/entity"
	"github.com/kurneo/go-template/internal/category/domain/usecase"
	"github.com/kurneo/go-template/pkg/database"
	"github.com/kurneo/go-template/pkg/log"
	"github.com/kurneo/go-template/pkg/support/http"
	"github.com/kurneo/go-template/pkg/support/slices"
	"github.com/kurneo/go-template/pkg/support/validator"
	"github.com/labstack/echo/v4"
	"strconv"
)

type Controller struct {
	l  log.Contract
	db database.Contract
	u  usecase.CategoryUseCaseContract
}

func (c Controller) List(context echo.Context) error {
	filters := http.GetFilterParams(context, []string{"name", "created_at", "status"})
	page, limit, errPaginate := http.GetPaginateParams(context)
	sorts := http.GetSortParams(context)

	errValidate := http.MergeErrorValidate(errPaginate, validator.ValidateStruct(struct {
		Status string `validate:"omitempty,oneof=1 2"`
	}{Status: filters["status"]}))

	if len(errValidate) > 0 {
		return http.ResponseUnprocessableEntity(context, errValidate)
	}

	list, err := c.u.List(context.Request().Context(), filters, sorts, page, limit)

	if err != nil {
		return http.ResponseError(context, err.GetMessage())
	}

	http.SetHeaderCountAndTotal(context, strconv.Itoa(int(list.Paginate.Total)), strconv.Itoa(list.Paginate.TotalPages))

	if len(list.List) == 0 {
		return http.ResponseEmptyList(context)
	}

	return context.JSON(
		200,
		slices.Map[entity.Category, map[string]interface{}](list.List, func(category entity.Category) map[string]interface{} {
			return category.ToMap()
		}),
	)
}

func (c Controller) Get(context echo.Context) error {
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

func (c Controller) Store(context echo.Context) error {
	body, err := http.ParseFormData[CategoryFormData](context)
	if err != nil {
		c.l.Error(err)
		return http.ResponseBadRequest(context, err.Error())
	}

	if errVald := validator.ValidateStruct(body); len(errVald) > 0 {
		return http.ResponseUnprocessableEntity(context, errVald)
	}

	errTrans := c.db.Begin()

	if errTrans != nil {
		c.l.Error(errTrans)
		return http.ResponseError(context, errTrans.Error())
	}
	cat, errCrt := c.u.Store(context.Request().Context(), body)

	if errCrt != nil {
		errTrans = c.db.Rollback()
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

	errTrans = c.db.Commit()
	if errTrans != nil {
		c.l.Error(errTrans)
		return http.ResponseError(context, errTrans.Error())
	}

	return http.ResponseOk(context, (*cat).ToMap())
}

func (c Controller) Update(context echo.Context) error {
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

	errTrans := c.db.Begin()
	if errTrans != nil {
		c.l.Error(errTrans)
		return http.ResponseError(context, errTrans.Error())
	}
	errUpdate := c.u.Update(context.Request().Context(), category, body)

	if errUpdate != nil {
		errTrans = c.db.Rollback()
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
	errTrans = c.db.Commit()
	if errTrans != nil {
		c.l.Error(errTrans)
		return http.ResponseError(context, errTrans.Error())
	}

	return http.ResponseOk(context, category.ToMap())

}

func (c Controller) Delete(context echo.Context) error {
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

	errTrans := c.db.Begin()
	if errTrans != nil {
		c.l.Error(errTrans)
		return http.ResponseError(context, errTrans.Error())
	}

	errDel := c.u.Delete(context.Request().Context(), category)

	if errDel != nil {
		errTrans = c.db.Rollback()
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
	errTrans = c.db.Commit()
	if errTrans != nil {
		c.l.Error(errTrans)
		return http.ResponseError(context, errTrans.Error())
	}

	return http.ResponseOk(context, true)
}

func (c Controller) RegisterRoute(group *echo.Group) {
	g := group.Group("/categories")
	g.GET("", c.List)
	g.POST("", c.Store)
	g.GET("/:id", c.Get)
	g.PUT("/:id", c.Update)
	g.DELETE("/:id", c.Delete)
}

func NewHttpV1Controller(
	l log.Contract,
	db database.Contract,
	u usecase.CategoryUseCaseContract,
) *Controller {
	return &Controller{l: l, u: u, db: db}
}
