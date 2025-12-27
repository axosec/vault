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

// UpdateFolderHandler godoc
// @Summary      Update Folder Metadata
// @Description  Update the encrypted name/icon/color blob for a folder.
// @Tags         Folders
// @Accept       json
// @Param        id   path      string true "Folder UUID"
// @Param        request body dto.UpdateFolderReq true "New Encrypted Metadata"
// @Success      200  {string}  string "OK"
// @Failure      400  {object}  map[string]string "Invalid request"
// @Failure      500  {object}  map[string]string "Internal server error"
// @Router       /folders/{id} [put]
func (h *Handler) UpdateFolderHandler(c *gin.Context) {
	idStr := c.Param("id")
	folderID, err := uuid.Parse(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid folder ID"})
		return
	}

	var req dto.UpdateFolderReq
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request: " + err.Error()})
		return
	}

	userID := c.MustGet("user_id").(uuid.UUID)

	if err := h.vaultService.UpdateFolder(c.Request.Context(), userID, folderID, req); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update folder"})
		return
	}

	c.Status(http.StatusOK)
}
