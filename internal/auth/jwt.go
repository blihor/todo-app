package auth

import (
	"time"

	"github.com/golang-jwt/jwt/v5"
	"go.mongodb.org/mongo-driver/v2/bson"
)

type jwtClaims struct {
	UserID string `json:"user_id"`
	jwt.RegisteredClaims
}

func GenerateToken(userID bson.ObjectID, secretJwt []byte) (string, error) {
	claims := &jwtClaims{
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

	tokenString, err := token.SignedString(secretJwt)
	if err != nil {
		return "", err
	}

	return tokenString, nil
}
