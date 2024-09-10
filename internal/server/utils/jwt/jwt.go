package jwt

import (
	"errors"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/dgrijalva/jwt-go"
	_ "golang.org/x/crypto/bcrypt"
)

const (
	_secretKey    = "HNG4wHwDkO5DvSqQK1vb8EetGPrfAcBuR3UwU6Nejms"
	JwtCookieName = "jwt_auth"
)

var (
	ErrIncorrectCookieName = errors.New("ParseJWT: incorrect cookie name")
)

func GenerateJWT(id int) (string, error) {
	claims := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.StandardClaims{
		Subject:   strconv.Itoa(id),
		ExpiresAt: time.Now().Add(time.Hour * 24).Unix(),
	})

	token, err := claims.SignedString([]byte(_secretKey))
	if err != nil {
		return "", err
	}

	return token, nil
}

func ParseJWT(jwtCookie *http.Cookie) (*jwt.StandardClaims, error) {
	token, err := jwt.ParseWithClaims(jwtCookie.Value, &jwt.StandardClaims{}, func(token *jwt.Token) (interface{}, error) {
		return []byte(_secretKey), nil
	})
	if err != nil || !token.Valid {
		log.Println(err)
		return nil, err
	}

	claims := token.Claims.(*jwt.StandardClaims)

	return claims, nil
}
