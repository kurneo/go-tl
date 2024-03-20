package usecase

import (
	"context"
	"errors"
	"github.com/kurneo/go-template/internal/auth/domain"
	"github.com/kurneo/go-template/internal/auth/domain/entity"
	"github.com/kurneo/go-template/internal/auth/domain/repository"
	"github.com/kurneo/go-template/pkg/error"
	"github.com/kurneo/go-template/pkg/jwt"
	"time"
)

var (
	ErrEmailNotFound    = errors.New("email not found")
	ErrPasswordNotMatch = errors.New("password not match")
)

type UserUseCaseContract interface {
	Login(ctx context.Context, e, p string) (*jwt.AccessToken[int64], error.Contract)
	Logout(ctx context.Context, token *jwt.AccessToken[int64]) error.Contract
	GetProfile(ctx context.Context, sub int64) (*entity.User, error.Contract)
}

type UserUseCase struct {
	r repository.UserRepositoryContact
	t *domain.TokenManager
	p *domain.PasswordChecker
}

func (u UserUseCase) Login(ctx context.Context, e, p string) (*jwt.AccessToken[int64], error.Contract) {
	user, err := u.r.GetUser(ctx, e)
	if err != nil {
		return nil, err
	}

	if user == nil {
		return nil, error.NewDomain(ErrEmailNotFound)
	}

	if u.p.Check(user.Password, p) == false {
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

func (u UserUseCase) Logout(ctx context.Context, token *jwt.AccessToken[int64]) error.Contract {
	return u.t.InvalidToken(ctx, token)
}

func (u UserUseCase) GetProfile(ctx context.Context, sub int64) (*entity.User, error.Contract) {
	user, err := u.r.GetUserById(ctx, sub)
	if err != nil {
		return nil, err
	}
	return user, nil
}

func NewUserUseCase(
	r repository.UserRepositoryContact,
	t *domain.TokenManager,
	p *domain.PasswordChecker,
) UserUseCaseContract {
	return &UserUseCase{
		r: r,
		t: t,
		p: p,
	}
}
