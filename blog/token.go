package blog

import (
	"fmt"
	"time"

	"github.com/golang-jwt/jwt"
)

var DntTellAnyone = []byte("IamBatman") //hmacSampleSecret

type JWTClaim struct {
	Email    string `json:"email"`
	Username string `json:"username"`
	jwt.StandardClaims
}

//generating JWT token from user credentials
func GenerateJWT(email, username string) (tokenString string, err error) {
	expirationTime := time.Now().Add(1 * time.Hour)
	claims := &JWTClaim{
		Email:    email,
		Username: username,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: expirationTime.Unix(),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	tokenString, err = token.SignedString(DntTellAnyone)

	return

}

func ValidateToken(signedToken string) {
	token, err := jwt.Parse(signedToken, func(token *jwt.Token) (interface{}, error) {
		return []byte(DntTellAnyone), nil
	})

	if token.Valid {
		fmt.Println("Successful")
	} else if ve, ok := err.(*jwt.ValidationError); ok {
		if ve.Errors&jwt.ValidationErrorMalformed != 0 {
			fmt.Println("That's not even a token")
		} else if ve.Errors&(jwt.ValidationErrorExpired|jwt.ValidationErrorNotValidYet) != 0 {

			fmt.Println("Token is either expired or not active yet")
		} else {
			fmt.Println("Couldn't handle this token:", err)
		}
	} else {
		fmt.Println("Couldn't handle this token:", err)
	}

}

func ExtractClaims(signedToken string) string {
	claims := &JWTClaim{}
	jwt.ParseWithClaims(signedToken, claims, func(token *jwt.Token) (interface{}, error) {
		return []byte(DntTellAnyone), nil
	})

	return claims.Email
}
