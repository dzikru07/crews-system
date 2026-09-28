package role_permission

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type Handler struct {
	service *Service
}

func NewHandler(
	service *Service,
) *Handler {
	return &Handler{
		service: service,
	}
}

func (h *Handler) AssignPermission(c *gin.Context) {

	roleID, err := uuid.Parse(
		c.Param("role_id"),
	)

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "invalid_role_id",
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

	// if err := h.service.AssignPermission(
	// 	c.Request.Context(),
	// 	roleID,
	// 	req.PermissionID,
	// ); err != nil {
	// 	c.JSON(http.StatusInternalServerError, gin.H{
	// 		"error": "assign_permission_failed",
	// 	})
	// 	return
	// }
	if err := h.service.AssignPermission(
		c.Request.Context(),
		roleID,
		req.PermissionID,
	); err != nil {

		if errors.Is(err, ErrInvalidAssignment) {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": "invalid_permission_assignment",
			})
			return
		}

		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "assign_permission_failed",
		})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"role_id":       roleID,
		"permission_id": req.PermissionID,
	})
}

func (h *Handler) RemovePermission(c *gin.Context) {

	roleID, err := uuid.Parse(
		c.Param("role_id"),
	)

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "invalid_role_id",
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

	if err := h.service.RemovePermission(
		c.Request.Context(),
		roleID,
		permissionID,
	); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "remove_permission_failed",
		})
		return
	}

	c.Status(http.StatusNoContent)
}

func (h *Handler) ListPermissions(c *gin.Context) {

	roleID, err := uuid.Parse(
		c.Param("role_id"),
	)

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "invalid_role_id",
		})
		return
	}

	permissions, err := h.service.ListPermissions(
		c.Request.Context(),
		roleID,
	)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "list_permissions_failed",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"role_id":     roleID,
		"permissions": permissions,
	})
}
