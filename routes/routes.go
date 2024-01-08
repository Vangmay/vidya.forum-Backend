package routes

import (
	"github.com/gofiber/fiber/v2"
	"github.com/vangmay/cvwo-backend/controllers"
)

func Setup(app *fiber.App) {

	// Authorization and User management Endpoints (Everything Works in AUTH)
	app.Get("/Auth/user/:id", controllers.GetUserById) // Gets user with the same UserId as parameter : "id"
	app.Get("/Auth/users", controllers.GetUsers)       // Gets all the users
	app.Get("/Auth/profile", controllers.Profile)      // Gets the struct containing the current logged in user
	app.Post("/Auth/register", controllers.Register)   // Registers a user
	app.Post("/Auth/login", controllers.Login)         // Logs a user in
	app.Post("/Auth/logout", controllers.Logout)       // Logout
	app.Patch("/Auth/profile", controllers.Edit)       // Edits the username, email and password
	app.Delete("/Auth/profile", controllers.Delete)    // Deletes the profile of the current logged in user

	// Post Endpoints (EVERYTHING WORKS AS EXPECTED IN POSTS)
	app.Get("/posts", controllers.GetAllPosts)             // Gets all the posts
	app.Get("/posts/popular", controllers.GetPopularPosts) // Gets all the posts (Sorted by likes)
	app.Get("/posts/:tag", controllers.GetPostsByTag)      // Gets all the posts based on a tag (Sorted by likes)
	app.Get("/post/:id", controllers.GetPostById)          // Gets a post based on Id parameter
	app.Post("/post", controllers.CreatePost)              // Creates a new post
	app.Patch("/post/:id", controllers.EditPost)           // Edits an existing post
	app.Delete("/post/:id", controllers.DeletePost)        // Deletes a post
	app.Post("/post/like/:id", controllers.LikePost)       // Like a post
	app.Post("/post/unlike/:id", controllers.UnlikePost)   // Unlike a post

	// Comment Endpoints
	app.Get("/comments/:PostId", controllers.GetCommentByPost)   // Gets all comments present under post with PostId : id
	app.Get("/comment/:commentId", controllers.GetCommentById)   // Gets comment based on commentId
	app.Post("/comment/:postId", controllers.CreateComment)      // Creates a new comment
	app.Patch("/comment/:commentId", controllers.EditComment)    // Edits a comment
	app.Delete("/comment/:commentId", controllers.DeleteComment) // Deletes a comments

}
