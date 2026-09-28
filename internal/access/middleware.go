package access

import (
	"net/http"
	"strings"

	"github.com/dzikru07/crews-system/internal/api_registry"
	"github.com/dzikru07/crews-system/internal/auth"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

func RequirePermission(
	service *Service,
	endpointID uuid.UUID,
) gin.HandlerFunc {

	return func(c *gin.Context) {

		authContext, ok := auth.GetAuthContext(c)

		if !ok {
			c.AbortWithStatusJSON(
				http.StatusUnauthorized,
				gin.H{
					"error": "authentication_context_missing",
				},
			)
			return
		}

		allowed, err := service.HasPermission(
			c.Request.Context(),
			authContext.IdentityID,
			endpointID,
		)

		if err != nil {
			c.AbortWithStatusJSON(
				http.StatusInternalServerError,
				gin.H{
					"error": "permission_check_failed",
				},
			)
			return
		}

		if !allowed {
			c.AbortWithStatusJSON(
				http.StatusForbidden,
				gin.H{
					"error": "access_denied",
				},
			)
			return
		}

		c.Next()
	}
}

func Middleware(
	service *Service,
) gin.HandlerFunc {

	return func(c *gin.Context) {

		// Authentication harus sudah dijalankan sebelumnya.
		authContext, ok := auth.GetAuthContext(c)

		if !ok {
			c.AbortWithStatusJSON(
				http.StatusUnauthorized,
				gin.H{
					"error": "authentication_context_missing",
				},
			)
			return
		}

		method := strings.ToUpper(
			c.Request.Method,
		)

		path := c.Request.URL.Path

		endpoint, err := service.ResolveEndpoint(
			c.Request.Context(),
			method,
			path,
		)

		if err != nil {
			c.AbortWithStatusJSON(
				http.StatusNotFound,
				gin.H{
					"error": "api_endpoint_not_registered",
				},
			)
			return
		}

		// Endpoint tidak membutuhkan authentication.
		if !endpoint.RequiresAuth {
			c.Next()
			return
		}

		allowed, err := service.HasPermission(
			c.Request.Context(),
			authContext.IdentityID,
			endpoint.ID,
		)

		if err != nil {
			c.AbortWithStatusJSON(
				http.StatusInternalServerError,
				gin.H{
					"error": "permission_check_failed",
				},
			)
			return
		}

		if !allowed {
			c.AbortWithStatusJSON(
				http.StatusForbidden,
				gin.H{
					"error": "access_denied",
				},
			)
			return
		}

		// Simpan endpoint ke context supaya handler/service
		// berikutnya tidak perlu resolve ulang.
		c.Set("crews_endpoint", endpoint)

		c.Next()
	}
}

func GetEndpoint(
	c *gin.Context,
) (*api_registry.APIEndpoint, bool) {

	value, exists := c.Get("crews_endpoint")

	if !exists {
		return nil, false
	}

	endpoint, ok := value.(*api_registry.APIEndpoint)

	return endpoint, ok
}

func Authorization(
	accessService *Service,
) gin.HandlerFunc {

	return func(c *gin.Context) {

		// =====================================================
		// AUTH CONTEXT
		// =====================================================

		authContext, ok := auth.GetAuthContext(c)

		if !ok {
			c.AbortWithStatusJSON(
				http.StatusUnauthorized,
				gin.H{
					"error": "authentication_context_missing",
				},
			)
			return
		}

		// =====================================================
		// FIND REGISTERED ENDPOINT
		// =====================================================

		endpoint, err := accessService.GetEndpointByRoute(
			c.Request.Context(),
			c.Request.Method,
			c.Request.URL.Path,
		)

		if err != nil {

			if err == ErrEndpointNotFound {
				c.AbortWithStatusJSON(
					http.StatusForbidden,
					gin.H{
						"error": "endpoint_not_registered",
					},
				)
				return
			}

			c.AbortWithStatusJSON(
				http.StatusInternalServerError,
				gin.H{
					"error": "endpoint_lookup_failed",
				},
			)
			return
		}

		// =====================================================
		// ENDPOINT DOES NOT REQUIRE AUTHORIZATION
		// =====================================================

		if !endpoint.RequiresAuth {
			c.Next()
			return
		}

		// =====================================================
		// CHECK IDENTITY PERMISSION
		// =====================================================

		allowed, err := accessService.HasPermission(
			c.Request.Context(),
			authContext.IdentityID,
			endpoint.ID,
		)

		if err != nil {
			c.AbortWithStatusJSON(
				http.StatusInternalServerError,
				gin.H{
					"error": "authorization_check_failed",
				},
			)
			return
		}

		if !allowed {
			c.AbortWithStatusJSON(
				http.StatusForbidden,
				gin.H{
					"error": "access_denied",
				},
			)
			return
		}

		c.Next()
	}
}
