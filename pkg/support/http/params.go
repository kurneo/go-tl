package http

import (
	"github.com/kurneo/go-template/pkg/support/validator"
	"github.com/labstack/echo/v4"
	"reflect"
	"strconv"
	"strings"
)

func GetPaginateParams(context echo.Context) (int, int, map[string][]string) {
	page := context.QueryParam("page")
	perPage := context.QueryParam("per_page")

	errorsValidate := validator.ValidateStruct(struct {
		Page    string `validate:"omitempty,numeric,gte=1" json:"page"`
		PerPage string `validate:"omitempty,numeric,gte=1" json:"per_page"`
	}{
		Page:    page,
		PerPage: perPage,
	})

	if len(errorsValidate) > 0 {
		return 0, 0, errorsValidate
	}

	if page == "" {
		page = "1"
	}

	if perPage == "" {
		perPage = "10"
	}

	intPage, _ := strconv.Atoi(page)
	intPerPage, _ := strconv.Atoi(perPage)

	return intPage, intPerPage, errorsValidate
}

func GetSortParams(context echo.Context) map[string]string {
	sort := context.QueryParam("sort")
	var sorts = map[string]string{}

	if sort == "" || sort == "-" {
		return sorts
	}

	split := strings.Split(sort, ",")
	for _, s := range split {
		dir := "asc"
		field := s
		if strings.HasPrefix(field, "-") {
			dir = "desc"
			field = field[1:]
		}
		sorts[field] = dir
	}

	return sorts
}

func GetFilterParams(context echo.Context, keys []string) map[string]string {
	filters := map[string]string{}
	query := context.Request().URL.Query()
	for _, key := range keys {
		filters[key] = ""
		filterKey := "filters[" + key + "]"
		if len(query[filterKey]) > 0 {
			filters[key] = query[filterKey][0]
		}
	}
	return filters
}

func GetIDRouteParam(context echo.Context, vars ...interface{}) (int, *map[string][]string) {
	idKey := "id"

	if len(vars) > 0 && vars[0] != nil && reflect.ValueOf(vars[0]).Kind() == reflect.String {
		idKey = vars[0].(string)
	}

	id := context.Param(idKey)

	errValidate := validator.ValidateValue(id, "numeric,gte=1")

	if len(errValidate) > 0 {
		return 0, &map[string][]string{
			idKey: errValidate,
		}
	}

	intID, _ := strconv.Atoi(id)

	return intID, nil
}
