package core

import (
	"memnixrest/database"
	"memnixrest/models"
	"time"

	"github.com/gofiber/fiber/v2"
)

// GetUser
func GetUser(userID uint) models.User {
	db := database.DBConn

	user := new(models.User)

	if err := db.First(&user, userID).Error; err != nil {
		return *user
	}

	return *user
}

// FetchNextTodayCard
func FetchNextTodayCard(c *fiber.Ctx, userID uint, deckID uint) models.ResponseHTTP {
	user := GetUser(userID)
	return FetchNextTodayMemByUserAndDeck(c, &user, deckID)

}

// FetchNextCard
func FetchNextCard(c *fiber.Ctx, userID uint, deckID uint) models.ResponseHTTP {
	user := GetUser(userID)
	return FetchNextMemByUserAndDeck(c, &user, deckID)
}

// UpdateMem
func UpdateMem(c *fiber.Ctx, r *models.Revision, mem *models.Mem) {
	db := database.DBConn

	if r.Result {
		if mem.Repetition == 0 {
			mem.Interval = 1
		} else if mem.Repetition == 1 {
			mem.Interval = 2
		} else if mem.Repetition == 2 {
			mem.Interval = 3
		} else if mem.Repetition == 3 {
			mem.Interval = 6
		} else {
			mem.Interval = uint((float32(mem.Interval) * mem.Efactor)) + 1
		}

		mem.Repetition += 1

		if mem.Repetition >= 3 {
			mem.Quality = 5
		} else {
			mem.Quality = 4
		}
	} else {
		mem.Repetition = 0
		mem.Interval = 0
		if mem.Repetition >= 3 {
			mem.Quality = 3
		} else if mem.Repetition < 3 && mem.Repetition >= 1 {
			mem.Quality = 2
		} else if mem.Total < 2 {
			mem.Quality = 1
		} else {
			mem.Quality = 0
		}
	}

	mem.Efactor = mem.Efactor + (0.1 - (5.0-float32(mem.Quality))*(0.08+(5-float32(mem.Quality)))*0.02)

	if mem.Efactor < 1.3 {
		mem.Efactor = 1.3
	}

	mem.NextDate = time.Now().AddDate(0, 0, int(mem.Interval))

	db.Save(mem)
}