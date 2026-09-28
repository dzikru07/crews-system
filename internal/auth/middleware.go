package auth

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

const (
	authContextKey = "crews_auth"
)

func Middleware(
	tokenManager *TokenManager,
) gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")

		if authHeader == "" {
			c.AbortWithStatusJSON(
				http.StatusUnauthorized,
				gin.H{
					"error": "missing_authorization",
				},
			)
			return
		}

		parts := strings.Fields(authHeader)

		if len(parts) != 2 ||
			!strings.EqualFold(parts[0], "Bearer") {
			c.AbortWithStatusJSON(
				http.StatusUnauthorized,
				gin.H{
					"error": "invalid_authorization",
				},
			)
			return
		}

		authContext, err := tokenManager.Validate(
			parts[1],
		)
		if err != nil {
			switch err {
			case ErrExpiredToken:
				c.AbortWithStatusJSON(
					http.StatusUnauthorized,
					gin.H{
						"error": "token_expired",
					},
				)

			default:
				c.AbortWithStatusJSON(
					http.StatusUnauthorized,
					gin.H{
						"error": "invalid_token",
					},
				)
			}

			return
		}

		c.Set(authContextKey, authContext)

		c.Next()
	}
}

func GetAuthContext(
	c *gin.Context,
) (*AuthContext, bool) {
	value, exists := c.Get(authContextKey)

	if !exists {
		return nil, false
	}

	authContext, ok := value.(*AuthContext)

	return authContext, ok
}

func GetTenantID(c *gin.Context) string {
	authContext, ok := GetAuthContext(c)

	if !ok {
		return ""
	}

	return authContext.TenantID.String()
}

func GetIdentityID(c *gin.Context) string {
	authContext, ok := GetAuthContext(c)

	if !ok {
		return ""
	}

	return authContext.IdentityID.String()
}

func GetClientID(c *gin.Context) string {
	authContext, ok := GetAuthContext(c)

	if !ok {
		return ""
	}

	return authContext.ClientID
}
