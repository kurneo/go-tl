package finder

import (
	"github.com/kurneo/go-template/cmd/app"
	"github.com/kurneo/go-template/internal/finder/datasource"
	v1 "github.com/kurneo/go-template/internal/finder/transport/v1"
	"github.com/kurneo/go-template/internal/finder/usecase"
)

func RegisterModule(app app.Contract) error {
	lg := app.GetLogger()
	cfg := app.GetConfig()

	r, err := datasource.New(cfg.Storage, lg)

	if err != nil {
		return err
	}

	v1.New(app, usecase.New(lg, r))
	return nil
}
