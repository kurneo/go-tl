package usecase

import (
	"context"
	"github.com/kurneo/go-template/internal/admin/entities"
	"github.com/kurneo/go-template/pkg/error"
)

type (
	AdminUseCaseContract interface {
		Login(ctx context.Context, e, p string) (*entities.AdminAccessToken, error.Contract)
		Logout(ctx context.Context, t string) error.Contract
		GetProfile(ctx context.Context, t string) (*entities.Admin, error.Contract)
		RefreshToken(ctx context.Context, t string) (*entities.AdminAccessToken, error.Contract)
	}

	AdminRepositoryContract interface {
		GetUser(ctx context.Context, e string) (*entities.Admin, error.Contract)
		GetUserById(ctx context.Context, id int) (*entities.Admin, error.Contract)
		UpdateLastLoginTime(ctx context.Context, u *entities.Admin) error.Contract
	}

	AdminTokenRepositoryContract interface {
		Create(ctx context.Context, t *entities.AdminAccessToken) error.Contract
		Get(ctx context.Context, t string) (*entities.AdminAccessToken, error.Contract)
		Invalid(ctx context.Context, t *entities.AdminAccessToken) error.Contract
	}

	TokenManagerContract interface {
		CreateToken(ctx context.Context, u *entities.Admin) (*entities.AdminAccessToken, error.Contract)
		CheckToken(ctx context.Context, t string) (*entities.AdminAccessToken, error.Contract)
		RefreshToken(ctx context.Context, t string) (*entities.AdminAccessToken, error.Contract)
		InvalidToken(ctx context.Context, t string) error.Contract
	}
)
