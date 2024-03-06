package usecase

import (
	"context"
	"github.com/kurneo/go-template/internal/auth/entities"
	"github.com/kurneo/go-template/pkg/error"
	"github.com/kurneo/go-template/pkg/support/jwt"
)

type (
	UserUseCaseContract[T jwt.SubType] interface {
		Login(ctx context.Context, e, p string) (*jwt.AccessToken[T], error.Contract)
		Logout(ctx context.Context, token *jwt.AccessToken[T]) error.Contract
		GetProfile(ctx context.Context, sub T) (*entities.User[T], error.Contract)
	}

	UserRepositoryContract[T jwt.SubType] interface {
		GetUser(ctx context.Context, e string) (*entities.User[T], error.Contract)
		GetUserById(ctx context.Context, id T) (*entities.User[T], error.Contract)
		UpdateLastLoginTime(ctx context.Context, u *entities.User[T]) error.Contract
	}

	TokenManagerContract[T jwt.SubType] interface {
		CreateToken(u *entities.User[T]) (*jwt.AccessToken[T], error.Contract)
		CheckToken(ctx context.Context, t string) (*jwt.AccessToken[T], error.Contract)
		InvalidToken(ctx context.Context, t *jwt.AccessToken[T]) error.Contract
	}
)
