package Handlers

import (
	"errors"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
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

func Login(c *gin.Context) {
	var credentials Models.User
	if err := c.ShouldBindJSON(&credentials); err != nil {
		Shared.SendGinGenericApiResponse[any](c, http.StatusBadRequest, "Invalid request", nil, "Invalid request body")
		return
	}

	var user Models.User
	// Retrieve the user from the database
	if err := db.Where("username = ?", credentials.Username).First(&user).Error; err != nil {
		Shared.SendGinGenericApiResponse[any](c, http.StatusUnauthorized, "Login failed", nil, "Invalid username or password")
		return
	}

	// Compare the hashed password
	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(credentials.Password)); err != nil {
		Shared.SendGinGenericApiResponse[any](c, http.StatusUnauthorized, "Login failed", nil, "Invalid username or password")
		return
	}

	// Generate JWT access token
	accessTokenString, err := Utils.GenerateToken(user)
	if err != nil {
		Shared.SendGinGenericApiResponse[any](c, http.StatusInternalServerError, "Error", nil, "Failed to generate access token")
		return
	}

	// Generate refresh token (could be another JWT or a random string)
	refreshTokenString, err := Utils.GenerateRefreshToken()
	if err != nil {
		Shared.SendGinGenericApiResponse[any](c, http.StatusInternalServerError, "Error", nil, "Failed to generate refresh token")
		return
	}

	// Save refresh token in the database with user association
	refreshTokenModel := Models.RefreshTokenModel{
		Token:     refreshTokenString,
		Username:  user.Username,
		ExpiresAt: time.Now().Add(Secret.RefreshTokenExpiry), // Set your desired expiry for refresh tokens
	}
	if result := db.Create(&refreshTokenModel); result.Error != nil {
		Shared.SendGinGenericApiResponse[any](c, http.StatusInternalServerError, "Error", nil, "Failed to save refresh token")
		return
	}

	// Return both the access token and refresh token to the client
	Shared.SendGinGenericApiResponse(c, http.StatusOK, "Login successful", map[string]interface{}{
		"accessToken":  accessTokenString,
		"refreshToken": refreshTokenString,
	}, "")
}

func Register(c *gin.Context) {
	var u Models.User
	if err := c.ShouldBindJSON(&u); err != nil {
		Shared.SendGinGenericApiResponse[any](c, http.StatusBadRequest, "Invalid request", nil, "Invalid request body")
		return
	}

	// Check if the role is not provided and default to Student
	if u.Role == nil {
		defaultRole := Models.Student
		u.Role = &defaultRole
	} else if *u.Role == Models.Admin {
		// Ensure an admin cannot be created through this endpoint
		Shared.SendGinGenericApiResponse[any](c, http.StatusBadRequest, "Invalid request", nil, "You can't create an admin")
		return
	}

	// Hash the password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(u.Password), bcrypt.DefaultCost)
	if err != nil {
		Shared.SendGinGenericApiResponse[any](c, http.StatusInternalServerError, "Error", nil, "Failed to hash password")
		return
	}
	u.Password = string(hashedPassword)

	//Create a student by default
	//u.Role = Models.Student

	// Create user in DB
	if result := db.Create(&u); result.Error != nil {
		Shared.SendGinGenericApiResponse[any](c, http.StatusInternalServerError, "Error", nil, "Failed to register user")
		return
	}

	// Mask the password in the response
	u.Password = ""
	Shared.SendGinGenericApiResponse[any](c, http.StatusCreated, "User registered successfully", u, "")
}

func RefreshToken(c *gin.Context) {
	refreshToken := strings.TrimPrefix(c.GetHeader("Authorization"), "Bearer ")

	username, err := validateRefreshToken(refreshToken)
	if err != nil {
		Shared.SendGinGenericApiResponse[any](c, http.StatusUnauthorized, "Invalid refresh token", nil, "Invalid or expired refresh token")
		return
	}

	user, err := GetUserByUsername(username)
	if err != nil {
		Shared.SendGinGenericApiResponse[any](c, http.StatusInternalServerError, "User not found", nil, err.Error())
		return
	}

	newAccessToken, err := Utils.GenerateToken(user)
	if err != nil {
		Shared.SendGinGenericApiResponse[any](c, http.StatusInternalServerError, "Error generating token", nil, err.Error())
		return
	}

	Shared.SendGinGenericApiResponse(c, http.StatusOK, "New access token generated", map[string]string{"accessToken": newAccessToken}, "")
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
	if result := db.Where("token = ? AND expires_at > ?", token, time.Now()).First(&refreshToken); result.Error != nil {
		return nil, result.Error
	}
	return &refreshToken, nil
}

func GetUserByUsername(username string) (Models.User, error) {
	var user Models.User
	result := db.Where("username = ?", username).First(&user)
	if result.Error != nil {
		return Models.User{}, result.Error
	}
	return user, nil
}
