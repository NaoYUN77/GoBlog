package models

import (
	"time"

	jwt "github.com/golang-jwt/jwt/v5"
)

type Claims struct {
	//自定义Playload
	Username string
	jwt.RegisteredClaims
}

// 定义加密盐
var Sercret = []byte("lolololo")

// 定义access token过期时间
const TokenExpireDuration = time.Minute * 10
