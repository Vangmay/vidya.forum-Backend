package routes

import (
	"github.com/gofiber/fiber/v2"
	"github.com/vangmay/cvwo-backend/controllers"
)

func Setup(app *fiber.App) {

	// Authorization and Accounts Endpoints
	app.Get("/Auth/users/:id", controllers.GetUserById)
	app.Get("/Auth/users", controllers.GetUsers)
	app.Get("/Auth/user", controllers.User)
	app.Post("/Auth/register", controllers.Register)
	app.Post("/Auth/login", controllers.Login)
	app.Post("/Auth/logout", controllers.Logout)
	app.Delete("/Auth/delete", controllers.Delete)
	app.Patch("/Auth/Edit/:id", controllers.Edit)

	// Post and Comment Endpoints
	// GetPosts
	// CreatePosts
	// EditPosts [OP]
	// DeletePosts [OP]
	app.Get("/posts", controllers.GetAllPosts)
	app.Get("/post/:id", controllers.GetPostById)
	app.Post("/posts/create", controllers.CreatePost)
	app.Patch("/posts/edit/:id", controllers.EditPost)
	app.Delete("/posts/:id", controllers.DeletePost)

	app.Get("/comment/:PostId", controllers.GetCommentByPost)
	app.Get("/comment/id/:commentId", controllers.GetCommentById)
	app.Post("/comment/create/:postId", controllers.CreateComment)
	app.Patch("/comment/edit/:commentId", controllers.EditComment)
	app.Delete("/comment/delete/:commentId", controllers.DeleteComment)

}
