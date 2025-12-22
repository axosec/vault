package api

import (
	"net/http"

	"github.com/axosec/vault/internal/dto"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

// DeleteResourceHandler godoc
// @Summary      Delete Resource
// @Description  Soft delete a folder or an item.
// @Tags         Management
// @Param        type path      string true "Resource Type (folder/item)"
// @Param        id   path      string true "Resource UUID"
// @Success      204  {object}  nil
// @Failure      400  {object}  map[string]string "Invalid request"
// @Failure      500  {object}  map[string]string "Internal server error"
// @Router       /resources/{type}/{id} [delete]
func (h *Handler) DeleteResourceHandler(c *gin.Context) {
	idStr := c.Param("id")
	resourceID, err := uuid.Parse(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid resource ID"})
		return
	}

	typeStr := c.Param("type")
	var resourceType dto.ResourceType
	switch typeStr {
	case "folder":
		resourceType = dto.TypeFolder
	case "item":
		resourceType = dto.TypeItem
	default:
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid resource type. Must be 'folder' or 'item'"})
		return
	}

	userID := c.MustGet("user_id").(uuid.UUID)

	if err := h.vaultService.DeleteResource(c.Request.Context(), userID, resourceID, resourceType); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete resource"})
		return
	}

	c.Status(http.StatusNoContent)
}

// ShareResourceHandler godoc
// @Summary      Share Resource
// @Description  Grant another user access to a resource by adding a new encrypted key.
// @Tags         Management
// @Accept       json
// @Param        request body dto.ShareParams true "Share details"
// @Success      200  {string}  string "OK"
// @Failure      400  {object}  map[string]string "Invalid request"
// @Failure      500  {object}  map[string]string "Internal server error"
// @Router       /share [post]
func (h *Handler) ShareResourceHandler(c *gin.Context) {
	var req dto.ShareParams
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid share parameters"})
		return
	}

	userID := c.MustGet("user_id").(uuid.UUID)

	if err := h.vaultService.ShareResource(c.Request.Context(), userID, req); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to share resource"})
		return
	}

	c.Status(http.StatusOK)
}

// RevokeAccessHandler godoc
// @Summary      Revoke Access
// @Description  Remove a specific user's access to a resource.
// @Tags         Management
// @Accept       json
// @Param        request body dto.RevokeReq true "Revocation details"
// @Success      200  {string}  string "OK"
// @Failure      400  {object}  map[string]string "Invalid request"
// @Failure      500  {object}  map[string]string "Internal server error"
// @Router       /share/revoke [post]
func (h *Handler) RevokeAccessHandler(c *gin.Context) {
	type RevokeReq struct {
		TargetUserID uuid.UUID `json:"target_user_id" binding:"required"`
		ResourceID   uuid.UUID `json:"resource_id" binding:"required"`
	}

	var req RevokeReq
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid parameters"})
		return
	}

	userID := c.MustGet("user_id").(uuid.UUID)

	if err := h.vaultService.RevokeAccess(c.Request.Context(), userID, req.TargetUserID, req.ResourceID); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to revoke access"})
		return
	}

	c.Status(http.StatusOK)
}
