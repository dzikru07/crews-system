package tenant

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

type createTenantRequest struct {
	Name string `json:"name" binding:"required"`
	Slug string `json:"slug" binding:"required"`
	Type string `json:"type" binding:"required"`
}

func (h *Handler) Create(c *gin.Context) {
	var req createTenantRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "invalid_request",
		})
		return
	}

	tenant, err := h.service.Create(
		c.Request.Context(),
		CreateInput{
			Name: req.Name,
			Slug: req.Slug,
			Type: req.Type,
		},
	)

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	c.JSON(http.StatusCreated, tenant)
}

func (h *Handler) GetByID(c *gin.Context) {
	id, err := uuid.Parse(c.Param("tenant_id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "invalid_tenant_id",
		})
		return
	}

	tenant, err := h.service.GetByID(
		c.Request.Context(),
		id,
	)

	if err != nil {
		if err == ErrNotFound {
			c.JSON(http.StatusNotFound, gin.H{
				"error": "tenant_not_found",
			})
			return
		}

		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "internal_server_error",
		})
		return
	}

	c.JSON(http.StatusOK, tenant)
}
