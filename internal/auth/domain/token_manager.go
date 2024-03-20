package domain

import (
	"context"
	"github.com/kurneo/go-template/internal/auth/domain/entity"
	"github.com/kurneo/go-template/pkg/error"
	"github.com/kurneo/go-template/pkg/jwt"
)

type TokenManager struct {
	t *jwt.TokenManager[int64]
}

func (t TokenManager) CreateToken(u *entity.User) (*jwt.AccessToken[int64], error.Contract) {
	token, err := t.t.CreateToken(u.ID)
	if err != nil {
		return nil, error.NewDomain(err)
	}
	return token, nil
}

func (t TokenManager) CheckToken(ctx context.Context, token string) (*jwt.AccessToken[int64], error.Contract) {
	accessToken, err := t.t.CheckToken(ctx, token)
	if err != nil {
		return nil, error.NewDomain(err)
	}
	return accessToken, nil
}

func (t TokenManager) InvalidToken(ctx context.Context, token *jwt.AccessToken[int64]) error.Contract {
	err := t.t.InvalidToken(ctx, token)
	if err != nil {
		return error.NewDomain(err)
	}
	return nil
}

func NewTokenManager(tm *jwt.TokenManager[int64]) *TokenManager {
	return &TokenManager{
		t: tm,
	}
}
