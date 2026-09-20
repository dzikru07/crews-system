package credential

import (
	"net/http"

	"github.com/dzikru07/crews-system/internal/identity"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type Handler struct {
	service      *Service
	identityRepo *identity.Repository
}

func NewHandler(
	service *Service,
	identityRepo *identity.Repository,
) *Handler {
	return &Handler{
		service:      service,
		identityRepo: identityRepo,
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

	identityID, err := uuid.Parse(c.Param("identity_id"))
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

	identityRecord, err := h.identityRepo.GetByID(
		c.Request.Context(),
		identityID,
		tenantID,
	)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"error": "identity_not_found",
		})
		return
	}

	if identityRecord.Type != identity.TypeMachine {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "identity_must_be_machine",
		})
		return
	}

	if identityRecord.Status != identity.StatusActive {
		c.JSON(http.StatusForbidden, gin.H{
			"error": "identity_not_active",
		})
		return
	}

	prefix := environmentPrefix(c)

	response, err := h.service.Create(
		c.Request.Context(),
		identityID,
		req,
		prefix,
	)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "credential_creation_failed",
		})
		return
	}

	c.JSON(http.StatusCreated, response)
}

func environmentPrefix(c *gin.Context) string {
	environment := c.GetHeader("X-Crews-Environment")

	switch environment {
	case "production":
		return "live"

	case "staging":
		return "stg"

	case "development":
		return "dev"

	default:
		return "dev"
	}
}
