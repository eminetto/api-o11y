package security

import (
	"fmt"
	"time"

	jwt "github.com/dgrijalva/jwt-go"
)

const (
	JWT_SECRET  = "d7830ad5791dsdsds"
	JWT_EXP_HOUR=1
	JWT_EXP_MIN = 0
	JWT_EXP_SEC = 30
)

//NewToken create a new token
func NewToken(email string) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"email":   email,
		"nbf":       time.Now().Unix(),
		"iat":       time.Now().Unix(),
		"exp":       time.Now().Local().Add(time.Hour*time.Duration(JWT_EXP_HOUR) + time.Minute*time.Duration(JWT_EXP_MIN) + time.Second*time.Duration(JWT_EXP_SEC)).Unix(),
	})
	// Sign and get the complete encoded token as a string using the secret
	sToken, err := token.SignedString([]byte(JWT_SECRET))
	if err != nil {
		return "", err
	}

	return sToken, nil
}

//ParseToken parse a token
func ParseToken(tokenString string) (*jwt.Token, error) {
	var token *jwt.Token
	var err error
	token, err = parseHS256(tokenString, token)
	if err != nil && err.Error() != "Token is expired" {
		token, err = parseHS256(tokenString, token)
	}

	return token, err
}

func parseHS256(tokenString string, token *jwt.Token) (*jwt.Token, error) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("Unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(JWT_SECRET), nil
	})
	return token, err
}

//GetClaims get claims information
func GetClaims(token *jwt.Token) (jwt.MapClaims, error) {
	if !token.Valid {
		return nil, fmt.Errorf("Unauthorized")
	}
	err := token.Claims.(jwt.MapClaims).Valid()
	if err != nil {
		return nil, err
	}

	return token.Claims.(jwt.MapClaims), nil
}