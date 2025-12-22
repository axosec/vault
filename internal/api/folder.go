package api

import (
	"net/http"

	"github.com/axosec/vault/internal/dto"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

// CreateFolderHandler godoc
// @Summary      Create Folder
// @Description  Creates a new encrypted folder.
// @Tags         Folders
// @Accept       json
// @Produce      json
// @Param        request body dto.CreateFolderReq true "Folder creation details"
// @Success      201  {object}  dto.FolderResponse
// @Failure      400  {object}  map[string]string "Invalid request"
// @Failure      500  {object}  map[string]string "Internal server error"
// @Router       /folders [post]
func (h *Handler) CreateFolderHandler(c *gin.Context) {
	var req dto.CreateFolderReq
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request: " + err.Error()})
		return
	}

	userID := c.MustGet("user_id").(uuid.UUID)

	folder, err := h.vaultService.CreateFolder(c.Request.Context(), userID, req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create folder"})
		return
	}

	c.JSON(http.StatusCreated, folder)
}

// ListFoldersHandler godoc
// @Summary      List Folders
// @Description  Get a flat of all folders the user has access to.
// @Tags         Folders
// @Produce      json
// @Success      200  {array}   dto.FolderSummary
// @Failure      500  {object}  map[string]string "Internal server error"
// @Router       /folders [get]
func (h *Handler) ListFoldersHandler(c *gin.Context) {
	userID := c.MustGet("user_id").(uuid.UUID)

	folders, err := h.vaultService.ListFolders(c.Request.Context(), userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch folders"})
		return
	}

	c.JSON(http.StatusOK, folders)
}
