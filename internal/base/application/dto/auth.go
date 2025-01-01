package dto

import "github.com/ares-cloud/ares-ddd-admin/pkg/token"

type Token struct {
}

type AuthDto struct {
	AccessToken           string `json:"access_token"`
	ExpiresIn             int64  `json:"expires_in"`
	RefreshToken          string `json:"refresh_token"`
	RefreshTokenExpiresIn int64  `json:"refresh_token_expires_in"`
}

func ToAuthDto(t *token.Token) *AuthDto {
	return &AuthDto{
		AccessToken:           t.AccessToken,
		ExpiresIn:             t.ExpiresIn,
		RefreshToken:          t.RefreshToken,
		RefreshTokenExpiresIn: t.RefreshTokenExpiresIn,
	}
}

type CaptchaDto struct {
	Key   string `json:"key"`
	Image string `json:"image"`
}
