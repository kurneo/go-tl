package usecase

import (
	"context"
	"errors"
	"github.com/kurneo/go-template/internal/auth/entities"
	"github.com/kurneo/go-template/pkg/error"
	"github.com/kurneo/go-template/pkg/logger"
	"github.com/kurneo/go-template/pkg/support/jwt"
	"time"
)

var (
	ErrEmailNotFound    = errors.New("email not found")
	ErrPasswordNotMatch = errors.New("password not match")
)

type AuthUseCase[T jwt.SubType] struct {
	r UserRepositoryContract[T]
	t TokenManagerContract[T]
	p entities.PasswordCheckerContract
	l logger.Contract
}

func (u AuthUseCase[T]) Login(ctx context.Context, e, p string) (*jwt.AccessToken[T], error.Contract) {
	user, err := u.r.GetUser(ctx, e)
	if err != nil {
		return nil, err
	}

	if user == nil {
		return nil, error.NewDomain(ErrEmailNotFound)
	}

	if user.CheckPassword(p, u.p) == false {
		return nil, error.NewDomain(ErrPasswordNotMatch)
	}

	token, err := u.t.CreateToken(user)

	if err != nil {
		return nil, err
	}

	loginTime := time.Now()
	user.LastLoginAt = &loginTime
	err = u.r.UpdateLastLoginTime(ctx, user)

	if err != nil {
		return nil, err
	}

	return token, nil
}

func (u AuthUseCase[T]) Logout(ctx context.Context, token *jwt.AccessToken[T]) error.Contract {
	return u.t.InvalidToken(ctx, token)
}

func (u AuthUseCase[T]) GetProfile(ctx context.Context, sub T) (*entities.User[T], error.Contract) {
	user, err := u.r.GetUserById(ctx, sub)
	if err != nil {
		return nil, err
	}
	return user, nil
}

func New[T jwt.SubType](
	r UserRepositoryContract[T],
	l logger.Contract,
	t TokenManagerContract[T],
	p entities.PasswordCheckerContract,
) UserUseCaseContract[T] {
	return &AuthUseCase[T]{
		r: r,
		l: l,
		t: t,
		p: p,
	}
}
