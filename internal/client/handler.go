package client

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

type createClientRequest struct {
	Name string `json:"name" binding:"required"`
	Slug string `json:"slug" binding:"required"`
	Type string `json:"type" binding:"required"`
}

func (h *Handler) Create(c *gin.Context) {
	environmentID, err := uuid.Parse(
		c.Param("environment_id"),
	)

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "invalid_environment_id",
		})
		return
	}

	var req createClientRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "invalid_request",
		})
		return
	}

	client, err := h.service.Create(
		c.Request.Context(),
		CreateInput{
			EnvironmentID: environmentID,
			Name:          req.Name,
			Slug:          req.Slug,
			Type:          req.Type,
		},
	)

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	c.JSON(http.StatusCreated, client)
}

func (h *Handler) GetByID(c *gin.Context) {
	clientID, err := uuid.Parse(
		c.Param("client_id"),
	)

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "invalid_client_id",
		})
		return
	}

	environmentID, err := uuid.Parse(
		c.Param("environment_id"),
	)

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "invalid_environment_id",
		})
		return
	}

	// client, err := h.service.GetByID(
	// 	c.Request.Context(),
	// 	clientID,
	// )

	client, err := h.service.GetByIDAndEnvironment(
		c.Request.Context(),
		clientID,
		environmentID,
	)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"error": "client_not_found",
		})
		return
	}

	c.JSON(http.StatusOK, client)
}
