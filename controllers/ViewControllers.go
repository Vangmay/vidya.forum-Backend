package controllers

import (
	"fmt"
	"net/http"
	"sort"

	"github.com/gofiber/fiber/v2"
	"github.com/vangmay/cvwo-backend/database"
	"github.com/vangmay/cvwo-backend/models"
)

/*
- Get Post based on likes
- Get post based on a certain tag (Most liked first)
Both of these functions first get a list of all the posts
Then sort them accordingly
*/
func GetPopularPosts(c *fiber.Ctx) error {
	posts := &[]models.Post{}

	err := database.DB.Preload("User").Preload("Comments").Preload("Comments.User").Find(posts).Error
	if err != nil {
		c.Status(http.StatusBadRequest).JSON(
			&fiber.Map{"message": "could not get the books"})
		return err
	}
	sortByLikes(*posts)

	c.Status(http.StatusOK).JSON(posts)
	return nil
}

func GetPostsByTag(c *fiber.Ctx) error {
	requestedTag := c.Params("tag")

	posts := &[]models.Post{}
	err := database.DB.Preload("User").Preload("Comments").Preload("Comments.User").Where("tag = ?", requestedTag).Find(posts).Error
	if err != nil {
		c.Status(http.StatusBadRequest).JSON(
			&fiber.Map{"message": "Could not get the books"})
		return err
	}
	sortByLikes(*posts)
	c.Status(http.StatusOK).JSON(posts)
	return nil
}

func GetPostsByUser(c *fiber.Ctx) error {
	userID := c.Params("id")

	// Convert userID to uint
	var uid uint
	if _, err := fmt.Sscan(userID, &uid); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid user ID",
		})
	}

	// Fetch posts by user ID from the database
	var posts []models.Post
	if err := database.DB.Where("user_id = ?", uid).Find(&posts).Error; err != nil {
		if err == fiber.ErrNotFound {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
				"error": "User not found or has no posts",
			})
		}
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to fetch posts",
		})
	}

	// Return the posts in the response
	return c.JSON(posts)
}

func GetLikesByUserID(c *fiber.Ctx) error {
	// Parse user ID from request parameters
	userID := c.Params("id")

	// Convert userID to uint
	var uid uint
	if _, err := fmt.Sscan(userID, &uid); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid user ID",
		})
	}

	// Fetch the sum of likes by user ID from the database
	var likesSum int
	if err := database.DB.Model(&models.Post{}).Where("user_id = ?", uid).Select("SUM(likes) as likes_sum").Scan(&likesSum).Error; err != nil {
		if err == fiber.ErrNotFound {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
				"error": "User not found or has no posts with likes",
			})
		}
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to fetch likes",
		})
	}

	// Return the sum of likes in the response
	return c.JSON(fiber.Map{"likes_sum": likesSum})
}

func sortByLikes(data []models.Post) {
	sorter := &postSorter{
		data: data,
	}
	sort.Sort(sorter)
}

type postSorter struct {
	data []models.Post
}

func (s *postSorter) Len() int {
	return len(s.data)
}
func (s *postSorter) Swap(i, j int) {
	s.data[i], s.data[j] = s.data[j], s.data[i]
}
func (s *postSorter) Less(i, j int) bool {
	return s.data[i].Likes > s.data[j].Likes
}
