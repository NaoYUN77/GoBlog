package jwt

//定义获取的参数

import (
	"Blog/models"
	"errors"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

func GenJWT(username string) (token string, err error) {
	var claims = models.Claims{
		Username: username,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(models.TokenExpireDuration)),
			Issuer:    "Blog",
		},
	}
	token, err = jwt.NewWithClaims(jwt.SigningMethodHS256, claims).SignedString(models.Sercret)
	if err != nil {
		return "", errors.New("jwt signing error")
	}
	return token, nil
}

func ParseJWT(token string) (*models.Claims, error) {
	var claims = new(models.Claims)
	tokenstruct, err := jwt.ParseWithClaims(token, claims, func(token *jwt.Token) (interface{}, error) {
		//需要Secret进行解码
		return models.Sercret, nil
	})
	if err != nil {
		return nil, errors.New("jwt parse error")
	}
	//验证token是否有效
	if !tokenstruct.Valid {
		return nil, errors.New("token invalid")
	}
	return claims, nil
}
