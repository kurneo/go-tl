package jwt

import (
	"context"
	"errors"
	"fmt"
	"github.com/golang-jwt/jwt/v4"
	"github.com/google/uuid"
	"github.com/kurneo/go-template/config"
	"github.com/kurneo/go-template/pkg/redis"
	echoJwt "github.com/labstack/echo-jwt/v4"
	"math"
	"strconv"
	"time"
)

type SubType = interface {
	int64 | string
}

type jwtMapClaims[T SubType] struct {
	Sub  T       `json:"sub"`
	UUID string  `json:"uuid"`
	Iat  int64   `json:"iat"`
	Ein  float64 `json:"ein"`
	Eat  int64   `json:"Eat"`
	jwt.MapClaims
}

type TokenManager[T SubType] struct {
	c config.JWT
	r redis.Contact
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

	_, err = t.r.Get(ctx, parsedToken.AccessToken)
	if err != nil {
		if err == redis.Nil {
			return parsedToken, nil
		} else {
			return nil, err
		}
	}

	return nil, jwt.ErrTokenMalformed
}

func (t TokenManager[T]) GetSecret() (string, error) {
	secret := t.c.Secret

	if secret == "" {
		return "", errors.New("jwt secret mismatch")
	}

	return secret, nil
}

func (t TokenManager[T]) GetTimeout() (int64, error) {
	timeout := t.c.Timeout

	if timeout == "" {
		return 0, errors.New("jwt token timeout mismatch")
	}

	timeoutInt, err := strconv.Atoi(timeout)

	if err != nil {
		return 0, err
	}

	return int64(timeoutInt), nil
}

func (t TokenManager[T]) InvalidToken(ctx context.Context, token *AccessToken[T]) error {
	err := t.r.Set(ctx, token.UUID, true, time.Duration(token.ExpiredAt-time.Now().Unix()))
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

func NewTokenManager[T SubType](cfg config.JWT, r redis.Contact) TokenManager[T] {
	return TokenManager[T]{
		c: cfg,
		r: r,
	}
}
