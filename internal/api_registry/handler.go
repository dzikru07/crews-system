package api_registry

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

func (h *Handler) CreateService(c *gin.Context) {
	applicationID, err := uuid.Parse(
		c.Param("application_id"),
	)

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "invalid_application_id",
		})
		return
	}

	var req CreateServiceRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "invalid_request",
		})
		return
	}

	service, err := h.service.CreateService(
		c.Request.Context(),
		applicationID,
		req,
	)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "api_service_creation_failed",
		})
		return
	}

	c.JSON(http.StatusCreated, service)
}

func (h *Handler) CreateEndpoint(c *gin.Context) {
	serviceID, err := uuid.Parse(
		c.Param("service_id"),
	)

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "invalid_service_id",
		})
		return
	}

	var req CreateEndpointRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "invalid_request",
		})
		return
	}

	endpoint, err := h.service.CreateEndpoint(
		c.Request.Context(),
		serviceID,
		req,
	)

	if err != nil {
		switch err {
		case ErrInvalidMethod:
			c.JSON(http.StatusBadRequest, gin.H{
				"error": "invalid_http_method",
			})

		case ErrInvalidPath:
			c.JSON(http.StatusBadRequest, gin.H{
				"error": "invalid_endpoint_path",
			})

		default:
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": "api_endpoint_creation_failed",
			})
		}

		return
	}

	c.JSON(http.StatusCreated, endpoint)
}
