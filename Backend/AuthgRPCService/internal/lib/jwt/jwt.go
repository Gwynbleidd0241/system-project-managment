package jwt

import (
	"authService/internal/domain/models"
	"github.com/dgrijalva/jwt-go"
	"time"
)

func NewToken(user models.User, duration time.Duration) (string, error) {

	// TODO: Add normal secret key

	secret := "secretKEY"

	token := jwt.New(jwt.SigningMethodHS256)
	claims := token.Claims.(jwt.MapClaims)

	claims["user_id"] = user.ID
	claims["email"] = user.Email
	claims["exp"] = time.Now().Add(duration).Unix()

	tokenString, err := token.SignedString([]byte(secret))
	if err != nil {
		return "", err
	}

	return tokenString, nil
}
