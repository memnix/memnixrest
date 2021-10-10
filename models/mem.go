package models

import (
	"time"

	"gorm.io/gorm"
)

// Mem structure
type Mem struct {
	gorm.Model
	UserID     uint `json:"user_id" example:"1"`
	User       User
	CardID     uint `json:"card_id" example:"1"`
	Card       Card
	DeckID     uint `json:"deck_id" example:"1"`
	Deck       Deck
	Quality    uint      `json:"quality" example:"0"` // [0: Blackout - 1: Error with choices - 2: Error with hints - 3: Error - 4: Good with hints - 5: Perfect]
	Repetition uint      `json:"repetition" example:"0" `
	Efactor    float32   `json:"e_factor" example:"2.5"`
	Interval   uint      `json:"interval" example:"0"`
	Total      uint      `json:"total" example:"0"`
	NextDate   time.Time `json:"next_date" example:"06/01/2003" gorm:"autoCreateTime"`
}
