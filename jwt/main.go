package main

import (
	"errors"
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

// 定义我们的密钥（在实际应用中应该从安全配置中获取）
var jwtKey = []byte("your_secret_key")

// 定义 Claims 结构体，可以添加自定义字段
type Claims struct {
	Username             string `json:"username"` // 自定义字段，例如用户名
	IsAdmin              bool   `json:"isAdmin"`  // 自定义字段，例如是否是管理员
	jwt.RegisteredClaims        // jwt 标准字段
}

func main() {
	// 创建一个新的 token
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, Claims{
		Username: "棉花糖",
		IsAdmin:  true,
		RegisteredClaims: jwt.RegisteredClaims{
			Issuer:    "码农夜读",
			Subject:   "棉花糖-JWT介绍",
			Audience:  []string{"programmer"},
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(24 * time.Hour)),
			IssuedAt:  jwt.NewNumericDate(time.Now().Add(1 * time.Hour)),
		},
	})

	// 使用密钥签名 token
	tokenString, err := token.SignedString(jwtKey)
	if err != nil {
		panic(err)
	}

	println("Generated Token:", tokenString)
	// 验证 token
	claims, err := ValidateToken(tokenString) // 故意篡改 token 来测试错误处理
	if err != nil {
		panic(err)
	}
	fmt.Println(claims.Username, claims.IsAdmin, claims.RegisteredClaims)
}

func ValidateToken(tokenString string) (*Claims, error) {
	claims := &Claims{}

	token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
		return jwtKey, nil
	})

	if err != nil {
		if errors.Is(err, jwt.ErrInvalidKey) {
			return nil, fmt.Errorf("jwt invalid")
		} else if errors.Is(err, jwt.ErrTokenMalformed) {
			return nil, fmt.Errorf("malformed token")
		} else if errors.Is(err, jwt.ErrTokenExpired) {
			return nil, fmt.Errorf("token has expired")
		}
		return nil, fmt.Errorf("couldn't parse token: %v", err)
	}

	if !token.Valid {
		return nil, fmt.Errorf("token is invalid")
	}

	return claims, nil
}
