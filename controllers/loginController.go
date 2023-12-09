package controllers

import (
	"boilerplate/backend/database"
	"boilerplate/backend/helpers/jwt"
	"boilerplate/backend/helpers/verification"
	"boilerplate/backend/models"
	"os"
	"time"

	"github.com/VDiPaola/go-backend-module/module_helpers/module_encryption"
	"github.com/VDiPaola/go-backend-module/module_helpers/module_jwt"
	"github.com/VDiPaola/go-backend-module/module_models"
	"github.com/gofiber/fiber/v2"
)

func SignUp(c *fiber.Ctx) error {
	//gets user object from body
	var user models.User

	if err := c.BodyParser(&user); err != nil {
		return c.Status(fiber.StatusBadRequest).SendString(err.Error())
	}

	//get hashed password
	encryptedPass, err := module_encryption.EncryptPassword(user.Password)

	if err != nil {
		return err
	}

	user.Password = encryptedPass
	user.HasPassword = true

	//save user to database
	if dbc := database.Connection.Create(&user); dbc.Error != nil {
		c.Status(fiber.StatusBadRequest)
		return c.JSON(fiber.Map{
			"message": dbc.Error,
		})
	}

	//generate and send
	if err := verification.GenerateCodeAndSend(user, 15); err != nil {
		c.Status(fiber.StatusInternalServerError)
		return c.JSON(fiber.Map{
			"message": err.Error(),
		})
	}

	return c.Status(200).JSON(&user)
}

func GoogleLogin(c *fiber.Ctx) error {
	// Extract the ID token from the request body
	var token module_models.GoogleResponse
	if err := c.BodyParser(&token); err != nil {
		return c.Status(fiber.StatusBadRequest).SendString(err.Error())
	}

	// Verify the ID token
	tokenInfo, err := module_jwt.VerifyIDToken(token.JWT, os.Getenv("GOOGLE_CLIENT_ID"))
	if err != nil {
		return c.Status(fiber.StatusUnauthorized).SendString("Invalid ID token: " + err.Error())
	}

	//get user from database if exists
	var user models.User
	dbc := database.Connection.Where("email = ?", tokenInfo.Email).First(&user)

	//if not exist create new account
	if dbc.RowsAffected <= 0 {
		user = models.User{
			Email:       tokenInfo.Email,
			HasPassword: false,
			Role:        models.Role.Member,
		}
		if dbc := database.Connection.Create(&user); dbc.Error != nil {
			c.Status(fiber.StatusBadRequest)
			return c.JSON(fiber.Map{
				"message": dbc.Error,
			})
		}
	}

	//set jwt
	if err := module_jwt.SetToken(c, user.Id, os.Getenv("JWT_SECRET"), time.Duration(time.Hour*24)); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "Error signing token",
		})
	}

	// Successful sign-in
	return c.Status(200).JSON(&user)
}

func Login(c *fiber.Ctx) error {
	//get user from body
	var user models.User
	if err := c.BodyParser(&user); err != nil {
		return c.Status(fiber.StatusBadRequest).SendString(err.Error())
	}

	//get user from database
	var dbUser models.User
	database.Connection.Where("email = ?", user.Email).First(&dbUser)

	//check if user exists
	if dbUser.Id == 0 {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"message": "User not found",
		})
	}

	//make sure they use a password
	if !dbUser.HasPassword || len(dbUser.Password) <= 0 {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "User did not sign up with a password",
		})
	}

	//check if password is correct
	if err := module_encryption.ComparePassword(dbUser.Password, user.Password); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "Incorrect password ",
		})
	}

	//check if they need a code to sign in before giving them jwt
	if dbUser.Role == models.Role.User {
		return c.Status(fiber.StatusUnauthorized).JSON(&dbUser)
	}

	//set jwt
	if err := module_jwt.SetToken(c, dbUser.Id, os.Getenv("JWT_SECRET"), time.Duration(time.Hour*24)); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "Error signing token",
		})
	}

	return c.Status(200).JSON(&dbUser)
}

func Logout(c *fiber.Ctx) error {
	//make cookie expired
	cookie := fiber.Cookie{
		Name:     "jwt",
		Value:    "",
		Expires:  time.Now().Add(-time.Hour),
		HTTPOnly: true,
	}

	c.Cookie(&cookie)

	return c.JSON(fiber.Map{
		"message": "Logout successful",
	})
}

func CheckJwt(c *fiber.Ctx) error {
	_, err := module_jwt.CheckToken(c, os.Getenv("JWT_SECRET"))
	if err != nil {
		c.Status(fiber.StatusUnauthorized)
		return c.JSON(fiber.Map{
			"message": "Unverified",
		})
	}

	user, err := jwt.GetUserFromToken(c)

	if err != nil {
		c.Status(fiber.StatusUnauthorized)
		return c.JSON(fiber.Map{
			"message": err.Error(),
		})
	}

	return c.Status(200).JSON(&user)
}
