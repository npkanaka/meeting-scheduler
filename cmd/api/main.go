package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/npkanaka/meeting-scheduler/internal/config"
	"github.com/npkanaka/meeting-scheduler/internal/handlers"
	"github.com/npkanaka/meeting-scheduler/internal/middleware" // Import the middleware package
	"github.com/npkanaka/meeting-scheduler/internal/repository"
	"github.com/npkanaka/meeting-scheduler/internal/service"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func main() {
	// Load configuration
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("Failed to load configuration: %v", err)
	}

	// Connect to database
	db, err := connectDB(cfg.Database.DSN)
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}

	// Initialize repositories
	eventRepo := repository.NewGormEventRepository(db)
	timeslotRepo := repository.NewGormTimeSlotRepository(db)
	userRepo := repository.NewGormUserRepository(db)
	availabilityRepo := repository.NewGormAvailabilityRepository(db)

	// Initialize services
	eventService := service.NewEventService(eventRepo)
	timeslotService := service.NewTimeSlotService(timeslotRepo, eventRepo)
	availabilityService := service.NewAvailabilityService(availabilityRepo, eventRepo, userRepo)
	recommendationService := service.NewRecommendationService(eventRepo, timeslotRepo, availabilityRepo, userRepo)

	// Initialize handlers
	eventHandler := handlers.NewEventHandler(eventService)
	timeslotHandler := handlers.NewTimeSlotHandler(timeslotService)
	availabilityHandler := handlers.NewAvailabilityHandler(availabilityService)
	recommendationHandler := handlers.NewRecommendationHandler(recommendationService)
	healthHandler := handlers.NewHealthHandler(db)

	// Create and configure Gin router
	router := gin.Default()

	// Apply middleware
	router.Use(gin.Logger())
	router.Use(gin.Recovery())
	router.Use(middleware.RequestLogger()) // Use your middleware
	router.Use(middleware.CORS())          // Use your CORS middleware

	// Register routes
	// Health check route
	router.GET("/health", healthHandler.Check)

	// Event routes
	router.POST("/events", eventHandler.Create)
	router.GET("/events", eventHandler.List)
	router.GET("/events/:id", eventHandler.Get)
	router.PUT("/events/:id", eventHandler.Update)
	router.DELETE("/events/:id", eventHandler.Delete)

	// Time slot routes - using :id consistently instead of :eventId
	router.POST("/events/:id/timeslots", timeslotHandler.Create)
	router.GET("/events/:id/timeslots", timeslotHandler.List)
	router.PUT("/timeslots/:id", timeslotHandler.Update)
	router.DELETE("/timeslots/:id", timeslotHandler.Delete)

	// Availability routes - using :id consistently instead of :eventId
	router.POST("/events/:id/availability", availabilityHandler.Create)
	router.GET("/events/:id/availability", availabilityHandler.GetEventAvailability)
	router.GET("/events/:id/availability/:userId", availabilityHandler.GetUserAvailability)
	router.PUT("/events/:id/availability/:userId", availabilityHandler.Update)
	router.DELETE("/availability/:id", availabilityHandler.Delete)

	// Recommendation routes - using :id consistently instead of :eventId
	router.GET("/events/:id/recommendations", recommendationHandler.GetRecommendations)

	// Start server
	srv := &http.Server{
		Addr:    ":" + cfg.Server.Port,
		Handler: router,
	}

	// Start the server in a goroutine
	go func() {
		log.Printf("Server listening on port %s", cfg.Server.Port)
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Failed to start server: %v", err)
		}
	}()

	// Graceful shutdown
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	log.Println("Shutting down server...")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := srv.Shutdown(ctx); err != nil {
		log.Fatalf("Server forced to shutdown: %v", err)
	}

	log.Println("Server exited properly")
}

func connectDB(dsn string) (*gorm.DB, error) {
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		return nil, err
	}
	return db, nil
}
