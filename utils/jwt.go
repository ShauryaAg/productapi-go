package utils

import (
	"fmt"
	"os"

	"github.com/ShauryaAg/ProductAPI/models"

	jwt "github.com/dgrijalva/jwt-go"
)

var (
	secret = os.Getenv("SECRET")
)

func CreateToken(user models.User) (string, error) {
	claims := jwt.MapClaims{}
	claims["authorized"] = true
	claims["userId"] = user.Id
	claims["Name"] = user.Name
	claims["email"] = user.Email
	at := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	token, err := at.SignedString([]byte(secret))
	if err != nil {
		return "", err
	}

	return token, nil
}

func ParseToken(tokenString string) (*jwt.MapClaims, error) {

	token, err := jwt.ParseWithClaims(tokenString, &jwt.MapClaims{}, func(token *jwt.Token) (interface{}, error) {
		return []byte(secret), nil
	})
	if err != nil {
		return nil, err
	}

	if claims, ok := token.Claims.(*jwt.MapClaims); ok && token.Valid {
		return claims, err
	}

	fmt.Println("err", err)
	return nil, err
}
