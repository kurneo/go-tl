package category

import (
	"github.com/google/wire"
	"github.com/kurneo/go-template/internal/category/data/datasource"
	"github.com/kurneo/go-template/internal/category/data/repository"
	domainRepository "github.com/kurneo/go-template/internal/category/domain/repository"
	"github.com/kurneo/go-template/internal/category/domain/usecase"
	v1 "github.com/kurneo/go-template/internal/category/transport/http/v1"
	"github.com/kurneo/go-template/pkg/database"
	"github.com/kurneo/go-template/pkg/log"
)

var WireSet = wire.NewSet(
	ResolveCatDatasource,
	ResolveCatRepo,
	ResolveCatUseCase,
	ResolveCatV1Controller,
)

func ResolveCatDatasource(db database.Contract) *datasource.CatDatasource {
	return datasource.NewCatDatasource(db)
}
func ResolveCatRepo(d *datasource.CatDatasource) domainRepository.CategoryRepositoryContract {
	return repository.NewCatRepo(d)
}

func ResolveCatUseCase(r domainRepository.CategoryRepositoryContract) usecase.CategoryUseCaseContract {
	return usecase.NewCatUseCase(r)
}

func ResolveCatV1Controller(
	l log.Contract,
	db database.Contract,
	u usecase.CategoryUseCaseContract,
) *v1.Controller {
	return v1.NewV1Controller(l, db, u)
}
