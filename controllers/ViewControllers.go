package controllers

import (
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
