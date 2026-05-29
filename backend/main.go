package main

import (
	"log"
	"securetask/database"
	"securetask/handlers"
	"securetask/models"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"

	"github.com/joho/godotenv"
)

func main() {
	// Loads .env to program
	if err := godotenv.Load("../.env"); err != nil {
		log.Println("Warning: Could not load .env file", err)
	}

	// Initialize database
	database.Connect()

	// Auto-migrate models
	err := database.DB.AutoMigrate(&models.User{}, &models.Task{})
	if err != nil {
		log.Fatal("Failed to migrate database:", err)
	}

	// Seed initial data
	seedData()

	// Setup Gin router
	r := gin.Default()

	// VULNERABILITY: Permissive CORS - allows all origins
	r.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"http://localhost:5173"},
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Authorization"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
	}))

	// Public routes (no authentication required)
	r.POST("/api/auth/register", handlers.Register)
	r.POST("/api/auth/login", handlers.Login)

	// VULNERABILITY #2: No authentication middleware on these routes!

	// Protected routes (with auth middleware)
	authorized := r.Group("/api")
	authorized.Use(handlers.AuthMiddleware())
	{
		authorized.GET("/tasks", handlers.GetTasks)
		authorized.POST("/tasks", handlers.CreateTask)
		authorized.PUT("/tasks/:id", handlers.UpdateTask)
		authorized.GET("/tasks/search", handlers.SearchTasks) // Should require auth
		authorized.DELETE("/tasks/:id", handlers.DeleteTask)  // Should require auth
		authorized.GET("/users/me", handlers.GetCurrentUser)
		authorized.PUT("/users/:id/profile", handlers.UpdateProfile)
		authorized.GET("/admin/users", handlers.GetAllUsers)
	}

	log.Println("🚀 Server starting on port 8080...")
	log.Println("⚠️  WARNING: This server contains intentional security vulnerabilities!")
	r.Run(":8080")
}

func seedData() {
	// Check if users already exist
	var count int64
	database.DB.Model(&models.User{}).Count(&count)
	if count > 0 {
		return // Data already seeded
	}

	// VULNERABILITY #5: Passwords stored in plain text (no hashing!)
	users := []models.User{
		{
			Email:    "admin@example.com",
			Password: hashPasswordOrPanic("admin123"), // Plain text password!
			Name:     "Admin User",
			Role:     "admin",
			Bio:      "I'm the administrator",
		},
		{
			Email:    "user@example.com",
			Password: hashPasswordOrPanic("password123"), // Plain text password!
			Name:     "Regular User",
			Role:     "user",
			Bio:      "Just a regular user",
		},
	}

	for _, user := range users {
		database.DB.Create(&user)
	}

	log.Println("✅ Database seeded with initial users")
}

func hashPasswordOrPanic(password string) string {
	hashed, err := handlers.HashPassword(password)
	if err != nil {
		panic(err)
	}
	return hashed
}
