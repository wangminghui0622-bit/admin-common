// Package jwtutil 签发与解析 admin-zero 的 JWT（HS256）。
package jwtutil

import (
	"time"

	"github.com/golang-jwt/jwt/v5"
)

type Claims struct {
	UserId   string `json:"userId"`
	RoleCode string `json:"roleCode"`
	Username string `json:"username"`
	jwt.RegisteredClaims
}

// NewToken 签发 token，返回 token 字符串与过期时间戳（秒）。
func NewToken(secret string, expireSec int64, userId, roleCode, username string) (string, int64, error) {
	expireAt := time.Now().Add(time.Duration(expireSec) * time.Second)
	claims := Claims{
		UserId:   userId,
		RoleCode: roleCode,
		Username: username,
		RegisteredClaims: jwt.RegisteredClaims{
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			ExpiresAt: jwt.NewNumericDate(expireAt),
		},
	}
	token, err := jwt.NewWithClaims(jwt.SigningMethodHS256, claims).SignedString([]byte(secret))
	if err != nil {
		return "", 0, err
	}
	return token, expireAt.Unix(), nil
}

// ParseToken 校验并解析 token。
func ParseToken(secret, tokenStr string) (*Claims, error) {
	claims := &Claims{}
	_, err := jwt.ParseWithClaims(tokenStr, claims, func(t *jwt.Token) (any, error) {
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, jwt.ErrSignatureInvalid
		}
		return []byte(secret), nil
	})
	if err != nil {
		return nil, err
	}
	return claims, nil
}
