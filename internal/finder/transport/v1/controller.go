package v1

import (
	"github.com/kurneo/go-template/cmd/app"
	"github.com/kurneo/go-template/internal/finder/entities"
	"github.com/kurneo/go-template/internal/finder/usecase"
	"github.com/kurneo/go-template/pkg/logger"
	"github.com/kurneo/go-template/pkg/support/http"
	"github.com/kurneo/go-template/pkg/support/slices"
	"github.com/kurneo/go-template/pkg/support/validator"
	"github.com/labstack/echo/v4"
)

type controller struct {
	a app.Contract
	l logger.Contract
	u usecase.FinderUseCaseContract
}

func (ctl controller) ListContents(context echo.Context) error {
	path := context.QueryParam("path")
	sortField := context.QueryParam("sort_field")
	sortDir := context.QueryParam("sort_dir")

	if sortField == "" {
		sortField = "name"
	}

	if sortDir == "" {
		sortDir = "asc"
	}

	files, directories, err := ctl.u.GetContents(path, sortField, sortDir)

	if err != nil {
		return http.ResponseError(context, err.GetMessage())
	} else {
		return http.ResponseOk(context, map[string]interface{}{
			"files": slices.Map[entities.File, map[string]interface{}](
				files,
				func(file entities.File) map[string]interface{} {
					return file.ToMap()
				},
			),
			"directories": slices.Map[entities.Directory, map[string]interface{}](
				directories,
				func(directory entities.Directory) map[string]interface{} {
					return directory.ToMap()
				},
			),
		})
	}
}

func (ctl controller) CreateDirectory(context echo.Context) error {
	body, errParseForm := http.ParseFormData[CreateDirFormData](context)

	if errParseForm != nil {
		ctl.l.Error(errParseForm)
		return http.ResponseBadRequest(context, errParseForm.Error())
	}

	if errorsValidate := validator.ValidateStruct(body); len(errorsValidate) > 0 {
		return http.ResponseUnprocessableEntity(context, errorsValidate)
	}

	if err := ctl.u.CreateDirectory(body); err != nil {
		if err.IsDomainError() {
			return http.ResponseBadRequest(context, err.GetMessage())
		}
		return http.ResponseError(context, err.GetMessage())
	}

	return http.ResponseOk(context, true)
}

func (ctl controller) Thumbnail(context echo.Context) error {
	filePath := context.QueryParam("path")

	file, err := ctl.u.Thumbnail(filePath)

	if err != nil {
		if err.IsDomainError() {
			return http.ResponseBadRequest(context, err.GetMessage())
		}
		return http.ResponseError(context, err.GetMessage())
	}

	return http.ResponseBlob(context, "image/jpg", file)
}

func (ctl controller) Preview(context echo.Context) error {
	filePath := context.QueryParam("path")

	file, contentType, err := ctl.u.Preview(filePath)

	if err != nil {
		if err.IsDomainError() {
			return http.ResponseBadRequest(context, err.GetMessage())
		}
		return http.ResponseError(context, err.GetMessage())
	}

	return http.ResponseBlob(context, contentType, file)
}

func (ctl controller) Upload(context echo.Context) error {
	path := context.FormValue("path")
	file, err := context.FormFile("file")

	if err != nil {
		ctl.l.Error(err)
		return http.ResponseError(context, err.Error())
	}

	if errU := ctl.u.Upload(path, file); errU != nil {
		if errU.IsDomainError() {
			return http.ResponseBadRequest(context, errU.GetMessage())
		}
		return http.ResponseError(context, errU.GetMessage())
	}
	return http.ResponseOk(context, true)
}

func (ctl controller) Rename(context echo.Context) error {
	body, errParseForm := http.ParseFormData[RenameFormData](context)
	if errParseForm != nil {
		ctl.l.Error(errParseForm)
		return http.ResponseBadRequest(context, errParseForm.Error())
	}

	if errorsValidate := validator.ValidateStruct(body); len(errorsValidate) > 0 {
		return http.ResponseUnprocessableEntity(context, errorsValidate)
	}

	if errRename := ctl.u.Rename(body); errRename != nil {
		if errRename.IsDomainError() {
			return http.ResponseBadRequest(context, errRename.GetMessage())
		}
		return http.ResponseError(context, errRename.GetMessage())
	}

	return context.JSON(200, true)
}

func (ctl controller) Delete(context echo.Context) error {
	body, errParseForm := http.ParseFormData[DeleteFormData](context)
	if errParseForm != nil {
		ctl.l.Error(errParseForm)
		return http.ResponseBadRequest(context, errParseForm.Error())
	}

	if errorsValidate := validator.ValidateStruct(body); len(errorsValidate) > 0 {
		return http.ResponseUnprocessableEntity(context, errorsValidate)
	}

	if errDelete := ctl.u.Delete(body); errDelete != nil {
		if errDelete.IsDomainError() {
			return http.ResponseBadRequest(context, errDelete.GetMessage())
		}
		return http.ResponseError(context, errDelete.GetMessage())
	}

	return http.ResponseOk(context, true)
}

func (ctl controller) Copy(context echo.Context) error {
	body, errParseForm := http.ParseFormData[CopyFormData](context)
	if errParseForm != nil {
		ctl.l.Error(errParseForm)
		return http.ResponseBadRequest(context, errParseForm.Error())
	}

	if errorsValidate := validator.ValidateStruct(body); len(errorsValidate) > 0 {
		return http.ResponseUnprocessableEntity(context, errorsValidate)
	}

	if errCopy := ctl.u.Copy(body); errCopy != nil {
		if errCopy.IsDomainError() {
			return http.ResponseBadRequest(context, errCopy.GetMessage())
		}
		return http.ResponseError(context, errCopy.GetMessage())
	}

	return http.ResponseOk(context, true)
}

func (ctl controller) Cut(context echo.Context) error {
	body, errParseForm := http.ParseFormData[CutFormData](context)
	if errParseForm != nil {
		ctl.l.Error(errParseForm)
		return http.ResponseBadRequest(context, errParseForm.Error())
	}

	if errorsValidate := validator.ValidateStruct(body); len(errorsValidate) > 0 {
		return http.ResponseUnprocessableEntity(context, errorsValidate)
	}

	if errCut := ctl.u.Cut(body); errCut != nil {
		if errCut.IsDomainError() {
			return http.ResponseBadRequest(context, errCut.GetMessage())
		}
		return http.ResponseError(context, errCut.GetMessage())
	}

	return http.ResponseOk(context, true)
}

func New(a app.Contract, u usecase.FinderUseCaseContract) {
	c := &controller{
		l: a.GetLogger(),
		a: a,
		u: u,
	}
	a.RegisterAdminV1Route(func(group *echo.Group, jwtMiddleware echo.MiddlewareFunc) {
		g := group.Group("/finder", jwtMiddleware)
		g.GET("/contents", c.ListContents)
		g.POST("/create-directory", c.CreateDirectory)
		g.GET("/thumbnail", c.Thumbnail)
		g.GET("/preview", c.Preview)
		g.POST("/upload", c.Upload)
		g.POST("/rename", c.Rename)
		g.POST("/delete", c.Delete)
		g.POST("/copy", c.Copy)
		g.POST("/cut", c.Cut)

		a.GetHttpHandler().Static("storage", "storage/app/public/files")
	})
}
