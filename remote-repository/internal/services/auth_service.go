package services

import (
	"database/sql"
	"errors"

	// "os"
	"regexp"

	"github.com/ziad-eliwa/jit-version-control-system/internal/database"
	"github.com/ziad-eliwa/jit-version-control-system/internal/middleware"
	"github.com/ziad-eliwa/jit-version-control-system/internal/models"
	"github.com/ziad-eliwa/jit-version-control-system/internal/pkg/hashing"
	"github.com/ziad-eliwa/jit-version-control-system/internal/utils"
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
	ErrUserNotFound        = errors.New("User not found")
	ErrIncorrectPassword   = errors.New("Password is not correct")
)

type AuthService struct {
	UserStore  database.UserStore
	TokenStore database.TokenStore
	// Middleware
	Authentication *middleware.AuthenticationMiddleware
	// OAuth
	googleClientID string
	googleSecret   string
	githubClientID string
	githubSecret   string
}

func NewAuthService(userstore database.UserStore, tokenstore database.TokenStore, authentication *middleware.AuthenticationMiddleware) *AuthService {
	return &AuthService{
		Authentication: authentication,
		UserStore:  userstore,
		TokenStore: tokenstore,
		// googleClientID: os.Getenv("GOOGLE_KEY"),
		// googleSecret:   os.Getenv("GOOGLE_SECRET"),
		// githubClientID: os.Getenv("GITHUB_KEY"),
		// githubSecret:   os.Getenv("GITHUB_SECRET"),
	}
}

func (ah *AuthService) Login(username, password string) (*models.TokenResponse, error) {
	user, err := ah.UserStore.GetUserbyUsername(username)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, ErrUserNotFound
		}
		return nil, err
	}

	if user.Username == "" || string(user.PasswordHash) == "" || user.EmailAddress == "" {
		return nil, ErrInvalidCredentials
	}

	var pass hashing.Password
	pass.Hash = []byte(user.PasswordHash)

	if ok,_ := pass.MatchPassword([]byte(password)); !ok {
		return nil, ErrIncorrectPassword
	}

	tokens, err := ah.GenerateAccessTokens(user.Username)

	if err != nil {
		return nil, err
	}

	err = ah.TokenStore.StoreRefreshToken(username, tokens.RefreshToken)

	if err != nil {
		return nil, err
	}
	//Return Tokens
	return &models.TokenResponse{
		RefreshToken: tokens.RefreshToken,
		AccessToken:  tokens.AccessToken,
	}, nil
}

func (ah *AuthService) Register(username, password, fullname, email string) (*models.TokenResponse, error) {
	_, err := ah.UserStore.GetUserbyUsername(username)

	if err != sql.ErrNoRows {
		if err != nil {
			return nil, err
		}
		return nil, ErrUserAlreadyExists
	}

	_, err = ah.UserStore.GetUserbyEmailAddress(email)

	if err != sql.ErrNoRows {
		if err != nil {
			return nil, err
		}
		return nil, ErrEmailAlreadyExists
	}

	// Check Credentials with Regular Expressions -> Done at client side ,too
	emailRegex := regexp.MustCompile(`^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`)
	if !emailRegex.MatchString(email) {
		return nil, ErrInvalidEmailAddress
	}

	usernameRegex := regexp.MustCompile(`^[a-z0-9_-]{5,20}$`)
	if !usernameRegex.MatchString(username) {
		return nil, ErrInvalidUsername
	}

	if !utils.IsValidPassword(password) {
		return nil, ErrInvalidPassword
	}

	// Create User
	var pass hashing.Password
	pass.Set(password)

	registeredUser := &database.User{
		Username:     username,
		PasswordHash: string(pass.Hash),
		FullName:     fullname,
		EmailAddress: email,
	}

	// Generate Tokens - Store Refresh Token
	tokens, err := ah.GenerateAccessTokens(registeredUser.Username)

	if err != nil {
		return nil, err
	}

	_, err = ah.UserStore.CreateUser(registeredUser)

	if err != nil {
		return nil, err
	}

	err = ah.TokenStore.StoreRefreshToken(username, tokens.RefreshToken)

	if err != nil {
		return nil, err
	}
	//Return Tokens
	return &models.TokenResponse{
		Status:       true,
		AccessToken:  tokens.AccessToken,
		RefreshToken: tokens.RefreshToken,
	}, nil
}

func (ah *AuthService) GenerateAccessTokens(username string) (*models.TokenResponse, error) {
	// Generate Access Tokens JWT
	accessToken, err := ah.Authentication.GenerateJWTToken(username)
	if err != nil {
		return nil, err
	}
	// Generate Refresh Tokens - Opaque
	refreshToken, err := ah.Authentication.GenerateAccessToken(username)
	if err != nil {
		return nil, err
	}

	return &models.TokenResponse{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
	}, nil
}
