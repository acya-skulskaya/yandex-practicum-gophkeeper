package auth

import (
	"fmt"
	"time"

	"github.com/acya-skulskaya/yandex-practicum-gophkeeper/internal/models"
	"github.com/golang-jwt/jwt/v5"
)

type claims struct {
	jwt.RegisteredClaims
	UserID uint
}

func (s *Service) Token(user models.User) (token string, err error) {
	token, err = BuildJWTString(s.SecretKey, user.ID)
	if err != nil {
		return "", fmt.Errorf("could not create token string: %w", err)
	}

	return token, nil
}

// BuildJWTString создаёт токен и возвращает его в виде строки.
func BuildJWTString(secretKey string, userID uint) (string, error) {
	// создаём новый токен с алгоритмом подписи HS256 и утверждениями — claims
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims{
		RegisteredClaims: jwt.RegisteredClaims{
			// когда создан токен
			IssuedAt: jwt.NewNumericDate(time.Now()),
		},
		UserID: userID,
	})

	// создаём строку токена
	tokenString, err := token.SignedString([]byte(secretKey))
	if err != nil {
		return "", fmt.Errorf("could not get signed string: %w", err)
	}

	// возвращаем строку токена
	return tokenString, nil
}

func GetUserIDFromToken(secretKey string, tokenString string) (userID uint, err error) {
	claims := &claims{}
	token, err := jwt.ParseWithClaims(tokenString, claims,
		func(t *jwt.Token) (interface{}, error) {
			if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fmt.Errorf("unexpected signing method: %s", t.Header["alg"])
			}
			return []byte(secretKey), nil
		})
	if err != nil {
		return 0, ErrTokenIsNotValid
	}

	if !token.Valid {
		return 0, ErrTokenIsNotValid
	}

	return claims.UserID, nil
}
