package token

import (
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
)

type IToken interface {
	GenerateToken(userID string, data interface{}) (*Token, error)
	Verify(token string, data interface{}) error
	DelToken(token string) error
	DelUserToken(userID string) error
}

// AccessToken //token
type AccessToken struct {
	UserId       string   `json:"userId"`                   // 刷新 token
	UserName     string   `json:"userName"`                 // 用户账号
	Platform     string   `json:"platform"`                 // 平台类型
	TenantId     string   `json:"tenantId"`                 //租户id
	AccessToken  string   `json:"access_token,omitempty"`   // 访问 token
	ExpiresAt    int64    `json:"expires_at,omitempty"`     // 过期时间
	RefreshToken string   `json:"refresh_token,omitempty"`  // 刷新 token
	RefExpiresAt int64    `json:"ref_expires_at,omitempty"` // refToken过期时间
	ServerCode   string   `json:"server_code"`              // 服务码
	Roles        []string `json:"roles"`                    // 角色CODE列表
}

func (a *AccessToken) MarshalBinary() (data []byte, err error) {
	return json.Marshal(a)
}

// 生成 Token 的 Hash 值
func generateTokenHash(token string) string {
	hash := sha256.New()
	hash.Write([]byte(token))
	return hex.EncodeToString(hash.Sum(nil))
}

type Token struct {
	AccessToken           string `json:"access_token"` //token
	ExpiresIn             int64  `json:"expires_in"`   //
	RefreshToken          string `json:"refresh_token"`
	RefreshTokenExpiresIn int64  `json:"refresh_token_expires_in"`
}
