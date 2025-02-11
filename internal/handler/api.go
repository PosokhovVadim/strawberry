package handler

import (
	"errors"
	"net/http"
	"time"

	"github.com/PosokhovVadim/stawberry/internal/app/apperror"
	"github.com/PosokhovVadim/stawberry/internal/handler/middleware"
	objectstorage "github.com/PosokhovVadim/stawberry/pkg/s3"

	"github.com/gin-gonic/gin"
)

func SetupRouter(
	productService ProductService,
	offerService OfferService,
	s3 *objectstorage.BucketBasics,
	authService AuthService,
	sessionManager middleware.SessionManager,
) *gin.Engine {
	router := gin.New()

	// Add default middleware
	router.Use(gin.Logger())
	router.Use(gin.Recovery())
	router.Use(middleware.CORS())

	// Health check endpoint
	router.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"status": "ok",
			"time":   time.Now().Unix(),
		})
	})

	// API routes group
	api := router.Group("/api")
	{
		authHandler := NewAuthHandler(authService)

		// Public routes
		public := api.Group("")
		{
			// Auth endpoints
			public.POST("/auth/register", authHandler.Register)
			public.POST("/auth/login", authHandler.Login)
			public.POST("/auth/refresh", authHandler.RefreshTokens)

			// Public product search
			// public.GET("/products/search", handlers.SearchProducts(db))
			// public.GET("/stores", handlers.GetStores(db))
		}

		// Protected routes
		protected := api.Group("")
		protected.Use(middleware.AuthMiddleware(sessionManager))
		{
			protected.GET("/auth/me", authHandler.Me)
			// User profile
			// protected.GET("/profile", handlers.GetProfile(db))
			// protected.PUT("/profile", handlers.UpdateProfile(db))

			// Store management
			// stores := protected.Group("/stores")
			// {
			// 	stores.GET("/:id", handlers.GetStore(db))
			// 	stores.GET("/:id/products", handlers.GetStoreProducts(db))
			// }

			// Product management
			// products := protected.Group("/products")
			// {
			// 	products.GET("", handlers.GetProducts(db))
			// 	products.GET("/:id", handlers.GetProduct(db))
			// 	products.PUT("/:id", handlers.UpdateProduct(db))
			// 	products.POST("", handlers.AddProduct(db))
			// }

			// Offer management
			// offers := protected.Group("/offers")
			// {
			// 	offers.POST("", handlers.CreateOffer(db))
			// 	offers.GET("", handlers.GetUserOffers(db))
			// 	offers.GET("/:id", handlers.GetOffer(db))
			// 	offers.PUT("/:id/status", handlers.UpdateOfferStatus(db))
			// 	offers.DELETE("/:id", handlers.CancelOffer(db))
			// }

			// Notification management
			// notifications := protected.Group("/notifications")
			// {
			// notifications.GET("", handlers.GetNotifications(db))
			// notifications.PUT("/:id/read", handlers.MarkNotificationRead(db))
			// notifications.DELETE("/:id", handlers.DeleteNotification(db))
			// }
		}
	}

	return router
}

func handleProductError(c *gin.Context, err error) {
	var productErr *apperror.ProductError
	if errors.As(err, &productErr) {
		status := http.StatusInternalServerError

		switch productErr.Code {
		case apperror.NotFound:
			status = http.StatusNotFound
		case apperror.DuplicateError:
			status = http.StatusConflict
		case apperror.DatabaseError:
			status = http.StatusInternalServerError
		}

		c.JSON(status, gin.H{
			"code":    productErr.Code,
			"message": productErr.Message,
		})
		return
	}

	c.JSON(http.StatusInternalServerError, gin.H{
		"code":    apperror.InternalError,
		"message": "An unexpected error occurred",
	})
}

func handleOfferError(c *gin.Context, err error) {
	var offerError *apperror.OfferError
	if errors.As(err, &offerError) {
		status := http.StatusInternalServerError

		switch offerError.Code {
		case apperror.NotFound:
			status = http.StatusNotFound
		case apperror.DuplicateError:
			status = http.StatusConflict
		case apperror.DatabaseError:
			status = http.StatusInternalServerError
		}

		c.JSON(status, gin.H{
			"code":    offerError.Code,
			"message": offerError.Message,
		})
		return
	}

	c.JSON(http.StatusInternalServerError, gin.H{
		"code":    apperror.InternalError,
		"message": "An unexpected error occurred",
	})
}

func handleBindError(c *gin.Context, err error) {
	c.JSON(http.StatusBadRequest, gin.H{"code": apperror.BadRequest, "message": err.Error()})
}

func handleAuthError(c *gin.Context, err error) {
	switch {
	// Register
	case errors.Is(err, apperror.ErrAuthUserEmailExists):
		c.JSON(http.StatusBadRequest, apperror.ErrAuthUserEmailExists)
	case errors.Is(err, apperror.ErrAuthEmailFormat):
		c.JSON(http.StatusBadRequest, apperror.ErrAuthEmailFormat)
	case errors.Is(err, apperror.ErrAuthPassword):
		c.JSON(http.StatusBadRequest, err)
	// Login
	case errors.Is(err, apperror.ErrAuthUserNotFound):
		c.JSON(http.StatusBadRequest, apperror.ErrAuthUserNotFound)
	default:
		c.JSON(http.StatusBadRequest, gin.H{"code": apperror.BadRequest, "message": err.Error()})
	}
}
