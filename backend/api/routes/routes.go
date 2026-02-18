package routes

import (
	"Server/controllers"
	"Server/validation"

	"github.com/gofiber/fiber/v3"
)

func SetupRoutes(app *fiber.App) {
	// auth
	app.Post("/user/signup", validation.ValidateUser, controllers.Register)
	app.Post("/user/signin", validation.ValidateUser, controllers.Login)

}
