package permission

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

	var req CreateRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "invalid_request",
		})
		return
	}

	result, err := h.service.Create(
		c.Request.Context(),
		tenantID,
		req,
	)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "create_permission_failed",
		})
		return
	}

	c.JSON(http.StatusCreated, result)
}

func (h *Handler) GetByID(c *gin.Context) {

	tenantID, err := uuid.Parse(
		c.Param("tenant_id"),
	)

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "invalid_tenant_id",
		})
		return
	}

	permissionID, err := uuid.Parse(
		c.Param("permission_id"),
	)

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "invalid_permission_id",
		})
		return
	}

	result, err := h.service.GetByID(
		c.Request.Context(),
		tenantID,
		permissionID,
	)

	if err != nil {
		if err == ErrPermissionNotFound {
			c.JSON(http.StatusNotFound, gin.H{
				"error": "permission_not_found",
			})
			return
		}

		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "get_permission_failed",
		})
		return
	}

	c.JSON(http.StatusOK, result)
}
