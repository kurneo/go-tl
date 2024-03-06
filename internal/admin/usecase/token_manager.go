package usecase

import (
	"context"
	"errors"
	"github.com/golang-jwt/jwt/v4"
	"github.com/kurneo/go-template/config"
	"github.com/kurneo/go-template/internal/admin/entities"
	pkgError "github.com/kurneo/go-template/pkg/error"
	"github.com/kurneo/go-template/pkg/logger"
	"math"
	"strconv"
	"time"
)

type tokenManager struct {
	c config.JWT
	l logger.Contract
	r AdminTokenRepositoryContract
}

func (t tokenManager) CreateToken(ctx context.Context, u *entities.Admin) (*entities.AdminAccessToken, pkgError.Contract) {
	secret, timeout, err := t.getConfig()
	if err != nil {
		t.l.Error(err)
		return nil, pkgError.NewDomain(err)
	}
	createdAt := time.Now()
	expiredAt := time.Now().Add(time.Minute * time.Duration(timeout))
	expiredIn := math.Floor(expiredAt.Sub(time.Now()).Seconds())
	tokenJwt := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{"sub": u.ID, "expired_in": expiredIn, "expired_at": expiredAt})
	token, err := tokenJwt.SignedString([]byte(secret))

	if err != nil {
		t.l.Error(err)
		return nil, pkgError.NewDomain(err)
	}

	tokenEntity := &entities.AdminAccessToken{
		AccessToken: token,
		ExpiredAt:   expiredAt,
		AdminID:     u.ID,
		CreatedAt:   &createdAt,
	}

	errCreate := t.r.Create(ctx, tokenEntity)

	if errCreate != nil {
		return nil, errCreate
	}

	tokenEntity.ExpiredIn = expiredIn

	return tokenEntity, nil
}

func (t tokenManager) CheckToken(ctx context.Context, token string) (*entities.AdminAccessToken, pkgError.Contract) {
	authToken, err := t.r.Get(ctx, token)
	if err != nil {
		return nil, err
	}

	if authToken == nil {
		return nil, pkgError.NewDomain(jwt.ErrTokenNotValidYet)
	}

	if authToken.IsExpired() {
		return nil, pkgError.NewDomain(jwt.ErrTokenExpired)
	}

	return authToken, nil
}

func (t tokenManager) RefreshToken(ctx context.Context, token string) (*entities.AdminAccessToken, pkgError.Contract) {
	authToken, err := t.r.Get(ctx, token)
	if err != nil {
		return nil, err
	}

	newToken, err := t.CreateToken(ctx, authToken.Admin)

	if err != nil {
		return nil, err
	}

	err = t.r.Invalid(ctx, authToken)

	if err != nil {
		return nil, err
	}

	return newToken, nil
}

func (t tokenManager) getConfig() (string, int, error) {
	secret := t.c.Secret
	timeout := t.c.Timeout

	if secret == "" {
		return "", 0, errors.New("jwt secret mismatch")
	}

	if timeout == "" {
		timeout = "1440"
	}

	timeoutInt, err := strconv.Atoi(timeout)

	if err != nil {
		return "", 0, err
	}

	return secret, timeoutInt, nil
}

func (t tokenManager) InvalidToken(ctx context.Context, token string) pkgError.Contract {
	authToken, err := t.CheckToken(ctx, token)
	if err != nil {
		return err
	}

	if err = t.r.Invalid(ctx, authToken); err != nil {
		return err
	}

	return nil
}

func NewTokenManager(cfg config.JWT, log logger.Contract, repository AdminTokenRepositoryContract) TokenManagerContract {
	return &tokenManager{
		c: cfg,
		l: log,
		r: repository,
	}
}
