package usecase

import (
	"context"
	"errors"
	"github.com/kurneo/go-template/internal/admin/entities"
	"github.com/kurneo/go-template/pkg/error"
	"github.com/kurneo/go-template/pkg/logger"
	"time"
)

var (
	ErrEmailNotFound    = errors.New("email not found")
	ErrPasswordNotMatch = errors.New("password not match")
)

type AuthUseCase struct {
	r AdminRepositoryContract
	t TokenManagerContract
	p entities.PasswordCheckerContract
	l logger.Contract
}

func (u AuthUseCase) Login(ctx context.Context, e, p string) (*entities.AdminAccessToken, error.Contract) {
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

	token, err := u.t.CreateToken(ctx, user)

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

func (u AuthUseCase) Logout(ctx context.Context, t string) error.Contract {
	return u.t.InvalidToken(ctx, t)
}

func (u AuthUseCase) GetProfile(ctx context.Context, t string) (*entities.Admin, error.Contract) {
	authToken, err := u.t.CheckToken(ctx, t)

	if err != nil {
		return nil, err
	}

	user, err := u.r.GetUserById(ctx, authToken.AdminID)

	if err != nil {
		return nil, err
	}
	return user, nil
}

func (u AuthUseCase) RefreshToken(ctx context.Context, t string) (*entities.AdminAccessToken, error.Contract) {
	return u.t.RefreshToken(ctx, t)
}

func New(r AdminRepositoryContract, l logger.Contract, t TokenManagerContract, p entities.PasswordCheckerContract) AdminUseCaseContract {
	return &AuthUseCase{
		r: r,
		l: l,
		t: t,
		p: p,
	}
}
