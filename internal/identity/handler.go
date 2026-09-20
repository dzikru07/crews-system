package identity

import (
	"errors"
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
	tenantID, err := uuid.Parse(c.Param("tenant_id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "invalid_tenant_id",
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

	identity, err := h.service.Create(
		c.Request.Context(),
		tenantID,
		req,
	)
	if err != nil {
		c.Error(err)

		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "internal_server_error",
		})
		return
	}

	c.JSON(http.StatusCreated, identity)
}

func (h *Handler) GetByID(c *gin.Context) {
	tenantID, err := uuid.Parse(c.Param("tenant_id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "invalid_tenant_id",
		})
		return
	}

	identityID, err := uuid.Parse(c.Param("identity_id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "invalid_identity_id",
		})
		return
	}

	identity, err := h.service.GetByID(
		c.Request.Context(),
		identityID,
		tenantID,
	)

	if errors.Is(err, ErrNotFound) {
		c.JSON(http.StatusNotFound, gin.H{
			"error": "identity_not_found",
		})
		return
	}

	if err != nil {
		c.Error(err)

		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "internal_server_error",
		})
		return
	}

	c.JSON(http.StatusOK, identity)
}
