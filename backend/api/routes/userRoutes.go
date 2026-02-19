package routes

import (
	"Server/controllers"
	"Server/middleware"

	"github.com/gofiber/fiber/v3"
)

func SetupUserRoutes(app *fiber.App) {
	// auth
	app.Get("/user/getUser/:id", controllers.GetUserById)
	// get slug
	// update
	app.Patch("/user/Update/:id", middleware.AuthMiddleware, controllers.UpdateUser)
	// following
	// delete
}
