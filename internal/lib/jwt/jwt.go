package jwt

import (
	"time"

	"github.com/Goga211/go-auth/internal/domain/model"
	"github.com/golang-jwt/jwt/v5"
)

func NewToken(user model.User, app model.App, duration time.Duration) (string, error) {
	token := jwt.New(jwt.SigningMethodHS256)
	claims := token.Claims.(jwt.MapClaims)
	claims["uid"] = user.ID
	claims["exp"] = time.Now().Add(duration).Unix()
	claims["email"] = user.Email
	claims["appID"] = app.ID

	tokenString, err := token.SignedString([]byte(app.Secret))
	if err != nil {
		return "", err
	}

	return tokenString, nil
}
