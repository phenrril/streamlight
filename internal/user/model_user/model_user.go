package model_user

import "time"

type UserInfo struct {
	Sub           string    `json:"sub"`
	GivenName     string    `json:"given_name"`
	Nickname      string    `json:"nickname"`
	Name          string    `json:"name"`
	Picture       string    `json:"picture"`
	Locale        string    `json:"locale"`
	UpdatedAt     time.Time `json:"updated_at"`
	Email         string    `json:"email"`
	EmailVerified bool      `json:"email_verified"`
}
