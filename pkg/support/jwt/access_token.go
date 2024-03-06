package jwt

import "time"

type AccessToken[T SubType] struct {
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
