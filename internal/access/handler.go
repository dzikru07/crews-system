package access

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

func (h *Handler) AssignPermission(c *gin.Context) {

	endpointID, err := uuid.Parse(
		c.Param("endpoint_id"),
	)

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "invalid_endpoint_id",
		})
		return
	}

	var req AssignPermissionRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "invalid_request",
		})
		return
	}

	err = h.service.AssignPermission(
		c.Request.Context(),
		endpointID,
		req.PermissionID,
	)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "assign_permission_failed",
		})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"endpoint_id":   endpointID,
		"permission_id": req.PermissionID,
	})
}

func (h *Handler) RemovePermission(c *gin.Context) {

	endpointID, err := uuid.Parse(
		c.Param("endpoint_id"),
	)

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "invalid_endpoint_id",
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

	err = h.service.RemovePermission(
		c.Request.Context(),
		endpointID,
		permissionID,
	)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "remove_permission_failed",
		})
		return
	}

	c.Status(http.StatusNoContent)
}

func (h *Handler) ListPermissions(c *gin.Context) {

	endpointID, err := uuid.Parse(
		c.Param("endpoint_id"),
	)

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "invalid_endpoint_id",
		})
		return
	}

	permissions, err := h.service.ListPermissions(
		c.Request.Context(),
		endpointID,
	)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "list_permissions_failed",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"endpoint_id": endpointID,
		"permissions": permissions,
	})
}
