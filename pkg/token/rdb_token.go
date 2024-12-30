package token

import (
	"context"
	"encoding/json"
	"errors"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/redis/go-redis/v9"
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
	userKey := "user:auth:" + userID
	if err := r.rdb.Del(ctx, userKey).Err(); err != nil {
		return nil, err
	}
	accessTokenHash := generateTokenHash(accessToken)
	refreshTokenHash := generateTokenHash(refreshToken)
	if err = r.rdb.SAdd(ctx, userKey, accessTokenHash, refreshTokenHash).Err(); err != nil {
		return nil, err
	}
	if err = r.rdb.Set(ctx, "token:"+accessTokenHash, data, time.Duration(r.expirationToken)*time.Second).Err(); err != nil {
		return nil, err
	}
	if err = r.rdb.Set(ctx, "refresh_token:"+refreshTokenHash, data, time.Duration(r.expirationRefresh)*time.Second).Err(); err != nil {
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
	ctx := context.Background()
	tokenHash := generateTokenHash(token)
	tokenKey := "token:" + tokenHash

	// 1. 从Redis获取token数据
	val, err := r.rdb.Get(ctx, tokenKey).Result()
	if err != nil {
		if errors.Is(err, redis.Nil) {
			return nil // token已经不存在，直接返回
		}
		return err
	}

	// 2. 解析数据获取userID
	var data map[string]interface{}
	if err = json.Unmarshal([]byte(val), &data); err != nil {
		return err
	}
	userID, ok := data["user_id"].(string)
	if !ok {
		return errors.New("invalid token data")
	}

	// 3. 获取用户的所有token
	userKey := "user:auth:" + userID
	tokenHashes, err := r.rdb.SMembers(ctx, userKey).Result()
	if err != nil {
		return err
	}

	// 4. 使用管道批量删除
	pipe := r.rdb.Pipeline()
	for _, hash := range tokenHashes {
		pipe.Del(ctx, "token:"+hash)
		pipe.Del(ctx, "refresh_token:"+hash)
	}
	pipe.Del(ctx, userKey)

	// 5. 执行管道命令
	_, err = pipe.Exec(ctx)
	return err
}

func (r *RdbToken) DelUserToken(userID string) error {
	ctx := context.Background()
	userKey := "user:auth:" + userID

	// 1. 获取用户的所有token hash
	tokenHashes, err := r.rdb.SMembers(ctx, userKey).Result()
	if err != nil {
		return err
	}

	// 2. 删除每个token的详细信息
	pipe := r.rdb.Pipeline()
	for _, hash := range tokenHashes {
		pipe.Del(ctx, "token:"+hash)
		pipe.Del(ctx, "refresh_token:"+hash)
	}

	// 3. 删除用户token集合
	pipe.Del(ctx, userKey)

	// 4. 执行管道命令
	_, err = pipe.Exec(ctx)
	return err
}
func (r *RdbToken) Verify(token string, data interface{}) error {
	ctx := context.Background()
	tokenHash := generateTokenHash(token)
	tokenKey := "token:" + tokenHash

	// 从 Redis 获取该 Token 是否存在
	val, err := r.rdb.Get(ctx, tokenKey).Result()
	if errors.Is(err, redis.Nil) {
		// Redis 中不存在该 Token，返回错误
		return ErrExpiredOrNotActive
	} else if err != nil {
		return ErrUnknown
	}

	// 解析存储的数据
	if data != nil {
		if err = json.Unmarshal([]byte(val), data); err != nil {
			return ErrCannotParseSubject
		}
	}

	// 验证JWT token
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

	return nil
}

// GenerateToken 生成令牌
func (r *RdbToken) generateTokenWithExpiration(data interface{}, expiration int64) (string, int64, error) {
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
			Time: time.Now().Add(time.Second * time.Duration(expiration)),
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

	return ss, expiration, nil
}

func (r *RdbToken) generateToken(data interface{}) (string, int64, error) {
	return r.generateTokenWithExpiration(data, r.expirationToken)
}

func (r *RdbToken) generateRefToken(data interface{}) (string, int64, error) {
	return r.generateTokenWithExpiration(data, r.expirationRefresh)
}
