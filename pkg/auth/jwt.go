package auth

import (
	"fmt"
	jwt "github.com/dgrijalva/jwt-go"
	"log"
	"redditclone/pkg/models"
	"time"
)

var jwtKey = []byte("secret")

type Claims struct {
	UserID   string `json:"id"`
	Username string `json:"username"`
	jwt.StandardClaims
}

func GenerateToken(userID, username string) (string, error) {
	exp := time.Now().Add(time.Hour * 72).Unix()
	claims := &Claims{
		UserID:   userID,
		Username: username,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: exp,
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(jwtKey)
}

func ParseToken(inToken string) (*models.Session, error) {
	hashSecretGetter := func(token *jwt.Token) (interface{}, error) {
		method, ok := token.Method.(*jwt.SigningMethodHMAC)
		if !ok || method.Alg() != "HS256" {
			return nil, fmt.Errorf("bad sign method")
		}
		return jwtKey, nil
	}
	token, err := jwt.Parse(inToken, hashSecretGetter)
	if err != nil {
		return nil, fmt.Errorf("invalid parse token")
	}
	fmt.Printf("\t\tpayload: %+v\n", token)
	log.Printf("\t\tpayload: %+v\n", token)

	payload, ok := token.Claims.(jwt.MapClaims)
	if !ok || !token.Valid {
		return nil, fmt.Errorf("invalid claims token")
	}
	fmt.Printf("\t\tpayload: %+v\n", payload)
	log.Printf("\t\tpayload: %+v\n", payload)
	session := &models.Session{
		ID:       payload["id"].(string),
		Username: payload["username"].(string),
	}
	return session, nil
}
