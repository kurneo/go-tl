package usecase

import (
	"context"
	"errors"
	"github.com/kurneo/go-template/internal/category/entities"
	"github.com/kurneo/go-template/pkg/error"
	"github.com/kurneo/go-template/pkg/logger"
	"github.com/kurneo/go-template/pkg/support/paginate"
	"time"
)

var (
	errDefaultCatMustPublish  = errors.New("default category must be published")
	errDefaultCatMustSet      = errors.New("default category must be set")
	errCannotDeleteDefaultCat = errors.New("cannot delete default category")
)

type catUseCase struct {
	l logger.Contract
	r CategoryRepositoryContract
}

func (c catUseCase) List(
	ctx context.Context,
	filter map[string]string,
	sort map[string]string,
	page,
	perPage int,
) ([]entities.Category, *paginate.Paginator, error.Contract) {
	return c.r.List(ctx, filter, sort, page, perPage)
}

func (c catUseCase) Store(ctx context.Context, dto CategoryDTO) (*entities.Category, error.Contract) {
	createdAt := time.Now()
	updatedAt := time.Now()
	cat := &entities.Category{
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

func (c catUseCase) Get(ctx context.Context, id int) (*entities.Category, error.Contract) {
	return c.r.Get(ctx, id)
}

func (c catUseCase) Update(ctx context.Context, cat *entities.Category, dto CategoryDTO) error.Contract {
	if cat.IsDefault && dto.GetStatus() == entities.StatusDraft {
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

func (c catUseCase) Delete(ctx context.Context, cat *entities.Category) error.Contract {
	if cat.IsDefault {
		return error.NewDomain(errCannotDeleteDefaultCat)
	}
	return c.r.Delete(ctx, cat)
}

func New(l logger.Contract, r CategoryRepositoryContract) CategoryUseCaseContract {
	return &catUseCase{
		l: l,
		r: r,
	}
}
