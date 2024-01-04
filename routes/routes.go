package routes

import (
	"github.com/gofiber/fiber/v2"
	"github.com/vangmay/cvwo-backend/controllers"
)

func Setup(app *fiber.App) {

	// Authorization and User management Endpoints
	app.Get("/Auth/user/:id", controllers.GetUserById)
	app.Get("/Auth/users", controllers.GetUsers)
	app.Get("/Auth/profile", controllers.User)
	app.Post("/Auth/register", controllers.Register)
	app.Post("/Auth/login", controllers.Login)
	app.Post("/Auth/logout", controllers.Logout)
	app.Patch("/Auth/Edit/:id", controllers.Edit)
	app.Delete("/Auth/delete", controllers.Delete)

	// Post Endpoints
	app.Get("/posts", controllers.GetAllPosts)
	app.Get("/posts/popular", controllers.GetPopularPosts)
	app.Get("/posts/:tag", controllers.GetPostsByTag)
	app.Get("/post/:id", controllers.GetPostById)
	app.Post("/post/create", controllers.CreatePost)
	app.Patch("/post/edit/:id", controllers.EditPost)
	app.Delete("/post/:id", controllers.DeletePost)
	app.Post("/post/like/:id", controllers.LikePost)
	app.Post("/post/unlike/:id", controllers.UnlikePost)

	// Comment Endpoints
	app.Get("/comments/:PostId", controllers.GetCommentByPost)
	app.Get("/comments/id/:commentId", controllers.GetCommentById)
	app.Post("/comment/create/:postId", controllers.CreateComment)
	app.Patch("/comment/edit/:commentId", controllers.EditComment)
	app.Delete("/comment/delete/:commentId", controllers.DeleteComment)

}
