package application

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

type createApplicationRequest struct {
	Name        string  `json:"name" binding:"required"`
	Slug        string  `json:"slug" binding:"required"`
	Description *string `json:"description"`
}

func (h *Handler) Create(c *gin.Context) {
	tenantID, err := uuid.Parse(c.Param("tenant_id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "invalid_tenant_id",
		})
		return
	}

	var req createApplicationRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "invalid_request",
		})
		return
	}

	app, err := h.service.Create(
		c.Request.Context(),
		CreateInput{
			TenantID:    tenantID,
			Name:        req.Name,
			Slug:        req.Slug,
			Description: req.Description,
		},
	)

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	c.JSON(http.StatusCreated, app)
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

	applicationID, err := uuid.Parse(c.Param("application_id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "invalid_application_id",
		})
		return
	}

	app, err := h.service.GetByIDAndTenant(
		c.Request.Context(),
		applicationID,
		tenantID,
	)

	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"error": "application_not_found",
		})
		return
	}

	c.JSON(http.StatusOK, app)
}
