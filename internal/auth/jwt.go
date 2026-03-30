package auth

import (
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"go.mongodb.org/mongo-driver/v2/bson"
)

type JwtClaims struct {
	UserID string `json:"user_id"`
	jwt.RegisteredClaims
}

func (s *service) GenerateToken(userID bson.ObjectID) (string, error) {
	claims := &JwtClaims{
		UserID: userID.Hex(),
		RegisteredClaims: jwt.RegisteredClaims{
			Issuer:  "todo-app",
			Subject: "login",
			ExpiresAt: &jwt.NumericDate{
				Time: time.Now().Add(time.Hour * 12),
			},
			IssuedAt: &jwt.NumericDate{
				Time: time.Now(),
			},
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	tokenString, err := token.SignedString([]byte(s.secretJwt))
	if err != nil {
		return "", err
	}

	return tokenString, nil
}

func (s *service) ValidateToken(tokenString string) (*JwtClaims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &JwtClaims{}, func(t *jwt.Token) (any, error) {
		return []byte(s.secretJwt), nil
	}, jwt.WithValidMethods([]string{jwt.SigningMethodHS256.Alg()}))
	if err != nil {
		return nil, err
	}

	if claims, ok := token.Claims.(*JwtClaims); ok {
		return claims, nil
	} else {
		return nil, fmt.Errorf("Unknown claims type")
	}
}
