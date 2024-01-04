package controllers

import (
	"github.com/gofiber/fiber/v2"
	"github.com/vangmay/cvwo-backend/database"
	"github.com/vangmay/cvwo-backend/models"
)

// Like controller
// Get post based on tags, sorted by likes
// Get post based on likes

func LikePost(c *fiber.Ctx) error {
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
	PostId, _ := c.ParamsInt("id")

	post := models.Post{
		Id: uint(PostId),
	}

	database.DB.Where("id = ?", PostId).First(&post)

	post.Likes = post.Likes - 1
	database.DB.Save(post)

	return nil
}
