package utils

import (
	"errors"
	"log"
	"time"

	"github.com/golang-jwt/jwt"
)

func GenerateToken(role string, userid int, username string) (tokens string, err error) {
	// 生成token
	token := jwt.New(jwt.SigningMethodHS256)

	// 设置令牌的声明（claims）
	claims := token.Claims.(jwt.MapClaims)
	claims["sub"] = role
	claims["iss"] = "envoyadmin"
	claims["userId"] = userid
	claims["username"] = username
	claims["exp"] = time.Now().Add(time.Hour * 6).Unix() // 设置令牌过期时间
	secret := []byte("envoyadmin.wangpc")
	tokens, err = token.SignedString(secret)
	if err != nil {
		log.Println("令牌生成失败:", err)
		return "", err
	}
	return tokens, nil
}

func ParToken(tokenString string) (claims jwt.MapClaims, err error) {
	// 解析token
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		return []byte("envoyadmin.wangpc"), nil
	})
	if err != nil {
		return nil, err
	}
	if err != nil {
		log.Println("令牌解码失败:", err)
		return
	}

	// 检查令牌是否有效
	if token.Valid {
		// 令牌有效，可以访问声明（claims）数据
		claims := token.Claims.(jwt.MapClaims)
		return claims, nil
	} else if ve, ok := err.(*jwt.ValidationError); ok {
		// 令牌无效，根据 ValidationError 类型进行处理
		if ve.Errors&jwt.ValidationErrorMalformed != 0 {
			return nil, errors.New("令牌格式错误")
		} else if ve.Errors&(jwt.ValidationErrorExpired|jwt.ValidationErrorNotValidYet) != 0 {
			return nil, errors.New("令牌已过期")
		} else {
			return nil, errors.New("令牌无效")
		}
	} else {
		return nil, errors.New("令牌无效")
	}
}
