package security

import (
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

func MakeJwt(userId uuid.UUID, secret string, duration time.Duration, now time.Time) (string, error) {
	token := jwt.NewWithClaims(
		jwt.SigningMethodHS256,
		jwt.RegisteredClaims{
			Issuer:    "chirpy",
			IssuedAt:  jwt.NewNumericDate(now),
			ExpiresAt: jwt.NewNumericDate(now.Add(duration)),
			Subject:   userId.String(),
		},
	)
	// Signing method HMAC needs bytes
	// https://golang-jwt.github.io/jwt/usage/signing_methods/#signing-methods-and-key-types
	ss, err := token.SignedString([]byte(secret))
	if err != nil {
		return "", err
	}
	return ss, nil
}

// Validates the given token and returns the user id on success
func ValidateJwt(tokenString, secret string) (uuid.UUID, error) {
	token, err := jwt.ParseWithClaims(tokenString, &jwt.RegisteredClaims{}, func(token *jwt.Token) (any, error) {
		// Must match the expected key type of the signing algo
		return []byte(secret), nil
	}, jwt.WithValidMethods([]string{jwt.SigningMethodHS256.Alg()}))
	if err != nil {
		return uuid.UUID{}, err
	}
	uuidStr, err := token.Claims.GetSubject()
	if err != nil {
		return uuid.UUID{}, err
	}
	userId, err := uuid.Parse(uuidStr)
	if err != nil {
		return uuid.UUID{}, err
	}
	return userId, nil
}
