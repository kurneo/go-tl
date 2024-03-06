package datasource

import (
	"context"
	"github.com/kurneo/go-template/internal/admin/datasource/models"
	"github.com/kurneo/go-template/internal/admin/entities"
	pkgError "github.com/kurneo/go-template/pkg/error"
	"github.com/kurneo/go-template/pkg/support/repository"
)

type TokenRepo struct {
	repository.Repository[models.AdminAccessToken, entities.AdminAccessToken]
}

func (repo TokenRepo) Create(ctx context.Context, token *entities.AdminAccessToken) pkgError.Contract {
	if err := repo.Insert(ctx, token); err != nil {
		return err
	}
	return nil
}

func (repo TokenRepo) Get(ctx context.Context, token string) (*entities.AdminAccessToken, pkgError.Contract) {
	t, err := repo.FirstBy(ctx, repository.Equal("token", token), repository.With("Admin"))
	if err != nil {
		return nil, err
	}
	return t, nil
}

func (repo TokenRepo) Invalid(ctx context.Context, token *entities.AdminAccessToken) pkgError.Contract {
	if err := repo.Delete(ctx, token); err != nil {
		return err
	}
	return nil
}
