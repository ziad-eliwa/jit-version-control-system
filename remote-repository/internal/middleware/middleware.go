package middleware

import (
	"crypto/rand"
	"encoding/base64"
	"errors"
	"log/slog"

	"net/http"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"github.com/ziad-eliwa/jit-version-control-system/internal/database"
	"github.com/ziad-eliwa/jit-version-control-system/internal/models"
)

var (
	ErrInvalidToken         = errors.New("Error Invalid Token")
	ErrExpiredToken         = errors.New("Error Expired Token")
	ErrUsernameNotInContext = errors.New("Error username not in context")
	ErrInvalidUsernameType  = errors.New("Error invalid username type")
)

type AuthenticationMiddleware struct {
	TokenStore database.TokenStore
	RepoStore  database.RepoStore

	JWTSecret   string
	Timeout     time.Duration
	MaxRefresh  time.Duration
	Logger      *slog.Logger
	IdentityKey string
}

// Generation
func (am *AuthenticationMiddleware) GenerateJWTToken(data any) (string, error) {
	expirationTime := time.Now().Add(am.Timeout)

	claims := jwt.MapClaims{
		am.IdentityKey: "ziadeliwa",
		"exp":          expirationTime.Unix(),
		"iat":          time.Now().Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	tokenString, err := token.SignedString(am.JWTSecret)

	if err != nil {
		return "", err
	}

	return tokenString, nil
}

func (am *AuthenticationMiddleware) GenerateAccessToken(username string) (string, error) {
	bytes := make([]byte, 32)
	if _, err := rand.Read(bytes); err != nil {
		return "", err
	}
	return base64.URLEncoding.EncodeToString(bytes), nil
}

// Token Validation
func (am *AuthenticationMiddleware) ValidateJWTToken(c *gin.Context) (jwt.MapClaims, error) {
	var token string
	var err error
	authHeader := c.GetHeader("Authorization")
	if authHeader == "" {
		c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Missing Authorization Header"})
	}
	const Prefix = "Bearer "

	if !strings.HasPrefix(authHeader, Prefix) {
		c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Missing Bearer Prefix"})
	}

	token = strings.TrimPrefix(authHeader, Prefix)

	jwtToken, err := jwt.Parse(token, func(t *jwt.Token) (any, error) {
		if jwt.SigningMethodHS256 != t.Method {
			return nil, jwt.ErrSignatureInvalid
		}

		c.Set("JWT_TOKEN", token)
		return am.JWTSecret, nil
	})

	if err != nil {
		if err == jwt.ErrTokenExpired {
			return nil, ErrExpiredToken
		}
		return nil, err
	}

	if claims, ok := jwtToken.Claims.(jwt.MapClaims); ok && jwtToken.Valid {
		return claims, nil
	}
	return nil, ErrInvalidToken
}

func (am *AuthenticationMiddleware) ValidateRefreshToken(token string) (bool, error) {
	refreshToken, err := am.TokenStore.GetRefreshToken(token)

	if err != nil {
		return false, err
	}

	if refreshToken.Revoked {
		return false, nil
	}

	return true, nil
}

func (am *AuthenticationMiddleware) RevokeRefreshToken(username string) error {
	err := am.TokenStore.RevokeAllTokens(username)

	if err != nil {
		return err
	}

	return nil
}

// Middleware
func (am *AuthenticationMiddleware) Autheticate() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		claims, err := am.ValidateJWTToken(ctx)

		if err != nil {
			am.Logger.Error("Error Authenticating")
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": err})
			return
		}

		ctx.Set(am.IdentityKey, claims[am.IdentityKey])
		ctx.Next()
	}
}

func (am *AuthenticationMiddleware) ExtractUserFromContext(c *gin.Context) (string, error) {
	user, ok := c.Get(am.IdentityKey)
	if !ok {
		return "", ErrUsernameNotInContext
	}
	username, ok := user.(string)

	if !ok {
		return "", ErrInvalidUsernameType
	}

	return username, nil
}

func ExtractRefreshTokenFromContext(c *gin.Context) (*models.LogoutRequest, error) {
	Refresh := &models.LogoutRequest{}
	err := c.BindJSON(&Refresh)
	if err != nil {
		return nil, err
	}
	if Refresh.RefreshToken == "" {
		return nil, errors.New("Refresh Token does not exist in Request")
	}
	return Refresh, nil
}

func (am *AuthenticationMiddleware) AuthorizePrivacy() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		user := ctx.Param("username")

		if user == "" {
			ctx.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "no username was specified in url"})
		}

		repo := ctx.Param("reponame")

		if repo == "" {
			ctx.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "no repository name was specified in url"})
		}

		privacy, err := am.RepoStore.GetRepoPrivacy(user, repo)

		if err != nil {
			ctx.AbortWithStatusJSON(http.StatusNotFound, gin.H{"error": "repository not found"})
		}

		ctx.Set("PRIVACY", privacy)
		ctx.Next()
	}
}

func (am *AuthenticationMiddleware) AuthorizeAccess() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		currentUser, err := am.ExtractUserFromContext(ctx)

		if err != nil {
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": ErrUsernameNotInContext})
		}

		user := ctx.Param("username")

		if user == "" {
			ctx.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "no username was specified in url"})
		}

		repo := ctx.Param("reponame")

		if repo == "" {
			ctx.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "no repository name was specified in url"})
		}

		authorized, err := am.RepoStore.GetAccessStatusOnRepo(user, repo, currentUser)

		if err != nil {
			ctx.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": "internal server error"})
		}

		if !authorized {
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "no access is granted to this repository"})
		}

		ctx.Next()
	}
}

func (am *AuthenticationMiddleware) AuthorizeOwnership() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		currentUser, err := am.ExtractUserFromContext(ctx)

		if err != nil {
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": ErrUsernameNotInContext})
		}

		user := ctx.Param("username")

		if user == "" {
			ctx.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "no username was specified in url"})
		}

		repo := ctx.Param("reponame")

		if repo == "" {
			ctx.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "no repository name was specified in url"})
		}

		if user == currentUser {
			ctx.Set("OWNER", true)
		} else {
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error":"you are not the owner of this repositories"})
		}

		ctx.Next()
	}
}