package controllers

import (
	"net/http"
	"strconv"

	"github.com/gofiber/fiber/v2"
	"github.com/vangmay/cvwo-backend/database"
	"github.com/vangmay/cvwo-backend/models"
)

//GetCommentsByPost
//GetCommentsByUser

//Create Comment

//Edit comment

//Delete comment

func GetCommentByPost(c *fiber.Ctx) error {
	id := c.Params("id") // Post id

	post := &models.Post{}
	err := database.DB.Find("id = ?", id).First(post).Error

	if err != nil {
		c.Status(http.StatusBadRequest).JSON(
			&fiber.Map{"message": "Could not get the comment"})
		return err
	}

	c.Status(http.StatusOK).JSON(post.Comments)
	return nil
}

func GetCommentById(c *fiber.Ctx) error {
	CommentId := c.Params("commentId")

	comment := models.Comment{}
	err := database.DB.Find("id = ?", CommentId).First(&comment).Error
	if err != nil {
		c.Status(http.StatusBadRequest).JSON(
			&fiber.Map{"message": "could not get comment"},
		)
		return err
	}

	return c.JSON(comment)
}

func CreateComment(c *fiber.Ctx) error {
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
	post.Comments = append(post.Comments, comment)
	comment.UserId = uint(userId)
	comment.PostId = uint(postId)
	comment.Post = post
	comment.User = currUser

	database.DB.Create(&comment)
	database.DB.Save(&post)
	return c.JSON(comment)
}

func EditComment(c *fiber.Ctx) error {
	commentId := c.Params("commentId")

	comment := models.Comment{}
	database.DB.Where("id = ?", commentId).First(&comment)

	updatedComment := models.Comment{}
	err := c.BodyParser(&updatedComment)
	if err != nil {
		return c.JSON(fiber.Map{
			"message": "Error parsing the comment",
		})
	}
	comment.Content = updatedComment.Content

	database.DB.Save(&comment)
	return c.JSON(comment)
}

func DeleteComment(c *fiber.Ctx) error {
	currUserId, _ := strconv.Atoi(GetCurrentUserId(c))
	comment := models.Comment{}
	c.BodyParser(&comment)
	if uint(currUserId) != comment.UserId {
		return c.JSON(fiber.Map{
			"message": "Unauthorized user",
		})
	}
	database.DB.Delete(comment)
	return nil
}
