package controllers

import (
	"net/http"
	"strconv"

	"github.com/gofiber/fiber/v2"
	"github.com/vangmay/cvwo-backend/database"
	"github.com/vangmay/cvwo-backend/models"
)

func GetCommentByPost(c *fiber.Ctx) error {
	/*
		Receives a PostId as URL Parameter for INPUT
		Searches for the post
		Returns all the comments for that specific post
	*/
	PostId, _ := c.ParamsInt("PostId")

	post := &models.Post{Id: uint(PostId)}

	err := database.DB.Preload("User").Preload("Comments").Preload("Comments.User").Preload("Comments.Post").Find(&post).Error

	if err != nil {
		c.Status(http.StatusBadRequest).JSON(
			&fiber.Map{"Message ": "Could not get the post"})
		return err
	}
	c.Status(http.StatusOK).JSON(post.Comments)
	return nil
}

func GetCommentById(c *fiber.Ctx) error {
	/*
		Receives the commentId as URL parameter
		Outputs the comment in JSON format
	*/
	CommentId, _ := strconv.Atoi(c.Params("commentId"))

	comment := models.Comment{
		Id: uint(CommentId),
	}

	database.DB.Preload("Post").Preload("User").Find(&comment)

	return c.JSON(comment)
}

func CreateComment(c *fiber.Ctx) error {
	/*
		2 inputs
		JSON OBJECT { content: "" } and postId as url parameter

		- Gets the post and user
		- Creates a comment
		- Assigns UserId and PostId
		- Returns the comment
	*/
	postId, _ := strconv.Atoi(c.Params("postId"))

	post := models.Post{}
	database.DB.Where("id = ?", postId).First(&post)

	userId, _ := strconv.Atoi(GetCurrentUserId(c))

	currUser := models.User{}
	database.DB.Where("id = ?", userId).First(&currUser)

	comment := models.Comment{}

	err := c.BodyParser(&comment)
	if err != nil {
		return c.JSON(fiber.Map{
			"message": "Error creating comment",
		})
	}
	comment.UserId = uint(userId)
	comment.PostId = uint(postId)

	database.DB.Create(&comment)
	database.DB.Preload("Post").Preload("User").Preload("Post.User").Find(&comment)
	return c.JSON(comment)
}

func EditComment(c *fiber.Ctx) error {
	/*
		Receives a json object {"content" : ""}
		Also recevies commentId as URL parameter
	*/
	commentId, _ := c.ParamsInt("commentId")

	comment := models.Comment{
		Id: uint(commentId),
	}
	database.DB.Find(&comment)

	updatedComment := models.Comment{
		Id: uint(commentId),
	}
	err := c.BodyParser(&updatedComment)
	if err != nil {
		return c.JSON(fiber.Map{
			"message": "Error parsing the comment",
		})
	}
	database.DB.Model(&updatedComment).Updates(&updatedComment)

	database.DB.Save(&updatedComment)
	return c.JSON(updatedComment)
}

func DeleteComment(c *fiber.Ctx) error {

	/*
		INPUT: UserId as a URL Parameter
		Deletes the comment are returns a message
	*/
	currUserId, _ := strconv.Atoi(GetCurrentUserId(c))
	commentId, _ := c.ParamsInt("commentId")
	comment := models.Comment{}
	database.DB.Where("Id = ?", commentId).First(&comment)
	if uint(currUserId) != comment.UserId {
		return c.JSON(fiber.Map{
			"message": "Unauthorized user",
		})
	}
	database.DB.Delete(comment)
	return c.JSON(fiber.Map{
		"message": "Delete the comment",
	})
}
