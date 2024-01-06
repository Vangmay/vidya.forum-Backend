package controllers

import (
	"fmt"
	"strconv"
	"strings"

	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/gofiber/fiber/v2"
	"github.com/vangmay/cvwo-backend/database"
	"github.com/vangmay/cvwo-backend/models"
	"golang.org/x/crypto/bcrypt"
)

type EditRequest struct { // This struct contains the layout of an edit request
	UserName string `json:"username"`
	Email    string `json:"email"`
	Password string `json: "password"`
}

const SECRETKEY = "secret"

func Register(c *fiber.Ctx) error {
	/*
		Body of register request
		{
			"username" : "",
			"email"		: "",
			"password"	: ""
		}

		- Read the INPUT body
		- Perform uniquiness checks on the username and email
		- If check fails: send appropriate error message
		- If check passes : create a user struct and assign proper values from "data"
		- Use bcrypt to generate a hash
		- Send database command to create new user

		OUTPUT:
		- Return JSON containing the user that has been created for acknowledgement

	*/

	var data map[string]string

	err := c.BodyParser(&data)

	if err != nil {
		return err
	}
	password, _ := bcrypt.GenerateFromPassword([]byte(data["password"]), 14)

	var tempUser models.User

	//Check if a user with same email already exists
	database.DB.Where("email = ?", data["email"]).First(&tempUser)
	if tempUser.Email == data["email"] {
		c.Status(fiber.StatusBadRequest)
		return c.JSON(fiber.Map{
			"message": "A user has already registered using this email",
		})
	}

	//Check if a user with same username already exists
	database.DB.Where("user_name = ?", data["username"]).First(&tempUser)
	if tempUser.UserName == data["username"] {
		c.Status(fiber.StatusBadRequest)
		return c.JSON(fiber.Map{
			"message": "A user has already registered using this username",
		})
	}

	user := models.User{
		UserName: strings.ToLower(data["password"]),
		Email:    strings.ToLower(data["email"]),
		Password: password,
		IsAdmin:  false,
	}

	database.DB.Create(&user)
	return c.JSON(user)
}

func Login(c *fiber.Ctx) error {
	/*
		INPUT Body of login request
		{
			"email" :
			"password" :
		}
		- Read the body
		- Search for a user in the database with the same email
		- If user not found => send the relevant message.
		- If found => Proceed
		- Hash the password entered by the user and compare to the user we searched. (Done using bcrypt)
		- Check fails => Send message "Incorrect password"
		- Check passes => Proceed
		- Create a new claim which expires in next 24 hours
		- Use it to generate a token
		- Create a cookie with the token
		OUTPUT:
		- Send a message with acknowledement that login has been sucessfull
	*/

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
	// If the compare and hash password function doesn't return any error,
	// this means the test has been passed
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
		"message": "You have been logged in",
	})

}
func Profile(c *fiber.Ctx) error {
	/*
		No body, the function can be called with a GET request
		uses the cookie generated in the previous function to get current user
		If no user, then returns a message "unauthenticated user"
		Else gets the current user and sends a json file containing the user details
		OUTPUT: JSON file with user details

	*/
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
	/*
			Destroy the cookie created when logging in
		 	by replacing it with a similar cookie which has a value of empty string

			NO INPUT
			Output => You have been logged out
	*/
	cookie := fiber.Cookie{
		Name:     "jwt",
		Value:    "",
		Expires:  time.Now().Add(-time.Hour),
		HTTPOnly: true,
	}
	c.Cookie(&cookie)

	return c.JSON(fiber.Map{
		"message": "You have been logged out",
	})
}

func Edit(c *fiber.Ctx) error {
	/*
		Input Body of the PATCH request
		{
			"username" : "",
			"email"	   : "",
			"password" : "",
		}

		- Gets the Id of the current user by calling "GetCurrentUserId"
		- Changes the username, email and password
		- Calls the save method on DB to save the user.
		- Outputs the updated user in JSON format.

		Future security update:
		I can generate another random secret key when registering a user
		and ask the user to store it in a .txt file,
		this key can be used to verify the user to ensure the correct user is editing
		and not someone posing as the user.
	*/
	edits := EditRequest{}

	err := c.BodyParser(&edits)
	fmt.Println(err)

	userId := GetCurrentUserId(c)
	userProfile := models.User{}
	database.DB.Where("id = ?", userId).First(&userProfile)

	if err != nil {
		return c.JSON(fiber.Map{
			"message": "Could not update user",
		})
	}
	newPassword, _ := bcrypt.GenerateFromPassword([]byte(edits.Password), 14)
	userProfile.UserName = edits.UserName
	userProfile.Email = edits.Email
	userProfile.Password = newPassword
	database.DB.Save(&userProfile)

	c.JSON(userProfile)
	return nil
}

func Delete(c *fiber.Ctx) error {
	/*
		Internally gets the user who is currently logged in. (To ensure user's can only delete their own profile)
		Deletes the user
		Deletes the cookie, forcing the user to log out.

		Output => "User deleted succesfully"
	*/

	user := models.User{}
	userId := GetCurrentUserId(c)
	err := database.DB.Where("Id = ?", userId).First(&user).Error

	if err != nil {
		return c.JSON(fiber.Map{
			"message": "Error deleting user, please check if this user is registered or that you are logged in",
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
	/*
		INPUT = "id", passed as a URL parameter
		- Searches the database for a user with same id
		- If not found: Outputs a message "Could not find user"
		- If found: Outputs the user struct in JSON format
	*/
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
	/*
		- Gets all the users that are registered
		- Outputs a slice of user struct in JSON format
	*/
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
