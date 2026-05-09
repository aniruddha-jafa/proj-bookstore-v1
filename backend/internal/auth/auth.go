package auth

import (
	"crypto/rand"
	"encoding/hex"
	"log"
	"net/http"
	"strings"

	"github.com/aniruddha-jafa/go-auth-v1/internal/apperrors"
)

const AUTHORIZATION string = "Authorization"
const BEARER string = "Bearer " // space is required
const REFRESH_TOKEN_BYTES = 32

func GetBearerToken(headers http.Header) (string, error) {
	if headers == nil || headers.Get(AUTHORIZATION) == "" {
		return "", apperrors.ErrMissingAuthorizationHeader
	}
	// Q: Would we ever need multiple AUTHORIZATION values?
	authHeader := strings.Join(headers.Values(AUTHORIZATION), "")
	if !strings.HasPrefix(authHeader, BEARER) {
		return "", apperrors.ErrMissingBearerToken
	}
	headersSlice := strings.Split(authHeader, " ")
	if len(headersSlice) != 2 {
		return "", apperrors.ErrInvalidCredentials
	}
	return headersSlice[1], nil
}

func MakeRefreshToken() (string, error) {
	buf := make([]byte, REFRESH_TOKEN_BYTES)
	_, err := rand.Read(buf)
	if err != nil {
		return "", err
	}
	refreshToken := hex.EncodeToString(buf)
	if len(refreshToken) != REFRESH_TOKEN_BYTES*2 {
		log.Printf("Unexpected refresh token length: %d, expected: %d", len(refreshToken), REFRESH_TOKEN_BYTES*2)
		return "", apperrors.HttpError{Status: http.StatusInternalServerError, Message: "unable to create refresh token"}
	}
	return refreshToken, nil
}
