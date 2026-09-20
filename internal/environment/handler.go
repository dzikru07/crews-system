package environment

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

type createEnvironmentRequest struct {
	Name string `json:"name" binding:"required"`
	Type string `json:"type" binding:"required"`
}

func (h *Handler) Create(c *gin.Context) {
	applicationID, err := uuid.Parse(
		c.Param("application_id"),
	)

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "invalid_application_id",
		})
		return
	}

	var req createEnvironmentRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "invalid_request",
		})
		return
	}

	env, err := h.service.Create(
		c.Request.Context(),
		CreateInput{
			ApplicationID: applicationID,
			Name:          req.Name,
			Type:          req.Type,
		},
	)

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	c.JSON(http.StatusCreated, env)
}

func (h *Handler) GetByID(c *gin.Context) {
	environmentID, err := uuid.Parse(
		c.Param("environment_id"),
	)

	applicationID, err := uuid.Parse(
		c.Param("application_id"),
	)

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "invalid_environment_id",
		})
		return
	}

	// env, err := h.service.GetByID(
	// 	c.Request.Context(),
	// 	environmentID,
	// )
	env, err := h.service.GetByIDAndApplication(
		c.Request.Context(),
		environmentID,
		applicationID,
	)

	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"error": "environment_not_found",
		})
		return
	}

	c.JSON(http.StatusOK, env)
}
