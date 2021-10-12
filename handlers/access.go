package handlers

import (
	"memnixrest/database"
	"memnixrest/models"
	"net/http"

	"github.com/gofiber/fiber/v2"
)

// GET

// GetAllAccesses
func GetAllAccesses(c *fiber.Ctx) error {
	db := database.DBConn

	var accesses []models.Access

	if res := db.Joins("User").Joins("Deck").Find(&accesses); res.Error != nil {

		return c.JSON(ResponseHTTP{
			Success: false,
			Message: "Get All accesses",
			Data:    nil,
		})
	}
	return c.JSON(ResponseHTTP{
		Success: true,
		Message: "Get All accesses",
		Data:    accesses,
	})
}

// GetAllAccesses
func GetAccessesByUserID(c *fiber.Ctx) error {
	db := database.DBConn

	userID := c.Params("userID")

	var accesses []models.Access

	if res := db.Joins("User").Joins("Deck").Where("accesses.user_id = ?", userID).Find(&accesses); res.Error != nil {

		return c.JSON(ResponseHTTP{
			Success: false,
			Message: "Get All accesses",
			Data:    nil,
		})
	}
	return c.JSON(ResponseHTTP{
		Success: true,
		Message: "Get All accesses",
		Data:    accesses,
	})

}

// GetAccessByID
func GetAccessByID(c *fiber.Ctx) error {
	id := c.Params("id")
	db := database.DBConn

	access := new(models.Access)

	if err := db.Joins("User").Joins("Deck").First(&access, id).Error; err != nil {
		return c.Status(http.StatusServiceUnavailable).JSON(ResponseHTTP{
			Success: false,
			Message: err.Error(),
			Data:    nil,
		})
	}

	return c.JSON(ResponseHTTP{
		Success: true,
		Message: "Success get access by ID.",
		Data:    *access,
	})
}

// GetAccessByUserAndDeckID
func GetAccessByUserAndDeckID(c *fiber.Ctx) error {
	userID := c.Params("userID")
	deckID := c.Params("deckID")

	db := database.DBConn

	access := new(models.Access)

	if err := db.Joins("User").Joins("Deck").Where("accesses.user_id = ? AND accesses.deck_id = ?", userID, deckID).First(&access).Error; err != nil {
		return c.Status(http.StatusServiceUnavailable).JSON(ResponseHTTP{
			Success: false,
			Message: err.Error(),
			Data:    nil,
		})
	}

	return c.JSON(ResponseHTTP{
		Success: true,
		Message: "Success get access by ID.",
		Data:    *access,
	})
}

// POST

// CreateNewAccess
func CreateNewAccess(c *fiber.Ctx) error {
	db := database.DBConn

	access := new(models.Access)

	if err := c.BodyParser(&access); err != nil {
		return c.Status(http.StatusBadRequest).JSON(ResponseHTTP{
			Success: false,
			Message: err.Error(),
			Data:    nil,
		})
	}

	db.Preload("User").Preload("Deck").Create(access)

	return c.JSON(ResponseHTTP{
		Success: true,
		Message: "Success register an access",
		Data:    *access,
	})
}

// PUT

// UpdateAccessByID
func UpdateAccessByID(c *fiber.Ctx) error {
	db := database.DBConn
	id := c.Params("id")

	access := new(models.Access)

	if err := db.First(&access, id).Error; err != nil {
		return c.Status(http.StatusServiceUnavailable).JSON(ResponseHTTP{
			Success: false,
			Message: err.Error(),
			Data:    nil,
		})
	}

	if err := UpdateAccess(c, access); err != nil {
		return c.Status(http.StatusServiceUnavailable).JSON(ResponseHTTP{
			Success: false,
			Message: "Couldn't update the access",
			Data:    nil,
		})
	}

	return c.JSON(ResponseHTTP{
		Success: true,
		Message: "Success update access by Id.",
		Data:    *access,
	})
}

// UpdateAccess
func UpdateAccess(c *fiber.Ctx, a *models.Access) error {
	db := database.DBConn

	if err := c.BodyParser(&a); err != nil {
		return c.Status(http.StatusBadRequest).JSON(ResponseHTTP{
			Success: false,
			Message: err.Error(),
			Data:    nil,
		})
	}

	db.Preload("User").Preload("Deck").Save(a)

	return nil
}
