package auth

import (
	"flag"
	"fmt"
	"os"
	"redditclone/pkg/models"
	"time"

	jwt "github.com/golang-jwt/jwt/v5"
)

var jwtKey []byte

func Init() {
	secret := os.Getenv("JWT_SECRET")

	flag.StringVar(&secret, "jwtSecret", secret, "JWT Secret")
	flag.Parse()

	if secret == "" {
		panic("JWT Secret is empty")
	}
	jwtKey = []byte(secret)

}

type Claims struct {
	UserID   string `json:"id"`
	Username string `json:"username"`
	jwt.RegisteredClaims
}

func GenerateToken(userID, username string) (string, error) {
	exp := time.Now().Add(time.Hour * 72)
	claims := &Claims{
		UserID:   userID,
		Username: username,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(exp),
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

	payload, ok := token.Claims.(jwt.MapClaims)
	if !ok || !token.Valid {
		return nil, fmt.Errorf("invalid claims token")
	}

	session := &models.Session{
		ID:       payload["id"].(string),
		Username: payload["username"].(string),
	}
	return session, nil
}
