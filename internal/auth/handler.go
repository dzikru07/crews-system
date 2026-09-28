package auth

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type Handler struct {
	service *Service
}

func NewHandler(service *Service) *Handler {
	return &Handler{
		service: service,
	}
}

func (h *Handler) Token(c *gin.Context) {
	var req TokenRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "invalid_request",
		})
		return
	}

	environment := c.GetHeader("X-Crews-Environment")

	if environment == "" {
		environment = "development"
	}

	response, err := h.service.Token(
		c.Request.Context(),
		req,
		environment,
	)
	if err != nil {
		switch err {
		case ErrUnsupportedGrantType:
			c.JSON(http.StatusBadRequest, gin.H{
				"error": "unsupported_grant_type",
			})

		case ErrInvalidClient:
			c.JSON(http.StatusUnauthorized, gin.H{
				"error": "invalid_client",
			})

		case ErrCredentialInactive:
			c.JSON(http.StatusUnauthorized, gin.H{
				"error": "credential_inactive",
			})

		case ErrIdentityInactive:
			c.JSON(http.StatusUnauthorized, gin.H{
				"error": "identity_inactive",
			})

		case ErrCredentialExpired:
			c.JSON(http.StatusUnauthorized, gin.H{
				"error": "credential_expired",
			})

		default:
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": "token_generation_failed",
			})
		}

		return
	}

	c.JSON(http.StatusOK, response)
}

func (h *Handler) Login(c *gin.Context) {
	var req LoginRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "invalid_request",
		})
		return
	}

	environment := c.GetHeader("X-Crews-Environment")

	if environment == "" {
		environment = "development"
	}

	response, err := h.service.Login(
		c.Request.Context(),
		req,
		environment,
	)

	if err != nil {
		switch err {
		case ErrInvalidLogin:
			c.JSON(http.StatusUnauthorized, gin.H{
				"error": "invalid_credentials",
			})

		default:
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": "login_failed",
			})
		}

		return
	}

	c.JSON(http.StatusOK, response)
}
