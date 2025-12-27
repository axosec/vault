package main

import (
	"fmt"
	"log"
	"time"

	"github.com/axosec/core/crypto/token"
	_ "github.com/axosec/vault/docs"
	"github.com/axosec/vault/internal/api"
	"github.com/axosec/vault/internal/config"
	"github.com/axosec/vault/internal/data/db"
	"github.com/axosec/vault/internal/service"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/swaggo/files"
	"github.com/swaggo/gin-swagger"
)

var (
	Version   = "dev"
	GitCommit = "unknown"
	BuildDate = "unknown"
)

// @title           Axosec vault API
// @version         0.1.0
// @description     Swagger definitions for axosec vault http api

// @host      localhost:8080
// @BasePath  /v1

func main() {
	// Version info
	fmt.Printf(`Starting Axosec Vault
Version: %s
Commit: %s
Built:  %s

`, Version, GitCommit, BuildDate)
	// Load config
	cfg, err := config.LoadConfig()
	if err != nil {
		fmt.Printf("failed to load config: %s\n", err)
		return
	}

	// Setup keys and JWT
	privateKey, publicKey, err := token.LoadKeysFromFiles(cfg.JWT.PrivateKeyPath, cfg.JWT.PublicKeyPath)
	if err != nil {
		fmt.Println(err)
		return
	}

	jwtManager := token.NewJWTManager(privateKey, publicKey, cfg.JWT.Issuer)

	// Setup db connection
	connPool, err := db.NewConnection(cfg.Database)
	if err != nil {
		fmt.Printf("failed to connect to database: %s", err)
		return
	}

	queries := db.New(connPool)

	// Initialize services
	vaultService := service.NewVaultService(connPool, queries)

	// Start http router
	apiHandler := api.NewHandler(jwtManager, vaultService)

	r := gin.Default()
	r.Use(cors.New(cors.Config{
		AllowOrigins: []string{"http://localhost:5173", "http://localhost:5174"},
		AllowMethods: []string{"GET", "POST", "PUT", "PATCH", "DELETE"},
		AllowHeaders: []string{
			"Content-Type",
		},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}))
	apiHandler.RegisterRouters(r)
	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
	if err := r.Run(":" + cfg.ServerPort); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
