package api

import (
	"net/http"

	"github.com/axosec/vault/internal/dto"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

// CreateItemHandler godoc
// @Summary      Create Item
// @Description  Creates a new encrypted item.
// @Tags         Items
// @Accept       json
// @Produce      json
// @Param        request body dto.CreateItemReq true "Item creation details"
// @Success      201  {object}  dto.ItemResponse
// @Failure      400  {object}  map[string]string "Invalid request"
// @Failure      500  {object}  map[string]string "Internal server error"
// @Router       /items [post]
func (h *Handler) CreateItemHandler(c *gin.Context) {
	var req dto.CreateItemReq
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request parameters"})
		return
	}

	userID := c.MustGet("user_id").(uuid.UUID)

	item, err := h.vaultService.CreateItem(c.Request.Context(), userID, req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create item"})
		return
	}

	c.JSON(http.StatusCreated, item)
}

// ListItemsHandler godoc
// @Summary      List Items
// @Description  List items within a specific folder. If folder_id is missing, lists root items.
// @Tags         Items
// @Produce      json
// @Param        folder_id query string false "Folder UUID (optional)"
// @Success      200  {array}   dto.ItemSummary
// @Failure      400  {object}  map[string]string "Invalid UUID"
// @Failure      500  {object}  map[string]string "Internal server error"
// @Router       /items [get]
func (h *Handler) ListItemsHandler(c *gin.Context) {
	userID := c.MustGet("user_id").(uuid.UUID)

	var folderID uuid.UUID
	folderIDStr := c.Query("folder_id")

	if folderIDStr != "" {
		parsed, err := uuid.Parse(folderIDStr)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid folder_id format"})
			return
		}
		folderID = parsed
	}

	items, err := h.vaultService.ListItems(c.Request.Context(), userID, folderID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch items"})
		return
	}

	c.JSON(http.StatusOK, items)
}

// GetItemHandler godoc
// @Summary      Get Item Details
// @Description  Fetch the full encrypted blob and keys for a specific item.
// @Tags         Items
// @Produce      json
// @Param        id   path      string true "Item UUID"
// @Success      200  {object}  dto.ItemDetail
// @Failure      400  {object}  map[string]string "Invalid UUID"
// @Failure      404  {object}  map[string]string "Item not found"
// @Failure      500  {object}  map[string]string "Internal server error"
// @Router       /items/{id} [get]
func (h *Handler) GetItemHandler(c *gin.Context) {
	idStr := c.Param("id")
	itemID, err := uuid.Parse(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid item ID"})
		return
	}

	userID := c.MustGet("user_id").(uuid.UUID)

	item, err := h.vaultService.GetItem(c.Request.Context(), userID, itemID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Item not found or access denied"})
		return
	}

	c.JSON(http.StatusOK, item)
}

// UpdateItemHandler godoc
// @Summary      Update Item
// @Description  Update the encrypted data blob of an item.
// @Tags         Items
// @Accept       json
// @Param        id   path      string true "Item UUID"
// @Param        request body dto.UpdateItemReq true "Update payload"
// @Success      200  {string}  string "OK"
// @Failure      400  {object}  map[string]string "Invalid request"
// @Failure      500  {object}  map[string]string "Internal server error"
// @Router       /items/{id} [put]
func (h *Handler) UpdateItemHandler(c *gin.Context) {
	idStr := c.Param("id")
	itemID, err := uuid.Parse(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid item ID"})
		return
	}

	var req dto.UpdateItemReq
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request parameters"})
		return
	}

	userID := c.MustGet("user_id").(uuid.UUID)

	if err := h.vaultService.UpdateItem(c.Request.Context(), userID, itemID, req); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update item"})
		return
	}

	c.Status(http.StatusOK)
}
