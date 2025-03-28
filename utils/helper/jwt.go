package helper

import (
	"errors"
	"os"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

var jwt_secret = []byte(os.Getenv("JWT_SECRET"))

type JWTClaims struct {
	UserID     uint   `json:"user_id"`
	Email      string `json:"email"`
	IsVerified bool   `json:"is_verified"`
	jwt.RegisteredClaims
}

func GenerateJWTRegister(userid uint, email string) (string, error) {
	claims := JWTClaims{
		UserID:     userid,
		Email:      email,
		IsVerified: false,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(15 * time.Minute)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(jwt_secret)
}

func GenerateJWTLogin(userid uint, email string, verified bool) (string, error) {
	claims := JWTClaims{
		UserID:     userid,
		Email:      email,
		IsVerified: verified,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(24 * time.Hour)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(jwt_secret)
}

func ParseJWT(tokenstring string) (*JWTClaims, error) {
	token, err := jwt.ParseWithClaims(tokenstring, &JWTClaims{}, func(t *jwt.Token) (interface{}, error) {
		return jwt_secret, nil
	})

	if err != nil {
		return nil, err
	}

	claims, ok := token.Claims.(*JWTClaims)
	if !ok || !token.Valid {
		return nil, errors.New("invalid token")
	}

	return claims, nil
}
