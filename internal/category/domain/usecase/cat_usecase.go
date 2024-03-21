package usecase

import (
	"context"
	"errors"
	"github.com/kurneo/go-template/internal/category/domain/entity"
	"github.com/kurneo/go-template/internal/category/domain/repository"
	"github.com/kurneo/go-template/pkg/error"
	"github.com/kurneo/go-template/pkg/log"
	"github.com/kurneo/go-template/pkg/support/page_list"
	"time"
)

var (
	errDefaultCatMustPublish  = errors.New("default category must be published")
	errDefaultCatMustSet      = errors.New("default category must be set")
	errCannotDeleteDefaultCat = errors.New("cannot delete default category")
)

type CategoryUseCaseContract interface {
	List(ctx context.Context, filter map[string]string, sort map[string]string, page, perPage int) (*page_list.PageList[entity.Category], error.Contract)
	Store(ctx context.Context, dto CategoryDTO) (*entity.Category, error.Contract)
	Get(ctx context.Context, id int64) (*entity.Category, error.Contract)
	Update(ctx context.Context, cat *entity.Category, dto CategoryDTO) error.Contract
	Delete(ctx context.Context, cat *entity.Category) error.Contract
}

type CategoryDTO interface {
	GetName() string
	GetDescription() *string
	GetStatus() int
	GetIsDefault() bool
}

type CatUseCase struct {
	l log.Contract
	r repository.CategoryRepositoryContract
}

func (c CatUseCase) List(
	ctx context.Context,
	filter map[string]string,
	sort map[string]string,
	page,
	perPage int,
) (*page_list.PageList[entity.Category], error.Contract) {
	return c.r.List(ctx, filter, sort, page, perPage)
}

func (c CatUseCase) Store(ctx context.Context, dto CategoryDTO) (*entity.Category, error.Contract) {
	createdAt := time.Now()
	updatedAt := time.Now()
	cat := &entity.Category{
		Name:        dto.GetName(),
		Description: dto.GetDescription(),
		Status:      dto.GetStatus(),
		IsDefault:   dto.GetIsDefault(),
		CreatedAt:   &createdAt,
		UpdatedAt:   &updatedAt,
	}

	if cat.IsDefault && !cat.IsPublic() {
		return nil, error.NewDomain(errDefaultCatMustPublish)
	}

	err := c.r.Store(ctx, cat)

	if err != nil {
		return nil, err
	}

	if cat.IsDefault {
		err = c.r.UpdateDefault(ctx, cat)
		if err != nil {
			return nil, err
		}
	}

	return cat, nil
}

func (c CatUseCase) Get(ctx context.Context, id int64) (*entity.Category, error.Contract) {
	return c.r.Get(ctx, id)
}

func (c CatUseCase) Update(ctx context.Context, cat *entity.Category, dto CategoryDTO) error.Contract {
	if cat.IsDefault && dto.GetStatus() == entity.StatusDraft {
		return error.NewDomain(errDefaultCatMustPublish)
	}

	if cat.IsDefault && !dto.GetIsDefault() {
		return error.NewDomain(errDefaultCatMustSet)
	}

	updatedAt := time.Now()

	if dto.GetIsDefault() {
		err := c.r.UpdateDefault(ctx, cat)
		if err != nil {
			return err
		}
	}

	cat.Name = dto.GetName()
	cat.Description = dto.GetDescription()
	cat.Status = dto.GetStatus()
	cat.IsDefault = dto.GetIsDefault()
	cat.UpdatedAt = &updatedAt

	return c.r.Update(ctx, cat)
}

func (c CatUseCase) Delete(ctx context.Context, cat *entity.Category) error.Contract {
	if cat.IsDefault {
		return error.NewDomain(errCannotDeleteDefaultCat)
	}
	return c.r.Delete(ctx, cat)
}

func NewCatUseCase(r repository.CategoryRepositoryContract) CategoryUseCaseContract {
	return &CatUseCase{
		r: r,
	}
}
