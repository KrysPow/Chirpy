package auth

import (
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

func HashPassword(password string) (string, error) {
	hashed_pwd, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		fmt.Println("Problem occured during hashing password: ", err)
		return "", err
	}

	return string(hashed_pwd), nil
}

func CheckPasswordHash(password, hash string) error {
	return bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
}

func MakeJWT(userID uuid.UUID, tokenSecret string, expiresIn time.Duration) (string, error) {
	string_uuid := userID.String()

	token := jwt.NewWithClaims(jwt.SigningMethodHS256,
		jwt.RegisteredClaims{Issuer: "chirpy",
			IssuedAt:  jwt.NewNumericDate(time.Now().UTC()),
			ExpiresAt: jwt.NewNumericDate(time.Now().UTC().Add(expiresIn)),
			Subject:   string_uuid})

	signed_string, err := token.SignedString([]byte(tokenSecret))
	if err != nil {
		return "", err
	}
	return signed_string, nil
}

func ValidateJWT(tokenString, tokenSecret string) (uuid.UUID, error) {
	var claims jwt.RegisteredClaims
	token, err := jwt.ParseWithClaims(tokenString, &claims,
		func(token *jwt.Token) (any, error) {
			if token.Method.Alg() != "HS256" {
				return nil, fmt.Errorf("unexpected signing method: %v", token.Method.Alg())
			}
			return []byte(tokenSecret), nil
		})
	if err != nil {
		return uuid.UUID{}, err
	}

	subject, err := token.Claims.GetSubject()
	if err != nil {
		return uuid.UUID{}, err
	}

	id, err := uuid.Parse(subject)
	if err != nil {
		return uuid.UUID{}, err
	}

	return id, nil
}

func GetBearerToken(header http.Header) (string, error) {
	auth := header.Get("Authorization")
	if auth == "" {
		return "", fmt.Errorf("no bearer")
	}
	return strings.Fields(auth)[1], nil
}

func MakeFreshToken() (string, error) {
	rand_data := make([]byte, 32)
	_, err := rand.Read(rand_data)
	if err != nil {
		return "", err
	}

	return hex.EncodeToString(rand_data), nil
}

func GetAPIKey(headers http.Header) (string, error) {
	string_apiKey := headers.Get("Authorization")
	if len(strings.Fields(string_apiKey)) != 2 {
		return "", fmt.Errorf("API-key in wrong format")
	}
	return strings.Fields(string_apiKey)[1], nil
}
