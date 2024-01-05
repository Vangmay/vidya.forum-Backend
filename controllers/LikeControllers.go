package controllers

import (
	"github.com/gofiber/fiber/v2"
	"github.com/vangmay/cvwo-backend/database"
	"github.com/vangmay/cvwo-backend/models"
)

/*
	Note: Like/Dislike will be handled at the frontend level
	In the frontEnd
	- If the user has already liked the post, they will be shown a dislike button
	- And Like button if it is the opposite
*/

// Like controller
// Get post based on tags, sorted by likes
// Get post based on likes

func LikePost(c *fiber.Ctx) error {

	/*
		Received INPUT id as a URL parameter
		Increments the like count
	*/
	PostId, _ := c.ParamsInt("id")

	post := models.Post{
		Id: uint(PostId),
	}

	database.DB.Where("id = ?", PostId).First(&post)

	post.Likes = post.Likes + 1
	database.DB.Save(post)
	return nil
}
func UnlikePost(c *fiber.Ctx) error {
	/*
		Received INPUT id as a URL parameter
		Decrements the like count
	*/
	PostId, _ := c.ParamsInt("id")

	post := models.Post{
		Id: uint(PostId),
	}

	database.DB.Where("id = ?", PostId).First(&post)

	post.Likes = post.Likes - 1
	database.DB.Save(post)

	return nil
}
