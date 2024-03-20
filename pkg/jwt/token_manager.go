package jwt

import (
	"context"
	"errors"
	"fmt"
	"github.com/golang-jwt/jwt/v4"
	"github.com/google/uuid"
	"github.com/kurneo/go-template/pkg/cache"
	"github.com/kurneo/go-template/pkg/support/repository"
	echoJwt "github.com/labstack/echo-jwt/v4"
	"github.com/spf13/viper"
	"math"
	"time"
)

type AccessToken[T repository.PrimaryKey] struct {
	Sub         T       `json:"-"`
	UUID        string  `json:"-"`
	AccessToken string  `json:"access_token"`
	ExpiredAt   int64   `json:"-"`
	ExpiredIn   float64 `json:"expired_in"`
	IssueAt     int64   `json:"-"`
}

func (t AccessToken[T]) IsExpired() bool {
	return t.ExpiredAt < time.Now().Unix()
}

type jwtMapClaims[T repository.PrimaryKey] struct {
	Sub  T       `json:"sub"`
	UUID string  `json:"uuid"`
	Iat  int64   `json:"iat"`
	Ein  float64 `json:"ein"`
	Eat  int64   `json:"Eat"`
	jwt.MapClaims
}

type TokenManager[T repository.PrimaryKey] struct {
	c cache.Contact
}

func (t TokenManager[T]) CreateToken(sub T) (*AccessToken[T], error) {
	secret, err := t.GetSecret()
	if err != nil {
		return nil, err
	}
	timeout, err := t.GetTimeout()
	if err != nil {
		return nil, err
	}

	iat := time.Now()
	eat := iat.Add(time.Minute * time.Duration(timeout))
	ein := math.Floor(eat.Sub(iat).Seconds())
	tokenJwt := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"sub":  sub,
		"uuid": uuid.New().String(),
		"iat":  iat.Unix(),
		"ein":  ein,
		"eat":  eat.Unix(),
	})
	token, err := tokenJwt.SignedString([]byte(secret))

	if err != nil {
		return nil, err
	}

	return &AccessToken[T]{
		AccessToken: token,
		ExpiredAt:   eat.Unix(),
		ExpiredIn:   ein,
		IssueAt:     iat.Unix(),
	}, nil
}

func (t TokenManager[T]) CheckToken(ctx context.Context, token string) (*AccessToken[T], error) {
	parsedToken, err := t.ParseToken(token)

	if err != nil {
		return nil, err
	}

	if parsedToken.IsExpired() {
		return nil, jwt.ErrTokenExpired
	}

	check, err := t.c.Get(ctx, parsedToken.UUID)
	if err != nil {
		return nil, err
	}

	if check != nil {
		return nil, jwt.ErrTokenMalformed
	}

	return parsedToken, nil
}

func (t TokenManager[T]) GetSecret() (string, error) {
	secret := viper.GetString("JWT_SECRET")

	if secret == "" {
		return "", errors.New("jwt secret mismatch")
	}

	return secret, nil
}

func (t TokenManager[T]) GetTimeout() (int64, error) {
	timeout := viper.GetInt("JWT_TOKEN_TIMEOUT")

	if timeout == 0 {
		return 0, errors.New("jwt token timeout mismatch")
	}

	return int64(timeout), nil
}

func (t TokenManager[T]) InvalidToken(ctx context.Context, token *AccessToken[T]) error {
	err := t.c.Set(ctx, token.UUID, true, time.Duration(token.ExpiredAt-time.Now().Unix())*time.Second)
	if err != nil {
		return err
	}
	return nil
}

func (t TokenManager[T]) ParseToken(token string) (*AccessToken[T], error) {
	var claim = &jwtMapClaims[T]{}
	tokenParsed, err := jwt.ParseWithClaims(
		token,
		claim,
		func(token *jwt.Token) (interface{}, error) {
			if token.Method.Alg() != echoJwt.AlgorithmHS256 {
				return nil, fmt.Errorf("unexpected jwt signing method=%v", token.Header["alg"])
			}
			secret, _ := t.GetSecret()
			return []byte(secret), nil
		},
	)

	if err != nil {
		if err.Error() == "token contains an invalid number of segments" {
			return nil, jwt.ErrTokenMalformed
		}
		return nil, err
	}

	if !tokenParsed.Valid {
		return nil, jwt.ErrTokenMalformed
	}

	return &AccessToken[T]{
		AccessToken: token,
		Sub:         claim.Sub,
		UUID:        claim.UUID,
		ExpiredAt:   claim.Eat,
		IssueAt:     claim.Iat,
		ExpiredIn:   claim.Ein,
	}, nil
}

func NewTokenManager[T repository.PrimaryKey](c cache.Contact) *TokenManager[T] {
	return &TokenManager[T]{
		c: c,
	}
}
