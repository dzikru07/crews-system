package identity_role

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

func (h *Handler) AssignRole(c *gin.Context) {

	identityID, err := uuid.Parse(
		c.Param("identity_id"),
	)

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "invalid_identity_id",
		})
		return
	}

	var req AssignRoleRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "invalid_request",
		})
		return
	}

	// if err := h.service.AssignRole(
	// 	c.Request.Context(),
	// 	identityID,
	// 	req.RoleID,
	// ); err != nil {
	// 	c.JSON(http.StatusInternalServerError, gin.H{
	// 		"error": "assign_role_failed",
	// 	})
	// 	return
	// }

	if err := h.service.AssignRole(
		c.Request.Context(),
		identityID,
		req.RoleID,
	); err != nil {

		if errors.Is(err, ErrInvalidAssignment) {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": "invalid_role_assignment",
			})
			return
		}

		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "assign_role_failed",
		})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"identity_id": identityID,
		"role_id":     req.RoleID,
	})
}

func (h *Handler) RemoveRole(c *gin.Context) {

	identityID, err := uuid.Parse(
		c.Param("identity_id"),
	)

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "invalid_identity_id",
		})
		return
	}

	roleID, err := uuid.Parse(
		c.Param("role_id"),
	)

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "invalid_role_id",
		})
		return
	}

	if err := h.service.RemoveRole(
		c.Request.Context(),
		identityID,
		roleID,
	); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "remove_role_failed",
		})
		return
	}

	c.Status(http.StatusNoContent)
}

func (h *Handler) ListRoles(c *gin.Context) {

	identityID, err := uuid.Parse(
		c.Param("identity_id"),
	)

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "invalid_identity_id",
		})
		return
	}

	roles, err := h.service.ListRoles(
		c.Request.Context(),
		identityID,
	)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "list_roles_failed",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"identity_id": identityID,
		"roles":       roles,
	})
}
