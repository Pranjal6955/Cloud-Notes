package main

import (
	"log"
	"net/http"
	"os"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"collaborative-notes/config"
	"collaborative-notes/controllers"
	"collaborative-notes/middleware"
	"collaborative-notes/websocket"
)

func main() {
	// Load environment variables
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found")
	}

	// Initialize database
	config.InitDB()

	// Initialize WebSocket hub
	hub := websocket.NewHub()
	go hub.Run()

	// Setup Gin router
	r := gin.Default()

	// CORS middleware
	r.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"*"},
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"*"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
	}))

	// Auth routes
	auth := r.Group("/api/auth")
	{
		auth.POST("/signup", controllers.SignUp)
		auth.POST("/login", controllers.Login)
		auth.POST("/refresh", controllers.RefreshToken)
	}

	// Protected routes
	api := r.Group("/api")
	api.Use(middleware.AuthMiddleware())
	{
		// User routes
		api.GET("/profile", controllers.GetProfile)
		api.PUT("/profile", controllers.UpdateProfile)

		// Note routes
		api.GET("/notes", controllers.GetNotes)
		api.POST("/notes", controllers.CreateNote)
		api.GET("/notes/:id", controllers.GetNote)
		api.PUT("/notes/:id", controllers.UpdateNote)
		api.DELETE("/notes/:id", controllers.DeleteNote)
		api.POST("/notes/:id/share", controllers.ShareNote)
		api.GET("/notes/:id/collaborators", controllers.GetCollaborators)
	}

	// WebSocket endpoint
	r.GET("/ws/:noteId", func(c *gin.Context) {
		websocket.HandleWebSocket(hub, c.Writer, c.Request, c.Param("noteId"))
	})

	// Health check
	r.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"status": "healthy"})
	})

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	log.Printf("Server starting on port %s", port)
	log.Fatal(http.ListenAndServe(":"+port, r))
}
