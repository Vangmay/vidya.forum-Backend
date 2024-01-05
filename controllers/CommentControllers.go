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
	PostId, _ := c.ParamsInt("PostId")

	post := &models.Post{Id: uint(PostId)}

	err := database.DB.Preload("User").Preload("Comments").Preload("Comments.User").Find(&post).Error

	if err != nil {
		c.Status(http.StatusBadRequest).JSON(
			&fiber.Map{"Message ": "Could not get the post"})
		return err
	}
	c.Status(http.StatusOK).JSON(post.Comments)
	return nil
}

func GetCommentById(c *fiber.Ctx) error {
	CommentId, _ := strconv.Atoi(c.Params("commentId"))

	comment := models.Comment{
		Id: uint(CommentId),
	}

	database.DB.Preload("Posts").Preload("User").Find(&comment)

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
	comment.UserId = uint(userId)
	comment.PostId = uint(postId)

	database.DB.Create(&comment)
	return c.JSON(comment)
}

func EditComment(c *fiber.Ctx) error {
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
	return nil
}
