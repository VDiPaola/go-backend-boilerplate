package jwt

import (
	"boilerplate/backend/database"
	"boilerplate/backend/models"
	"os"

	"github.com/VDiPaola/go-backend-module/module_helpers/module_jwt"
	"github.com/gofiber/fiber/v2"
)

func GetUserFromToken(c *fiber.Ctx) (models.User, error) {
	//check token
	token, err := module_jwt.CheckToken(c, os.Getenv("JWT_SECRET"))

	if err != nil {
		return models.User{}, err
	}

	//get claim in correct format
	claims := module_jwt.GetClaims(token)

	var user models.User

	//get user from claim
	database.Connection.Where("id = ?", claims.Issuer).First(&user)

	return user, nil
}
