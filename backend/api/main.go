package main

import (
	"Server/database"
	_ "Server/docs"
	"Server/routes"
	"log"

	swaggo "github.com/gofiber/contrib/v3/swaggo"
	"github.com/gofiber/fiber/v3"
	"github.com/gofiber/fiber/v3/middleware/cors"
	"github.com/joho/godotenv"
)

// @title Fiber Golang Mongo GRPC WEBSOCKET etc...
// @version 1.0
// @description This is Swagger docs for rest api golang fiber
// @host localhost:5000
// @BasePath
// @Schemes http
// @SecurityDefinitions.apiKey BearerAuth
// @In header
// @name Authorization
// @description Type "Bearer" followed by a space and the token
func main() {
	// load .env file
	if err := godotenv.Load(); err != nil {
		log.Fatal("Error loading .env file")
	}

	database.Connect()
	app := fiber.New()

	app.Use(cors.New(cors.Config{
		AllowCredentials: true,
		AllowOriginsFunc: func(origin string) bool {
			return true
		},
	}))

	app.Get("/", func(c fiber.Ctx) error {
		return c.SendString("Hello, World!")
	})

	// Setup routes
	routes.SetupAuthRoutes(app)
	routes.SetupUserRoutes(app)

	// Mount the UI with the default configuration under /swagger
	app.Get("/swagger/*", swaggo.HandlerDefault)

	app.Listen(":8080")
}
