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
	{
		v1.GET("/health", h.Helth)

	}
}
