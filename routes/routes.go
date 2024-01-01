package routes

import (
	"github.com/gofiber/fiber/v2"
	"github.com/vangmay/cvwo-backend/controllers"
)

func Setup(app *fiber.App) {

	// Authorization and Accounts Endpoints
	app.Get("/Auth/user/:id", controllers.GetUserById)
	app.Get("/Auth/users", controllers.GetUsers)
	app.Post("/Auth/register", controllers.Register)
	app.Post("/Auth/login", controllers.Login)
	app.Get("/Auth/user", controllers.User)
	app.Post("/Auth/logout", controllers.Logout)
	app.Delete("/Auth/delete", controllers.Delete)
	app.Patch("/Auth/Edit/:id", controllers.Edit)

	// Post and Comment Endpoints

	// Tags endpoints
}
