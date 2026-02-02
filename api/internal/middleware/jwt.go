package middleware

import (
	"github.com/gin-gonic/gin"
	"github.com/huya_live/api/pkg/jwt"
	"github.com/huya_live/api/pkg/response"
)

func JWTRequired(jwtManager *jwt.Manager) gin.HandlerFunc {
	return func(c *gin.Context) {
		token := c.GetHeader("Authorization")
		if token == "" || len(token) < 8 || token[:7] != "Bearer " {
			response.Unauthorized(c, "missing or invalid authorization header")
			c.Abort()
			return
		}

		claims, err := jwtManager.ValidateToken(token[7:])
		if err != nil {
			if err == jwt.ErrExpiredToken {
				response.Unauthorized(c, "token has expired")
			} else {
				response.Unauthorized(c, "invalid token")
			}
			c.Abort()
			return
		}

		c.Set("user_id", claims.UserID)
		c.Set("user_role", claims.Role)
		c.Set("user_level", claims.Level)
		c.Next()
	}
}
