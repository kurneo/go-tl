package entities

import "time"

type AdminAccessToken struct {
	ID          int        `json:"-"`
	AccessToken string     `json:"access_token"`
	ExpiredAt   time.Time  `json:"-"`
	ExpiredIn   float64    `json:"expired_in"`
	AdminID     int        `json:"-"`
	CreatedAt   *time.Time `json:"-"`
	Admin       *Admin     `json:"-"`
}

func (t AdminAccessToken) IsExpired() bool {
	return t.ExpiredAt.Unix() < time.Now().Unix()
}

func (t AdminAccessToken) ToMap() map[string]interface{} {
	return map[string]interface{}{
		"access_token": t.AccessToken,
		"expired_in":   t.ExpiredIn,
	}
}
