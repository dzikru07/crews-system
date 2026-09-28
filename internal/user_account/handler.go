package user_account

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type Handler struct {
	service *Service
}

func NewHandler(service *Service) *Handler {
	return &Handler{
		service: service,
	}
}

func (h *Handler) Create(c *gin.Context) {
	tenantID, err := uuid.Parse(
		c.Param("tenant_id"),
	)

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "invalid_tenant_id",
		})
		return
	}

	identityID, err := uuid.Parse(
		c.Param("identity_id"),
	)

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "invalid_identity_id",
		})
		return
	}

	var req CreateRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "invalid_request",
		})
		return
	}

	response, err := h.service.Create(
		c.Request.Context(),
		tenantID,
		identityID,
		req,
	)

	if err != nil {
		switch err {

		case ErrInvalidEmail:
			c.JSON(http.StatusBadRequest, gin.H{
				"error": "invalid_email",
			})

		case ErrInvalidPassword:
			c.JSON(http.StatusBadRequest, gin.H{
				"error": "invalid_password",
			})

		case ErrInvalidUsername:
			c.JSON(http.StatusBadRequest, gin.H{
				"error": "invalid_username",
			})

		case ErrIdentityNotFound:
			c.JSON(http.StatusNotFound, gin.H{
				"error": "identity_not_found",
			})

		case ErrIdentityNotHuman:
			c.JSON(http.StatusBadRequest, gin.H{
				"error": "identity_must_be_human",
			})

		case ErrIdentityNotActive:
			c.JSON(http.StatusForbidden, gin.H{
				"error": "identity_not_active",
			})

		case ErrUserAlreadyExists:
			c.JSON(http.StatusConflict, gin.H{
				"error": "user_account_already_exists",
			})

		default:
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": "user_account_creation_failed",
			})
		}

		return
	}

	c.JSON(http.StatusCreated, response)
}
