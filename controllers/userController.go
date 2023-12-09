package controllers

import (
	"boilerplate/backend/database"
	"boilerplate/backend/models"
	"os"

	"github.com/VDiPaola/go-backend-module/module_helpers/module_jwt"
	"github.com/gofiber/fiber/v2"
)

func GetUserFromToken(c *fiber.Ctx) error {
	//check token
	token, err := module_jwt.CheckToken(c, os.Getenv("JWT_SECRET"))

	if err != nil {
		c.Status(fiber.StatusUnauthorized)
		return c.JSON(fiber.Map{
			"message": err,
		})
	}

	//get claim in correct format
	claims := module_jwt.GetClaims(token)

	var user models.User

	//get user from claim
	database.Connection.Where("id = ?", claims.Issuer).First(&user)

	return c.JSON(claims)
}

func GetUsers(c *fiber.Ctx) error {
	//get all users
	var _users []models.User

	database.Connection.Find(&_users)

	return c.JSON(&_users)
}

func GetUser(c *fiber.Ctx) error {
	//get user from id
	id := c.Params("id")
	var _user models.User

	result := database.Connection.Find(&_user, id)

	if result.RowsAffected == 0 {
		return c.SendStatus(404)
	}

	return c.Status(200).JSON(&_user)
}
