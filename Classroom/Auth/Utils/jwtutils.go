package Utils

import (
	"crypto/rand"
	"encoding/hex"
	"errors"
	"github.com/golang-jwt/jwt/v5"
	"strings"
	"time"
	"webapptrials/Classroom/Auth/Models"
	"webapptrials/Classroom/Secret"
)

// Modify your Claims struct to include the role
type Claims struct {
	Username string          `json:"username"`
	Role     Models.UserRole `json:"role"`
	jwt.RegisteredClaims
}

func GenerateToken(user Models.User) (string, error) {
	expirationTime := time.Now().Add(Secret.AccessTokenExpiry)
	claims := &Claims{
		Username: user.Username,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(expirationTime),
		},
		Role: *user.Role,
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString(Secret.JwtKey)

	return tokenString, err
}

func GenerateRefreshToken() (string, error) {
	// 32 bytes gives a 256-bit string, considered very secure
	bytes := make([]byte, 32)
	if _, err := rand.Read(bytes); err != nil {
		return "", err
	}
	return hex.EncodeToString(bytes), nil
}

func ValidateToken(tokenString string) (*Claims, error) {
	claims := &Claims{}

	token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
		return Secret.JwtKey, nil
	})

	if err != nil {
		return nil, err
	}

	if !token.Valid {
		return nil, errors.New("invalid token, ValidationErrorMalformed")
	}

	return claims, nil
}

// ParseToken parses a JWT token string, validates it, and returns the custom claims.
func ParseToken(tokenStr string) (*Claims, error) {
	// Remove the Bearer prefix, if present
	tokenStr = strings.TrimPrefix(tokenStr, "Bearer ")

	// Parse the token
	token, err := jwt.ParseWithClaims(tokenStr, &Claims{}, func(token *jwt.Token) (interface{}, error) {
		// Ensure the token's algorithm matches "HS256"
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("unexpected signing method")
		}

		return Secret.JwtKey, nil
	})

	if err != nil {
		return nil, err
	}

	// Validate the token and ensure the token claims match the expected type
	claims, ok := token.Claims.(*Claims)
	if !ok || !token.Valid {
		return nil, errors.New("invalid token")
	}

	return claims, nil
}
