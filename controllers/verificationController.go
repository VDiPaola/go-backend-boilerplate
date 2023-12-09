package controllers

import (
	"boilerplate/backend/database"
	"boilerplate/backend/helpers/verification"
	"boilerplate/backend/models"
	"os"
	"time"

	"github.com/VDiPaola/go-backend-module/module_helpers/module_jwt"
	"github.com/gofiber/fiber/v2"
)

func VerifyCode(c *fiber.Ctx) error {
	//get code
	code := c.Params("code")
	userId := c.Params("id")

	if len(code) <= 0 {
		c.Status(fiber.StatusBadRequest)
		return c.JSON(fiber.Map{
			"message": "Invalid Code",
		})
	}

	//get user from id
	var user models.User
	if dbc := database.Connection.Where("id = ?", userId).First(&user); dbc.Error != nil {
		c.Status(fiber.StatusBadRequest)
		return c.JSON(fiber.Map{
			"message": dbc.Error,
		})
	}

	//verify code
	if err := verification.VerifyCode(user.Code, code); err != nil {
		c.Status(fiber.StatusBadRequest)
		return c.JSON(fiber.Map{
			"message": err.Error(),
		})
	}

	//invalidate code
	user.Code.ExpiresAt = time.Date(0, 0, 0, 0, 0, 0, 0, time.UTC).Unix()
	//upgrade role
	user.Role = "member"
	//update user
	database.Connection.Updates(&user)

	//set jwt in cookie
	if err := module_jwt.SetToken(c, user.Id, os.Getenv("JWT_SECRET"), time.Duration(time.Hour*24)); err != nil {
		c.Status(fiber.StatusInternalServerError)
		return c.JSON(fiber.Map{
			"message": "Error signing token",
		})
	}

	//return user
	return c.Status(200).JSON(&user)
}

func RequestCode(c *fiber.Ctx) error {
	userId := c.Params("id")

	//get user from id
	var user models.User
	if dbc := database.Connection.Where("id = ?", userId).First(&user); dbc.Error != nil {
		return dbc.Error
	}

	//generate and send
	if err := verification.GenerateCodeAndSend(user, 15); err != nil {
		c.Status(fiber.StatusInternalServerError)
		return c.JSON(fiber.Map{
			"message": err.Error(),
		})
	}

	return nil
}
