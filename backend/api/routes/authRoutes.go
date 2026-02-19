package routes

import (
	"Server/controllers"
	"Server/validation"

	"github.com/gofiber/fiber/v3"
)

func SetupAuthRoutes(app *fiber.App) {
	// auth
	app.Post("/user/signup", validation.ValidateUser, controllers.Register)
	app.Post("/user/signin", validation.ValidateUser, controllers.Login)

}
