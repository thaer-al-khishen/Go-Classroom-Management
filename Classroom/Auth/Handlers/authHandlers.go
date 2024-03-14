package Handlers

import (
	"encoding/json"
	"errors"
	"github.com/golang-jwt/jwt/v5"
	"github.com/julienschmidt/httprouter"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
	"net/http"
	"strings"
	"time"
	"webapptrials/Classroom/Auth/Models"
	"webapptrials/Classroom/Auth/Utils"
	"webapptrials/Classroom/Secret"
	"webapptrials/Classroom/Shared"
)

var db *gorm.DB

func InitializeDB(d *gorm.DB) {
	db = d
}

const JwtKey = "your_secret_key"

type Claims struct {
	Username string `json:"username"`
	jwt.RegisteredClaims
}

func Login(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	var credentials Models.User
	if err := json.NewDecoder(r.Body).Decode(&credentials); err != nil {
		Shared.SendApiResponse[any](w, http.StatusBadRequest, "Invalid request", nil, "Invalid request body")
		return
	}

	var user Models.User
	// Retrieve the user from the database
	if err := db.Where("username = ?", credentials.Username).First(&user).Error; err != nil {
		Shared.SendApiResponse[any](w, http.StatusUnauthorized, "Login failed", nil, "Invalid username or password")
		return
	}

	// Compare the hashed password
	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(credentials.Password)); err != nil {
		Shared.SendApiResponse[any](w, http.StatusUnauthorized, "Login failed", nil, "Invalid username or password")
		return
	}

	// Generate JWT access token
	accessTokenString, err := Utils.GenerateToken(user.Username)
	if err != nil {
		Shared.SendApiResponse[any](w, http.StatusInternalServerError, "Error", nil, "Failed to generate access token")
		return
	}

	// Generate refresh token (could be another JWT or a random string)
	refreshTokenString, err := Utils.GenerateRefreshToken()
	if err != nil {
		Shared.SendApiResponse[any](w, http.StatusInternalServerError, "Error", nil, "Failed to generate refresh token")
		return
	}

	// Save refresh token in the database with user association
	refreshTokenModel := Models.RefreshTokenModel{
		Token:     refreshTokenString,
		Username:  user.Username,
		ExpiresAt: time.Now().Add(Secret.RefreshTokenExpiry), // Set your desired expiry for refresh tokens
	}
	if result := db.Create(&refreshTokenModel); result.Error != nil {
		Shared.SendApiResponse[any](w, http.StatusInternalServerError, "Error", nil, "Failed to save refresh token")
		return
	}

	// Return both the access token and refresh token to the client
	Shared.SendApiResponse(w, http.StatusOK, "Login successful", map[string]interface{}{
		"accessToken":  accessTokenString,
		"refreshToken": refreshTokenString,
	}, "")
}

func Register(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	var u Models.User
	if err := json.NewDecoder(r.Body).Decode(&u); err != nil {
		Shared.SendApiResponse[any](w, http.StatusBadRequest, "Invalid request", nil, "Invalid request body")
		return
	}

	// Hash the password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(u.Password), bcrypt.DefaultCost)
	if err != nil {
		Shared.SendApiResponse[any](w, http.StatusInternalServerError, "Error", nil, "Failed to hash password")
		return
	}
	u.Password = string(hashedPassword)

	// Create user in DB
	if result := db.Create(&u); result.Error != nil {
		Shared.SendApiResponse[any](w, http.StatusInternalServerError, "Error", nil, "Failed to register user")
		return
	}

	// Mask the password in the response
	u.Password = ""
	Shared.SendApiResponse[any](w, http.StatusCreated, "User registered successfully", u, "")
}

func RefreshToken(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	refreshToken := strings.TrimPrefix(r.Header.Get("Authorization"), "Bearer ")

	username, err := validateRefreshToken(refreshToken)
	if err != nil {
		Shared.SendApiResponse[any](w, http.StatusUnauthorized, "Invalid refresh token", nil, "Invalid or expired refresh token")
		return
	}

	newAccessToken, err := Utils.GenerateToken(username)
	if err != nil {
		Shared.SendApiResponse[any](w, http.StatusInternalServerError, "Error generating token", nil, err.Error())
		return
	}

	Shared.SendApiResponse(w, http.StatusOK, "New access token generated", map[string]string{"accessToken": newAccessToken}, "")
}

func validateRefreshToken(token string) (string, error) {
	refreshTokenRecord, err := FindRefreshToken(token)
	if err != nil {
		return "", errors.New("refresh token not found or expired")
	}

	// The token is valid; return the associated username
	// Optionally, here you might want to issue a new refresh token and invalidate the old one
	return refreshTokenRecord.Username, nil
}

func FindRefreshToken(token string) (*Models.RefreshTokenModel, error) {
	var refreshToken Models.RefreshTokenModel
	result := db.Where("token = ? AND expires_at > ?", token, time.Now()).First(&refreshToken)
	if result.Error != nil {
		return nil, result.Error
	}
	return &refreshToken, nil
}
