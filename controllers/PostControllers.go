package controllers

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/dgrijalva/jwt-go"
	"github.com/gofiber/fiber/v2"
	"github.com/vangmay/cvwo-backend/database"
	"github.com/vangmay/cvwo-backend/models"
)

func GetCurrentUserId(c *fiber.Ctx) string {
	/*
		This is a helper function meant for internal use
		Cannot be accessed via endpoint
	*/
	cookie := c.Cookies("jwt")
	token, _ := jwt.ParseWithClaims(cookie, &jwt.StandardClaims{}, func(token *jwt.Token) (interface{}, error) {
		return []byte(SECRETKEY), nil
	})
	claims := token.Claims.(*jwt.StandardClaims)
	userId := claims.Issuer
	return userId
}

func GetAllPosts(c *fiber.Ctx) error {
	/*
		- Preloads the User and Comments field in every post
		- Outputs an array containing every post
	*/
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
	/*
		INPUT - Id, Passed as a URL parameter
		- Searches the database for a post with same Id
		- Preloads the user and comments
		- Outputs the post in JSON format
	*/
	id, _ := c.ParamsInt("id")
	post := &models.Post{
		Id: uint(id),
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
	/*
		INPUT json object
		{
			"title": "",
			"body" : "",
			"tag" : ""
		}
		- Assign the object to newPost
		- Get current user id
		- Set it was newPost.UserId
		- Preload users to populate User
		- OUTPUT the new post as JSON object
	*/
	userId, _ := strconv.Atoi(GetCurrentUserId(c))

	newPost := models.Post{}

	c.BodyParser(&newPost)

	newPost.UserId = uint(userId)
	newPost.Likes = 0
	newPost.IsEdited = false

	database.DB.Create(&newPost)
	database.DB.Preload("User").First(&newPost)

	return c.JSON(newPost)
}

func EditPost(c *fiber.Ctx) error {
	/*
		INPUT JSON OBJECT, id is received as URL parameter
		{
			Title : "",
			Body : "",
			Tag : ""
		}
		Fetches the current user id
		Gets the post based on it
		Checks if the post.UserId == CurrentUserId
		If no, send unauthenticated
		If yes, edit the post
		OUTPUT: Returns the edited post in JSON format
	*/
	type EditRequest struct {
		Title string `json:"title"`
		Body  string `json:"body"`
		Tag   string `json: "tag"`
	}

	//Get current user id
	userId, _ := strconv.Atoi(GetCurrentUserId(c))

	postId, err := c.ParamsInt("id")

	postToEdit := models.Post{}
	database.DB.Where("id = ?", postId).First(&postToEdit)

	edits := EditRequest{}
	c.BodyParser(&edits)

	if uint(userId) != postToEdit.UserId {
		c.Status(fiber.StatusForbidden)
		return c.JSON(fiber.Map{
			"message": "Unauthorized",
		})
	}

	postToEdit.Title = edits.Title
	postToEdit.Body = edits.Body
	postToEdit.Tag = edits.Tag
	postToEdit.IsEdited = true

	if err != nil {
		return c.JSON(fiber.Map{
			"message": "error creating post",
		})
	}

	database.DB.Save(&postToEdit)
	database.DB.Preload("User").Preload("Comments").Preload("Comments.User").Find(&postToEdit)
	return c.JSON(postToEdit)
}

func DeletePost(c *fiber.Ctx) error {
	// Get current user id
	// Check if the userId matched the userId in post
	// If Not: Throw unauthenticated error
	// Else delete the post
	// Send acknowledgement that post has been deleted

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
	fmt.Println(currentUser)

	post := models.Post{}
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
		return nil
	}
	err = database.DB.Delete(&models.Post{}, PostId).Error

	if err != nil {
		c.Status(http.StatusBadRequest).JSON(&fiber.Map{
			"message": "could not delete book",
		})
		return err

	}
	c.Status(http.StatusOK).JSON(fiber.Map{
		"message": "Post has been deleted",
	})
	return nil

}
