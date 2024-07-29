package auth

import (
	"fmt"
	jwt "github.com/golang-jwt/jwt/v5"
	"net/http"
	"os"
	"time"
)

type Claims struct {
	//in the future role
	jwt.RegisteredClaims
}

func DefaultCookie(c *http.Cookie) {
	c.Value = ""
	c.Expires = time.Now().Add(time.Hour * (-1))

}
func CreateJWTTokenCookieUser(w http.ResponseWriter, id string) error {
	claims := &Claims{
		RegisteredClaims: jwt.RegisteredClaims{
			Issuer:    "localhost:8080",
			Subject:   id,
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Hour * 24)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}
	/*	err := godotenv.Load()
		if err != nil {
			log.Fatal("Error loading .env file")
		}*/

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString([]byte(os.Getenv("JWT_SECRET")))
	if err != nil {
		fmt.Println(err)
		return err
	}
	http.SetCookie(w, &http.Cookie{
		HttpOnly: true,
		Name:     "jwt-token",
		Value:    tokenString,
		Expires:  time.Now().Add(time.Hour * 24)},
	)
	return nil
}

// isLogedIn func
