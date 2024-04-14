package jwt

import (
	"banner/internal/domain"
	"github.com/golang-jwt/jwt/v5"
	"time"
)

func NewToken(user *domain.User, dur time.Duration, secret string) (string, error) {
	token := jwt.New(jwt.SigningMethodHS256)

	claims := token.Claims.(jwt.MapClaims)
	claims["login"] = user.Login
	claims["exp"] = time.Now().Add(dur).Unix()
	claims["is_admin"] = user.IsAdmin

	tokenString, err := token.SignedString([]byte(secret))
	if err != nil {
		return "", err
	}

	return tokenString, nil
}

func ParseToken(token, secret string) (jwt.MapClaims, error) {

	t, err := jwt.Parse(token, func(token *jwt.Token) (interface{}, error) {
		return []byte(secret), nil
	})
	if err != nil {
		return nil, err
	}
	if !t.Valid {
		return nil, domain.ErrInvalidToken
	}

	claims, ok := t.Claims.(jwt.MapClaims)
	if !ok {
		return nil, domain.ErrInvalidToken
	}

	return claims, nil
}
