package services

import (
	"context"
	"database/sql"
	"errors"
	"os"
	"regexp"
	"time"

	"github.com/ziad-eliwa/jit-version-control-system/internal/database"
	"github.com/ziad-eliwa/jit-version-control-system/internal/middleware"
	"github.com/ziad-eliwa/jit-version-control-system/internal/models"
	"github.com/ziad-eliwa/jit-version-control-system/internal/pkg/hashing"
)

var (
	ErrInvalidToken        = errors.New("Invalid Token")
	ErrExpiredToken        = errors.New("Expired Token")
	ErrInvalidCredentials  = errors.New("Invalid Credentials")
	ErrInvalidUsername     = errors.New("Invalid Username")
	ErrInvalidPassword     = errors.New("Invalid Password")
	ErrInvalidEmailAddress = errors.New("Invalid Email Address")
	ErrUserAlreadyExists   = errors.New("User Already Exists")
	ErrEmailAlreadyExists  = errors.New("Email Already Exists")
)

type AuthService struct {
	UserStore  database.UserStore
	TokenStore database.TokenStore
	// JWT Tokens
	JWTSecret      string
	AccessTokenTTL time.Duration
	TokenGenerator middleware.JWTGenerator
	// OAuth
	googleClientID string
	googleSecret   string
	githubClientID string
	githubSecret   string
}

func NewAuthService(userstore database.UserStore) *AuthService {
	return &AuthService{
		UserStore:      userstore,
		googleClientID: os.Getenv("GOOGLE_KEY"),
		googleSecret:   os.Getenv("GOOGLE_SECRET"),
		githubClientID: os.Getenv("GITHUB_KEY"),
		githubSecret:   os.Getenv("GITHUB_SECRET"),
	}
}

func (ah *AuthService) Login(username, password string) (string, error) {
	user, err := ah.UserStore.GetUserbyUsername(username)

	if err == nil || user.Username != "" {
		return "", ErrUserAlreadyExists
	}

	return "", nil
}

func (ah *AuthService) Register(username, password, fullname, email string) (*models.RegisterTokenResponse, error) {
	retrievedUser, err := ah.UserStore.GetUserbyUsername(username)

	if err != sql.ErrNoRows {
		if err != nil {
			return nil, ErrUserAlreadyExists
		}
		return nil, err
	}

	if retrievedUser.EmailAddress == email {
		return nil, ErrEmailAlreadyExists
	}

	// Check Credentials with Regular Expressions -> Done at client side ,too
	emailRegex := regexp.MustCompile(`^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`)
	if !emailRegex.MatchString(email) {
		return nil, ErrInvalidEmailAddress
	}
	usernameRegex := regexp.MustCompile(`^[a-zA-Z0-9]{5,20}$`)
	if !usernameRegex.MatchString(username) {
		return nil, ErrInvalidUsername
	}
	passRegex := regexp.MustCompile(`^(?=.*[a-z])(?=.*[A-Z])(?=.*\d)(?=.*[!@#$%*])[A-Za-z\d!@#$%*]{8,}$`)
	if !passRegex.MatchString(password) {
		return nil, ErrInvalidPassword
	}

	// Create User
	var pass hashing.Password
	pass.Set(password)

	registeredUser := &database.User{
		Username:     username,
		Password:     pass,
		FullName:     fullname,
		EmailAddress: email,
	}

	createdUser, err := ah.UserStore.CreateUser(registeredUser)
	if err != nil {
		return nil, err
	}
	// Generate Tokens - Store Refresh Token
	tokens, err := ah.GenerateAccessTokens(createdUser)

	if err != nil {
		return nil, err
	}

	err = ah.TokenStore.StoreRefreshToken(username, tokens.RefreshToken)

	if err != nil {
		// e := ah.UserStore.DeleteUser(username)
		// if e != nil {
		// 	return nil, e
		// }
		return nil, err
	}
	//Return Tokens
	return &models.RegisterTokenResponse{
		Status:       true,
		AccessToken:  tokens.AccessToken,
		RefreshToken: tokens.RefreshToken,
	}, nil
}

func (ah *AuthService) GenerateAccessTokens(user *database.User) (*models.TokenResponse, error) {
	// Generate JWT Tokens
	tokens, err := ah.TokenGenerator(context.Background(), user.Username)
	if err != nil {
		return nil, err
	}
	return &models.TokenResponse{
		AccessToken:  tokens.AccessToken,
		RefreshToken: tokens.RefreshToken,
	}, nil
}

func ValidateAccessToken(token string) (bool, error) {

	return true, nil
}
