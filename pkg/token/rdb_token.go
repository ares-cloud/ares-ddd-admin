package token

import (
	"context"
	"encoding/json"
	"errors"
	"github.com/golang-jwt/jwt/v5"
	"github.com/redis/go-redis/v9"
	"time"
)

// RdbToken 默认的token实现
type RdbToken struct {
	issuer            string
	signingKey        string
	expirationToken   int64
	expirationRefresh int64
	rdb               *redis.Client
}

func NewRdbToken(rdb *redis.Client, issuer, signingKey string, expirationToken, expirationRefresh int64) IToken {
	return &RdbToken{rdb: rdb, issuer: issuer, signingKey: signingKey, expirationToken: expirationToken, expirationRefresh: expirationRefresh}
}

func (r *RdbToken) GenerateToken(userID string, data interface{}) (*Token, error) {
	accessToken, expiration, err := r.generateToken(data)
	if err != nil {
		return nil, err
	}
	refreshToken, refreshTokenExpiration, err := r.generateRefToken(data)
	if err != nil {
		return nil, err
	}
	ctx := context.Background()
	accessTokenHash := generateTokenHash(accessToken)
	refreshTokenHash := generateTokenHash(refreshToken)
	// 保存 Access Token Hash -> User ID
	err = r.rdb.Set(ctx, accessTokenHash, userID, time.Duration(expiration)*time.Second).Err()
	if err != nil {
		return nil, err
	}
	// 保存 Refresh Token Hash -> User ID
	err = r.rdb.Set(ctx, refreshTokenHash, userID, time.Duration(refreshTokenExpiration)*time.Second).Err() // Refresh Token 一般时间更长
	if err != nil {
		return nil, err
	}
	// 保存 User ID -> Token Hash (包含 Access Token 和 Refresh Token)
	err = r.rdb.SAdd(ctx, "user:"+userID, accessTokenHash, refreshTokenHash).Err()
	if err != nil {
		return nil, err
	}
	return &Token{
		AccessToken:           accessToken,
		ExpiresIn:             expiration,
		RefreshToken:          refreshToken,
		RefreshTokenExpiresIn: refreshTokenExpiration,
	}, nil
}

func (r *RdbToken) DelToken(token string) error {
	tokenHash := generateTokenHash(token)
	ctx := context.Background()
	// 根据 Access Token Hash 获取 User ID
	userID, err := r.rdb.Get(ctx, tokenHash).Result()
	if err != nil {
		return err
	}

	// 获取用户关联的所有 Token（Access 和 Refresh）
	tokenHashes, err := r.rdb.SMembers(ctx, "user:"+userID).Result()
	if err != nil {
		return err
	}

	// 删除每个 Token Hash
	for _, th := range tokenHashes {
		err = r.rdb.Del(ctx, th).Err()
		if err != nil {
			return err
		}
	}

	// 删除 User ID -> Token Hash 集合
	err = r.rdb.Del(ctx, "user:"+userID).Err()
	return err
}

func (r *RdbToken) DelUserToken(userID string) error {
	ctx := context.Background()
	// 获取用户关联的所有 Token Hash（Access 和 Refresh）
	tokenHashes, err := r.rdb.SMembers(ctx, "user:"+userID).Result()
	if err != nil {
		return err
	}

	// 删除每个 Token Hash
	for _, tokenHash := range tokenHashes {
		err = r.rdb.Del(ctx, tokenHash).Err()
		if err != nil {
			return err
		}
	}

	// 删除 User ID -> Token Hash 集合
	err = r.rdb.Del(ctx, "user:"+userID).Err()
	return err
}
func (r *RdbToken) Verify(token string, data interface{}) error {
	tokenHash := generateTokenHash(token)
	// 从 Redis 获取该 Token 是否存在
	userID, err := r.rdb.Get(context.Background(), tokenHash).Result()
	if errors.Is(err, redis.Nil) {
		// Redis 中不存在该 Token，返回错误
		return ErrExpiredOrNotActive
	} else if err != nil {
		return ErrUnknown
	}
	if userID == "" {
		return ErrExpiredOrNotActive
	}
	t, err := jwt.ParseWithClaims(token, &jwt.RegisteredClaims{}, func(*jwt.Token) (interface{}, error) {
		return []byte(r.signingKey), nil
	})
	// 无效时检查错误
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
		if err = r.parse(t, data); err != nil {
			return err
		}
	}

	return nil
}

// GenerateToken 生成令牌
func (r *RdbToken) generateToken(data interface{}) (string, int64, error) {
	bytes, err := json.Marshal(data)
	if err != nil {
		return "", 0, err
	}
	claims := &jwt.RegisteredClaims{
		Issuer: r.issuer,
		IssuedAt: &jwt.NumericDate{
			Time: time.Now(),
		},
		ExpiresAt: &jwt.NumericDate{
			Time: time.Now().Add(time.Second * time.Duration(r.expirationToken)),
		},
		NotBefore: &jwt.NumericDate{
			Time: time.Now(),
		},
		Subject: string(bytes),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS512, claims)
	ss, err := token.SignedString([]byte(r.signingKey))
	if err != nil {
		return "", 0, err
	}

	return ss, r.expirationToken, nil
}

// GenerateRefToken 生成令牌
func (r *RdbToken) generateRefToken(data interface{}) (string, int64, error) {
	bytes, err := json.Marshal(data)
	if err != nil {
		return "", 0, err
	}
	claims := &jwt.RegisteredClaims{
		Issuer: r.issuer,
		IssuedAt: &jwt.NumericDate{
			Time: time.Now(),
		},
		ExpiresAt: &jwt.NumericDate{
			Time: time.Now().Add(time.Second * time.Duration(r.expirationRefresh)),
		},
		NotBefore: &jwt.NumericDate{
			Time: time.Now(),
		},
		Subject: string(bytes),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS512, claims)
	ss, err := token.SignedString([]byte(r.signingKey))
	if err != nil {
		return "", 0, err
	}

	return ss, r.expirationRefresh, nil
}
func (r *RdbToken) parse(t *jwt.Token, data interface{}) error {
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
