package api

import (
	"github.com/axosec/core/crypto/token"
	"github.com/axosec/vault/internal/service"
	"github.com/gin-gonic/gin"
)

type Handler struct {
	jwt          *token.JWTManager
	vaultService *service.VaultService
}

func NewHandler(jwt *token.JWTManager, vaultService *service.VaultService) *Handler {
	return &Handler{
		jwt:          jwt,
		vaultService: vaultService,
	}
}

// @Summary      Helthcheck
// @Description  returns ok if api up
// @Success      200
// @Router       /health [get]
func (h *Handler) Helth(c *gin.Context) {
	c.JSON(200, gin.H{"status": "ok"})
}

func (h *Handler) RegisterRouters(e *gin.Engine) {
	v1 := e.Group("/v1")

	v1.GET("/health", h.Helth)

	protected := v1.Group("/")
	protected.Use(h.AuthenticatedMiddleware())
	{
		protected.POST("/folders", h.CreateFolderHandler)
		protected.GET("/folders", h.ListFoldersHandler)
		protected.PUT("/folders/:id", h.UpdateFolderHandler)

		protected.POST("/items", h.CreateItemHandler)
		protected.GET("/items", h.ListItemsHandler)
		protected.GET("/items/:id", h.GetItemHandler)
		protected.PUT("/items/:id", h.UpdateItemHandler)

		protected.DELETE("/resources/:type/:id", h.DeleteResourceHandler)
		protected.POST("/share", h.ShareResourceHandler)
		protected.POST("/share/revoke", h.RevokeAccessHandler)
	}
}
