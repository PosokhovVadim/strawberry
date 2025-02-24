package main

import (
	"log"
	"os"

	"github.com/PosokhovVadim/stawberry/internal/repository"

	"github.com/PosokhovVadim/stawberry/internal/app"
	"github.com/PosokhovVadim/stawberry/internal/config"
	"github.com/PosokhovVadim/stawberry/internal/domain/service/offer"
	"github.com/PosokhovVadim/stawberry/internal/domain/service/product"
	"github.com/PosokhovVadim/stawberry/internal/handler"
	objectstorage "github.com/PosokhovVadim/stawberry/pkg/s3"
	"github.com/gin-gonic/gin"
)

// Global variables for application state
var (
	router *gin.Engine
)

func main() {
	// Initialize application
	if err := initializeApp(); err != nil {
		log.Fatalf("Failed to initialize application: %v", err)
	}

	// Get port from environment variable or use default
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	// Start server
	if err := app.StartServer(router, port); err != nil {
		log.Fatalf("Server error: %v", err)
	}
}

// initializeApp initializes all application components
func initializeApp() error {
	// Load configuration
	cfg := config.LoadConfig()

	// Set Gin mode based on environment
	if os.Getenv("GIN_MODE") == "release" {
		gin.SetMode(gin.ReleaseMode)
	}

	// Initialize database connection
	db := repository.InitDB(cfg)
	productRepository := repository.NewProductRepository(db)
	offerRepository := repository.NewOfferRepository(db)

	productService := product.NewProductService(productRepository)
	offerService := offer.NewOfferService(offerRepository)

	// Initialize object storage s3
	s3 := objectstorage.ObjectStorageConn(cfg)

	// Initialize router
	router = handler.SetupRouter(productService, offerService, s3)

	return nil
}
