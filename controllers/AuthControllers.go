package controllers

import (
	"fmt"
	"strconv"

	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/gofiber/fiber/v2"
	"github.com/vangmay/cvwo-backend/database"
	"github.com/vangmay/cvwo-backend/models"
	"golang.org/x/crypto/bcrypt"
)

type EditRequest struct {
	UserName string `json:"username"`
	Email    string `json:"email"`
}

const SECRETKEY = "secret"

func Register(c *fiber.Ctx) error {
	var data map[string]string

	// {
	// 	"username" :
	// 	"email" :
	// 	"password"
	// } THIS SHOULD BE THE BODY INPUT

	err := c.BodyParser(&data)

	if err != nil {
		return err
	}
	password, _ := bcrypt.GenerateFromPassword([]byte(data["password"]), 14)

	//Check if a user with same email already exists

	var tempUser models.User

	database.DB.Where("email = ?", data["email"]).First(&tempUser)

	if tempUser.Email == data["email"] {
		c.Status(fiber.StatusBadRequest)
		return c.JSON(fiber.Map{
			"message": "A user has already registered using this email",
		})
	}

	database.DB.Where("user_name = ?", data["username"]).First(&tempUser)
	if tempUser.UserName == data["username"] {
		c.Status(fiber.StatusBadRequest)
		return c.JSON(fiber.Map{
			"message": "A user has already registered using this username",
		})
	}

	user := models.User{
		UserName: data["username"],
		Email:    data["email"],
		Password: password,

		IsAdmin: false,
	}

	database.DB.Create(&user)
	return c.JSON(user)
}

func Login(c *fiber.Ctx) error {
	// Body of login request
	// {
	// 	"email" :
	// 	"password" :
	// }

	var data map[string]string

	err := c.BodyParser(&data)

	if err != nil {
		return err
	}

	var user models.User

	database.DB.Where("email = ?", data["email"]).First(&user)

	if user.Id == 0 {
		c.Status(fiber.StatusNotFound)
		return c.JSON(fiber.Map{
			"message": "User not found",
		})
	}

	if err := bcrypt.CompareHashAndPassword(user.Password, []byte(data["password"])); err != nil {
		c.Status(fiber.StatusBadRequest)
		return c.JSON(fiber.Map{
			"message": "incorrect password",
		})
	}
	// If the compare and hash password function doesn't return any error, this means the test has been passed
	// We can proceed forward and assign a few claims

	claims := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.StandardClaims{
		Issuer:    strconv.Itoa(int(user.Id)),
		ExpiresAt: time.Now().Add(time.Hour * 24).Unix(),
	})

	token, err := claims.SignedString([]byte(SECRETKEY))

	if err != nil {
		c.Status(fiber.StatusInternalServerError)
		return c.JSON(fiber.Map{
			"message": "Could not login",
		})
	}

	// We generate a cookie for this user and authenticate him

	cookie := fiber.Cookie{
		Name:     "jwt",
		Value:    token,
		Expires:  time.Now().Add(time.Hour * 24),
		HTTPOnly: true,
	}

	c.Cookie(&cookie)

	return c.JSON(fiber.Map{
		"message": "great success",
	})

}
func User(c *fiber.Ctx) error {
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

	var user models.User

	database.DB.Where("id = ?", claims.Issuer).First(&user)

	return c.JSON(user)
}

func Logout(c *fiber.Ctx) error {
	// Destroy the cookie by replacing it with a similar cookie which has a value of empty string
	cookie := fiber.Cookie{
		Name:     "jwt",
		Value:    "",
		Expires:  time.Now().Add(-time.Hour),
		HTTPOnly: true,
	}
	c.Cookie(&cookie)

	return c.JSON(fiber.Map{
		"message": "success",
	})
}

func Edit(c *fiber.Ctx) error {

	// user can edit the following values
	// userName, bio, profilePhoto
	edits := EditRequest{}

	err := c.BodyParser(&edits)
	fmt.Println(err)

	if err != nil {
		return c.JSON(fiber.Map{
			"message": "Could not update user",
		})
	}

	user := models.User{}
	id := c.Params("id")

	database.DB.Where("id = ?", id).First(&user)

	user.UserName = edits.UserName
	user.Email = edits.Email
	database.DB.Save(&user)

	c.JSON(user)
	return nil
}

func Delete(c *fiber.Ctx) error {
	user := models.User{}
	err := c.BodyParser(&user)

	if err != nil {
		return c.JSON(fiber.Map{
			"message": "Error deleting user, please check if this user is registered",
		})
	}
	database.DB.Delete(&user)

	cookie := fiber.Cookie{
		Name:     "jwt",
		Value:    "",
		Expires:  time.Now().Add(-time.Hour),
		HTTPOnly: true,
	}
	c.Cookie(&cookie)

	return c.JSON(fiber.Map{
		"message": "User deleted successfully",
	})
}

func GetUserById(c *fiber.Ctx) error {

	id := c.Params("id")

	user := models.User{}

	err := database.DB.Where("id = ?", id).First(&user).Error
	if err != nil {
		return c.JSON(fiber.Map{
			"message": "Could not find user",
		})
	}
	c.JSON(user)
	return nil
}

func GetUsers(c *fiber.Ctx) error {
	users := &[]models.User{}

	err := database.DB.Find(users).Error

	if err != nil {
		return c.JSON(fiber.Map{
			"message": "could not get users",
		})
	}

	c.JSON(users)
	return nil
}
