package middleware

import (
	"context"
	"net/http"
	"os"
	"time"
	jwt "github.com/appleboy/gin-jwt/v3"
	gojwt "github.com/golang-jwt/jwt/v5"
	"github.com/appleboy/gin-jwt/v3/core"
)

func InitAuthMiddleware() *jwt.GinJWTMiddleware {
	authMiddleware := &jwt.GinJWTMiddleware{
		Key:            []byte(os.Getenv("JWT_SECRET")), // From environment
		Timeout:        time.Minute * 15,                // Short-lived access tokens
		MaxRefresh:     time.Hour * 24 * 7,              // 1 week refresh validity
		SecureCookie:   true,                            // HTTPS only
		CookieHTTPOnly: true,                            // Prevent XSS
		CookieSameSite: http.SameSiteStrictMode,         // CSRF protection
		SendCookie:     true,                            // Enable secure cookies
		PayloadFunc: func(data any) gojwt.MapClaims {
			return gojwt.MapClaims{
				"username": data,
			}
		},
	}
	return authMiddleware
}

type JWTGenerator func(context.Context,any)(*core.Token,error)

func InitTokenGenerator(mdw *jwt.GinJWTMiddleware) JWTGenerator {
	return mdw.TokenGenerator
}
