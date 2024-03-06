package http

import (
	"github.com/kurneo/go-template/config"
	"github.com/labstack/echo/v4"
	"net/http"
	"os"
	"reflect"
)

func ResponseUnprocessableEntity(context echo.Context, errors interface{}) error {
	return context.JSON(http.StatusUnprocessableEntity, map[string]interface{}{
		"message": "the given data was invalid",
		"errors":  errors,
	})
}

func ResponseEmptyList(context echo.Context) error {
	return context.JSON(http.StatusOK, make([]string, 0))
}

func ResponseOk(context echo.Context, data interface{}) error {
	return context.JSON(http.StatusOK, data)
}

func ResponseError(context echo.Context, vars ...interface{}) error {
	message := "internal server error"
	cfg, _ := config.NewConfig()
	if cfg.Debug && len(vars) > 0 && vars[0] != nil && reflect.ValueOf(vars[0]).Kind() == reflect.String {
		message = vars[0].(string)
	}
	return context.JSON(
		http.StatusInternalServerError,
		map[string]interface{}{"message": message},
	)
}

func ResponseBadRequest(context echo.Context, vars ...interface{}) error {
	message := "bad request"
	if len(vars) > 0 && vars[0] != nil && reflect.ValueOf(vars[0]).Kind() == reflect.String {
		message = vars[0].(string)
	}
	return context.JSON(
		http.StatusBadRequest,
		map[string]interface{}{"message": message},
	)
}

func ResponseUnauthorized(context echo.Context, vars ...interface{}) error {
	message := "unauthorized"
	if len(vars) > 0 && vars[0] != nil && reflect.ValueOf(vars[0]).Kind() == reflect.String {
		message = vars[0].(string)
	}
	return context.JSON(
		http.StatusUnauthorized,
		map[string]interface{}{"message": message},
	)
}

func ResponseBlob(context echo.Context, contentType string, b []byte) error {
	return context.Blob(http.StatusOK, contentType, b)
}

func ResponseSteam(context echo.Context, contentType string, file *os.File) error {
	return context.Stream(http.StatusOK, contentType, file)
}

func ResponseNotFound(context echo.Context, vars ...interface{}) error {
	message := "not found"
	if len(vars) > 0 && vars[0] != nil && reflect.ValueOf(vars[0]).Kind() == reflect.String {
		message = vars[0].(string)
	}

	return context.JSON(
		http.StatusNotFound,
		map[string]interface{}{"message": message},
	)
}

func ResponseNoContent(context echo.Context) error {
	return context.NoContent(http.StatusNoContent)
}
