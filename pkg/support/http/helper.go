package http

import (
	"github.com/labstack/echo/v4"
	"github.com/monoculum/formam/v3"
)

func ParseFormData[T any](context echo.Context, vars ...interface{}) (*T, error) {
	var body T
	tagName := "form"
	if len(vars) > 0 && vars[0] != nil {
		tagName = vars[0].(string)
	}

	formParams, err := context.FormParams()

	if err != nil {
		return nil, err
	}

	decoder := formam.NewDecoder(&formam.DecoderOptions{TagName: tagName})
	if errDecode := decoder.Decode(formParams, &body); errDecode != nil {
		return nil, errDecode
	}
	return &body, nil
}

func SetHeaderCountAndTotal(context echo.Context, total, totalPage string) {
	context.Response().Header().Set("x-Total-Count", total)
	context.Response().Header().Set("X-Total-Pages", totalPage)
}

func MergeErrorValidate(errors ...map[string][]string) map[string][]string {
	errs := make(map[string][]string)
	for _, err := range errors {
		if len(err) == 0 {
			continue
		}
		for field, errorDetail := range err {
			errs[field] = errorDetail
		}
	}
	return errs
}
