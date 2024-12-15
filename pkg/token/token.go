package token

import (
	"encoding/json"
	"errors"
	"github.com/golang-jwt/jwt/v5"
	"time"
)

var (
	ErrUnknown            = errors.New("couldn't handle this token")
	ErrMalformed          = errors.New("that's not even a token")
	ErrExpiredOrNotActive = errors.New("token is either expired or not active yet")
	ErrNotStandardClaims  = errors.New("claims not standard")
	ErrCannotParseSubject = errors.New("cannot parse subject")
)

// DefToken 默认的token实现
type DefToken struct {
	issuer            string
	signingKey        string
	expirationToken   int64
	expirationRefresh int64
}

func Def() *DefToken {
	return &DefToken{issuer: "gd-dev", signingKey: "uAYnaSgAiYzAiGwLFe", expirationToken: 360000, expirationRefresh: 720000}
}

func NewDefToken(issuer, signingKey string, expirationToken, expirationRefresh int64) *DefToken {
	return &DefToken{issuer: issuer, signingKey: signingKey, expirationToken: expirationToken, expirationRefresh: expirationRefresh}
}

// GenerateToken 生成令牌
func (to *DefToken) GenerateToken(userId string, data interface{}) (*Token, error) {
	bytes, err := json.Marshal(data)
	if err != nil {
		return nil, err
	}
	claims := &jwt.RegisteredClaims{
		Issuer: to.issuer,
		IssuedAt: &jwt.NumericDate{
			Time: time.Now(),
		},
		ExpiresAt: &jwt.NumericDate{
			Time: time.Now().Add(time.Second * time.Duration(to.expirationToken)),
		},
		NotBefore: &jwt.NumericDate{
			Time: time.Now(),
		},
		Subject: string(bytes),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS512, claims)
	ss, err := token.SignedString([]byte(to.signingKey))
	if err != nil {
		return nil, err
	}
	refToken, i, err := to.GenerateRefToken(userId, data)
	if err != nil {
		return nil, err
	}
	return &Token{
		AccessToken:           ss,
		ExpiresIn:             to.expirationToken,
		RefreshToken:          refToken,
		RefreshTokenExpiresIn: i,
	}, nil
}

// GenerateRefToken 生成令牌
func (to *DefToken) GenerateRefToken(_ string, data interface{}) (string, int64, error) {
	bytes, err := json.Marshal(data)
	if err != nil {
		return "", 0, err
	}
	claims := &jwt.RegisteredClaims{
		Issuer: to.issuer,
		IssuedAt: &jwt.NumericDate{
			Time: time.Now(),
		},
		ExpiresAt: &jwt.NumericDate{
			Time: time.Now().Add(time.Second * time.Duration(to.expirationRefresh)),
		},
		NotBefore: &jwt.NumericDate{
			Time: time.Now(),
		},
		Subject: string(bytes),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS512, claims)
	ss, err := token.SignedString([]byte(to.signingKey))
	if err != nil {
		return "", 0, err
	}

	return ss, to.expirationRefresh, nil
}

// Verify 验证令牌
func (to *DefToken) Verify(token string, data interface{}) error {
	t, err := jwt.ParseWithClaims(token, &jwt.RegisteredClaims{}, func(*jwt.Token) (interface{}, error) {
		return []byte(to.signingKey), nil
	})

	// 检查解析过程中的错误
	if err != nil {
		if errors.Is(err, jwt.ErrTokenMalformed) {
			return ErrMalformed
		} else if errors.Is(err, jwt.ErrTokenExpired) || errors.Is(err, jwt.ErrTokenNotValidYet) {
			return ErrExpiredOrNotActive
		} else {
			return ErrUnknown
		}
	}

	// 检查令牌是否有效
	if t == nil || !t.Valid {
		return ErrUnknown
	}

	// 有效时解析数据
	if data != nil {
		if err := to.parse(t, data); err != nil {
			return err
		}
	}

	return nil
}

func (to *DefToken) DelToken(token string) error {
	return nil
}
func (to *DefToken) DelUserToken(userID string) error {
	return nil
}

func (to *DefToken) parse(t *jwt.Token, data interface{}) error {
	clm, ok := t.Claims.(*jwt.RegisteredClaims)
	if !ok {
		return ErrNotStandardClaims
	}

	err := json.Unmarshal([]byte(clm.Subject), data)
	if err != nil {
		return ErrCannotParseSubject
	}

	return nil
}
