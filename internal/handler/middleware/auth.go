package middleware

import (
	"errors"
	"net/http"
	"strings"

	"github.com/PosokhovVadim/stawberry/internal/app/apperror"
	"github.com/PosokhovVadim/stawberry/internal/domain/entity"
	"github.com/gin-gonic/gin"
	"golang.org/x/net/context"
)

type SessionManager interface {
	VerifyAccessToken(ctx context.Context, token string) (claims entity.Claims, err error)
}

func AuthMiddleware(sessionManager SessionManager) gin.HandlerFunc {
	return func(c *gin.Context) {
		tokenString := c.GetHeader("Authorization")
		if tokenString == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Authorization header required"})
			c.Abort()
			return
		}

		parts := strings.Split(tokenString, " ")
		if len(parts) != 2 || strings.ToLower(parts[0]) != "bearer" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid authorization header format"})
			c.Abort()
			return
		}

		claims, err := sessionManager.VerifyAccessToken(c.Request.Context(), parts[1])
		if err != nil {
			switch {
			case errors.Is(err, apperror.ErrSessionExpired):
				c.JSON(http.StatusUnauthorized, gin.H{"code": apperror.TokenExpired, "message": "Token expired"})
			default:
				c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token"})
			}
			c.Abort()
			return
		}

		c.Set("userID", claims.UserID)
		c.Next()
	}
}
