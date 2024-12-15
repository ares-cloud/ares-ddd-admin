package jwt

import (
	"errors"
	"fmt"

	"github.com/ares-cloud/ares-ddd-admin/pkg/constant"
	"github.com/ares-cloud/ares-ddd-admin/pkg/token"
	"github.com/cloudwego/hertz/pkg/app"
)

var (
	ErrNotFound = errors.New("not found from context")
)

// ParseToken 解析TOKEN
func ParseToken(c *app.RequestContext) (*token.AccessToken, error) {
	val, ok := c.Get(constant.KeyAccessToken)
	if !ok {
		return nil, ErrNotFound
	}

	accessToken, ok := val.(token.AccessToken)
	if !ok {
		return nil, fmt.Errorf("parse error")
	}

	return &accessToken, nil
}
