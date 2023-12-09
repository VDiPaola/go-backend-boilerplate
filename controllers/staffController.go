package controllers

import (
	"boilerplate/backend/database"
	"boilerplate/backend/models"
	"os"
	"time"

	"github.com/VDiPaola/go-backend-module/module_helpers/module_encryption"
	"github.com/VDiPaola/go-backend-module/module_helpers/module_jwt"
	"github.com/gofiber/fiber/v2"
)

func StaffLogin(c *fiber.Ctx) error {
	//get user from body
	var user models.User
	if err := c.BodyParser(&user); err != nil {
		return c.Status(503).SendString(err.Error())
	}

	//get user from database
	var dbUser models.User
	database.Connection.Where("email = ?", user.Email).First(&dbUser)

	//check if user exists
	if dbUser.Id == 0 {
		c.Status(fiber.StatusNotFound)
		return c.JSON(fiber.Map{
			"message": "User not found",
		})
	}

	//check if password is correct
	if err := module_encryption.ComparePassword(dbUser.Password, user.Password); err != nil {
		c.Status(fiber.StatusBadRequest)
		return c.JSON(fiber.Map{
			"message": "Incorrect password ",
		})
	}

	if err := module_jwt.SetToken(c, dbUser.Id, os.Getenv("JWT_SECRET"), time.Duration(time.Hour*24)); err != nil {
		c.Status(fiber.StatusInternalServerError)
		return c.JSON(fiber.Map{
			"message": "Error signing token",
		})
	}

	return c.Status(200).JSON(&dbUser)
}
