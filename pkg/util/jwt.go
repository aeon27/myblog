package util

import (
	"time"

	"github.com/aeon27/myblog/pkg/setting"
	"github.com/dgrijalva/jwt-go"
)

var jwtSecret = []byte(setting.JwtSecret)

type Claims struct {
	Username string `json:"username"`
	Password string `json:"password"`
	jwt.StandardClaims
}

func GenerateToken(username, password string) (string, error) {
	// token有效期为3H
	expireTime := time.Now().Add(3 * time.Hour)

	// 声明 payload
	claims := Claims{
		username,
		password,
		jwt.StandardClaims{
			ExpiresAt: expireTime.Unix(),
			Issuer:    "my-blog",
		},
	}

	// 生成 JWT 结构体，采用 HMAC-SHA256 方法
	tokenClaims := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	// 运用给定密钥生成加密后的 JWT，格式为： header.payload.signature
	token, err := tokenClaims.SignedString(jwtSecret)

	return token, err
}

func ParseToken(token string) (*Claims, error) {
	tokenClaims, err := jwt.ParseWithClaims(token, &Claims{}, func(token *jwt.Token) (interface{}, error) {
		return jwtSecret, nil
	})

	if tokenClaims != nil {
		if claims, ok := tokenClaims.Claims.(*Claims); ok && tokenClaims.Valid {
			return claims, nil
		}
	}

	return nil, err
}
