package main

import (
	"log"
	"os"

	"github.com/gin-gonic/gin"

	"github.com/PosokhovVadim/stawberry/internal/app"
	"github.com/PosokhovVadim/stawberry/internal/config"
	"github.com/PosokhovVadim/stawberry/internal/domain/service/auth"
	"github.com/PosokhovVadim/stawberry/internal/domain/service/offer"
	"github.com/PosokhovVadim/stawberry/internal/domain/service/product"
	"github.com/PosokhovVadim/stawberry/internal/domain/service/session"
	"github.com/PosokhovVadim/stawberry/internal/handler"
	"github.com/PosokhovVadim/stawberry/internal/repository"
	repo_session "github.com/PosokhovVadim/stawberry/internal/repository/session"
	"github.com/PosokhovVadim/stawberry/pkg/hasher"
	objectstorage "github.com/PosokhovVadim/stawberry/pkg/s3"
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

	// initialize password hasher
	bcryptHasher, err := hasher.NewBcryptHasher(&hasher.Config{Cost: hasher.MinCost})
	if err != nil {
		log.Fatalf("failed to initialize bcrypt hasher: %s", err.Error())
	}

	// Initialize database connection
	db := repository.InitDB(cfg)
	productRepository := repository.NewProductRepository(db)
	offerRepository := repository.NewOfferRepository(db)
	userRepository := repository.NewUserRepository(db)
	sessionRepository := repo_session.New(db)

	productService := product.NewProductService(productRepository)
	offerService := offer.NewOfferService(offerRepository)
	sessionManager := session.New(&session.Dependencies{
		SessionRepository: sessionRepository,
	})
	authService := auth.New(&auth.Dependencies{
		Config:         auth.DefaultConfig,
		UserRepo:       userRepository,
		SessionManager: sessionManager,
		PasswordHasher: bcryptHasher,
		PasswordPolicy: auth.DefaultPasswordPolicy,
	})

	// Initialize object storage s3
	s3 := objectstorage.ObjectStorageConn(cfg)

	// Initialize router
	router = handler.SetupRouter(productService, offerService, s3, authService, sessionManager)

	return nil
}
