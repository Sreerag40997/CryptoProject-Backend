package utils

import (
	"errors"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

// var secret = []byte("SUPER_SECRET")

type Claims struct {
	UserID uint   `json:"userID"`
	Email  string `json:"email"`
	Role   string `json:"role"`
	jwt.RegisteredClaims
}

// Token Generating
func generate(userID uint, role string, jwtSecret string,  ttl time.Duration) (string, error) {
	
	exp := time.Now().Add(ttl)
	claims := &Claims{
		UserID: userID,
		Role:   role,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(exp),
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	signed, err := token.SignedString([]byte(jwtSecret))
	return signed, err
}

// AccessToken
func GenerateAccess(userID uint, role string, jwtSecret string) (string, error) {
	return generate(userID, role, jwtSecret, 15*time.Minute)
}

// RefreshToken
func GenerateRefresh(userID uint, role string, jwtSecret string) (string, error) {
	return generate(userID, role, jwtSecret, 7*24*time.Hour)
}

// Parse
func Parse(token string, jwtSecret string) (*jwt.Token, error) {
	return jwt.Parse(token, func(t *jwt.Token) (interface{}, error) {
		return jwtSecret, nil
	})
}

// Token Verifiction Func
func VerifyToken(tokenStr string, jwtSecret string) (*Claims, error) {
	claims := &Claims{}
	token, err := jwt.ParseWithClaims(tokenStr, claims, func(t *jwt.Token) (any, error) {

		 if _,ok:=t.Method.(*jwt.SigningMethodHMAC);!ok{
				return nil,jwt.ErrSignatureInvalid
			 }

		return []byte(jwtSecret),nil
	})
	if err != nil {
		return nil, err
	}
	if !token.Valid {
		return nil, errors.New("invalid token")
	}
	return claims, nil
}