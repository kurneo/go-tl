package usecase

import (
	"context"
	"github.com/kurneo/go-template/config"
	"github.com/kurneo/go-template/internal/auth/entities"
	"github.com/kurneo/go-template/pkg/error"
	"github.com/kurneo/go-template/pkg/logger"
	"github.com/kurneo/go-template/pkg/support/jwt"
)

type tokenManager[T jwt.SubType] struct {
	cfg config.JWT
	l   logger.Contract
	t   jwt.TokenManager[T]
}

func (t tokenManager[T]) CreateToken(u *entities.User[T]) (*jwt.AccessToken[T], error.Contract) {
	token, err := t.t.CreateToken(u.ID)
	if err != nil {
		return nil, error.NewDomain(err)
	}
	return token, nil
}

func (t tokenManager[T]) CheckToken(ctx context.Context, token string) (*jwt.AccessToken[T], error.Contract) {
	return nil, nil
}

func (t tokenManager[T]) InvalidToken(ctx context.Context, token *jwt.AccessToken[T]) error.Contract {
	err := t.t.InvalidToken(ctx, token)
	if err != nil {
		return error.NewDomain(err)
	}
	return nil
}

func NewTokenManager[T jwt.SubType](cfg config.JWT, l logger.Contract) TokenManagerContract[T] {
	return &tokenManager[T]{
		cfg: cfg,
		l:   l,
	}
}
