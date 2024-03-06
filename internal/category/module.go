package category

import (
	"github.com/kurneo/go-template/cmd/app"
	"github.com/kurneo/go-template/internal/category/datasource"
	v1 "github.com/kurneo/go-template/internal/category/transport/http/v1"
	"github.com/kurneo/go-template/internal/category/usecase"
)

func RegisterModule(app app.Contract) error {
	lg := app.GetLogger()
	useCase := usecase.New(lg, datasource.NewCatRepo(lg, app.GetDB()))
	v1.New(app, useCase)
	return nil
}
