package controllers

import (
	"net/http"

	"github.com/dgrijalva/jwt-go"
	"github.com/gofiber/fiber/v2"
	"github.com/vangmay/cvwo-backend/database"
	"github.com/vangmay/cvwo-backend/models"
)

func GetCurrentUserId(c *fiber.Ctx) string {
	cookie := c.Cookies("jwt")
	token, _ := jwt.ParseWithClaims(cookie, &jwt.StandardClaims{}, func(token *jwt.Token) (interface{}, error) {
		return []byte(SECRETKEY), nil
	})
	claims := token.Claims.(*jwt.StandardClaims)
	userId := claims.Issuer
	return userId
}

func GetAllPosts(c *fiber.Ctx) error {
	posts := &[]models.Post{}

	err := database.DB.Preload("User").Preload("Comments").Preload("Comments.User").Find(posts).Error
	if err != nil {
		c.Status(http.StatusBadRequest).JSON(
			&fiber.Map{"message": "could not get the books"})
		return err
	}

	c.Status(http.StatusOK).JSON(posts)
	return nil
}

func GetPostById(c *fiber.Ctx) error {
	id := c.Params("id")
	post := &models.Post{}
	if id == "" {
		c.Status(http.StatusInternalServerError).JSON(&fiber.Map{
			"message": "id cannot be empty",
		})
		return nil
	}

	err := database.DB.Preload("User").Preload("Comments").Preload("Comments.User").Find(&post).Error

	if err != nil {
		c.Status(http.StatusBadRequest).JSON(
			&fiber.Map{"Message ": "Could not get the post"})
		return err
	}
	c.Status(http.StatusOK).JSON(post)
	return nil
}

func CreatePost(c *fiber.Ctx) error {
	database.DB.Debug()
	userId := GetCurrentUserId(c)
	author := models.User{}
	err := database.DB.Where("id = ?", userId).First(&author).Error
	if err != nil {
		return c.JSON(fiber.Map{
			"message": "Cannot fetch post",
		})
	}

	newPost := models.Post{}

	c.BodyParser(&newPost)

	newPost.UserId = author.Id

	newPost.User = author

	database.DB.Preload("users").Create(&newPost)

	return c.JSON(newPost)
}

func EditPost(c *fiber.Ctx) error {
	type EditRequest struct {
		Title   string `json:"title"`
		Content string `json:"content"`
		Tag     string `json: "tag"`
	}
	// Only title content and tag can
	// be edited after a post has been created

	//Get user id
	userId := GetCurrentUserId(c)
	//

	author := models.User{}
	database.DB.Where("id = ?", userId).First(&author)

	postId, err := c.ParamsInt("id")
	newPost := models.Post{}
	database.DB.Where("id = ?", postId).First(&newPost)

	edits := EditRequest{}
	c.BodyParser(&edits)

	newPost.Title = edits.Title
	newPost.Content = edits.Content
	newPost.Tag = edits.Tag
	newPost.User = author

	if err != nil {
		return c.JSON(fiber.Map{
			"message": "error creating post",
		})
	}

	database.DB.Save(&newPost)
	database.DB.Preload("User").Preload("Comments").Preload("Comments.User").Find(&newPost)
	return c.JSON(newPost)
}

func DeletePost(c *fiber.Ctx) error {
	// Get curret user id
	// Check if the userId matched the userId in post
	// If Not: Throw unauthenticated error
	// Else delete the post

	cookie := c.Cookies("jwt")
	token, err := jwt.ParseWithClaims(cookie, &jwt.StandardClaims{}, func(token *jwt.Token) (interface{}, error) {
		return []byte(SECRETKEY), nil
	})
	if err != nil {
		c.Status(fiber.StatusUnauthorized)
		return c.JSON(fiber.Map{
			"message": "unauthenticated user",
			"Name":    "",
		})
	}
	claims := token.Claims.(*jwt.StandardClaims)
	id := claims.Issuer

	currentUser := models.User{}
	database.DB.Where("id = ?", id).First(&currentUser)

	post := models.Post{}
	posts := &[]models.Post{}
	PostId := c.Params("id")
	database.DB.Where("id = ?", PostId).First(&post)

	if PostId == "" {
		c.Status(http.StatusInternalServerError).JSON(&fiber.Map{
			"message": "id cannot be empty",
		})
		return nil
	}

	if post.UserId != currentUser.Id {
		c.Status(http.StatusUnauthorized).JSON(&fiber.Map{
			"message": "Unauthorized to delete the post",
		})
	}
	err = database.DB.Delete(&models.Post{}, PostId).Error

	if err != nil {
		c.Status(http.StatusBadRequest).JSON(&fiber.Map{
			"message": "could not delete book",
		})
		return err

	}
	c.Status(http.StatusOK).JSON(posts)
	return nil

}
