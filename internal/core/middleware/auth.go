package middleware

import (
	"context"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	coreHttp "github.com/hogiabao7725/gin-auth-playground/internal/core/delivery/http"
	"github.com/hogiabao7725/gin-auth-playground/internal/infra/token"
)

type contextKey string

const (
	UserIDKey contextKey = "userID"
	RoleKey   contextKey = "role"
)

type AuthMiddleware struct {
	jwt *token.JWT
}

func NewAuthMiddleware(jwt *token.JWT) *AuthMiddleware {
	return &AuthMiddleware{jwt: jwt}
}

func (m *AuthMiddleware) RequireAuth() gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			coreHttp.AbortWithError(c, http.StatusUnauthorized, "authorization header is required")
			return
		}

		parts := strings.SplitN(authHeader, " ", 2)
		if len(parts) != 2 || strings.ToLower(parts[0]) != "bearer" {
			coreHttp.AbortWithError(c, http.StatusUnauthorized, "authorization header format must be Bearer {token}")
			return
		}

		accessToken := parts[1]
		claims, err := m.jwt.ParseAccessToken(accessToken)
		if err != nil {
			coreHttp.AbortWithError(c, http.StatusUnauthorized, "invalid or expired token")
			return
		}

		// 1. Store in Gin's internal map (c.Keys)
		// This is for outer middlewares or HTTP handlers that need quick access via c.GetString().
		// Note: This data will not propagate down to the Core layers (UseCase/Domain).
		c.Set(string(UserIDKey), claims.Subject)
		c.Set(string(RoleKey), claims.Role)

		// 2. Wrap and store in the Native Go Context
		// Used in handler (usecase)
		ctx := c.Request.Context()
		ctx = context.WithValue(ctx, UserIDKey, claims.Subject)
		ctx = context.WithValue(ctx, RoleKey, claims.Role)
		c.Request = c.Request.WithContext(ctx)

		c.Next()
	}
}

func GetUserID(ctx context.Context) string {
	if val, ok := ctx.Value(UserIDKey).(string); ok {
		return val
	}
	return ""
}

func GetRole(ctx context.Context) string {
	if val, ok := ctx.Value(RoleKey).(string); ok {
		return val
	}
	return ""
}

// RequireRole checks if the user has one of the allowed roles.
// Must be used AFTER RequireAuth middleware.
func RequireRole(allowedRoles ...string) gin.HandlerFunc {
	return func(c *gin.Context) {
		role := GetRole(c.Request.Context())

		for _, r := range allowedRoles {
			if role == r {
				c.Next()
				return
			}
		}

		coreHttp.AbortWithError(c, http.StatusForbidden, "you do not have permission to access this resource")
	}
}
