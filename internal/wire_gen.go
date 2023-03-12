// Code generated by Wire. DO NOT EDIT.

//go:generate go run github.com/google/wire/cmd/wire
//go:build !wireinject
// +build !wireinject

package internal

import (
	"github.com/memnix/memnix-rest/app/http/controllers"
	"github.com/memnix/memnix-rest/app/meilisearch"
	"github.com/memnix/memnix-rest/infrastructures"
	"github.com/memnix/memnix-rest/internal/auth"
	"github.com/memnix/memnix-rest/internal/deck"
	"github.com/memnix/memnix-rest/internal/user"
)

// Injectors from wire.go:

func InitializeUser() controllers.UserController {
	db := infrastructures.GetDBConn()
	iRepository := user.NewRepository(db)
	client := infrastructures.GetRedisClient()
	iRedisRepository := user.NewRedisRepository(client)
	iUseCase := user.NewUseCase(iRepository, iRedisRepository)
	userController := controllers.NewUserController(iUseCase)
	return userController
}

func InitializeAuth() controllers.AuthController {
	db := infrastructures.GetDBConn()
	iRepository := user.NewRepository(db)
	iUseCase := auth.NewUseCase(iRepository)
	authController := controllers.NewAuthController(iUseCase)
	return authController
}

func InitializeJWT() controllers.JwtController {
	db := infrastructures.GetDBConn()
	iRepository := user.NewRepository(db)
	client := infrastructures.GetRedisClient()
	iRedisRepository := user.NewRedisRepository(client)
	iUseCase := user.NewUseCase(iRepository, iRedisRepository)
	jwtController := controllers.NewJwtController(iUseCase)
	return jwtController
}

func InitializeOAuth() controllers.OAuthController {
	db := infrastructures.GetDBConn()
	iRepository := user.NewRepository(db)
	iUseCase := auth.NewUseCase(iRepository)
	client := infrastructures.GetRedisClient()
	iAuthRedisRepository := auth.NewRedisRepository(client)
	oAuthController := controllers.NewOAuthController(iUseCase, iAuthRedisRepository)
	return oAuthController
}

func InitializeDeck() controllers.DeckController {
	db := infrastructures.GetDBConn()
	iRepository := deck.NewRepository(db)
	client := infrastructures.GetRedisClient()
	iRedisRepository := deck.NewRedisRepository(client)
	iUseCase := deck.NewUseCase(iRepository, iRedisRepository)
	deckController := controllers.NewDeckController(iUseCase)
	return deckController
}

func InitializeMeiliSearch() meilisearch.MeiliSearch {
	db := infrastructures.GetDBConn()
	iRepository := deck.NewRepository(db)
	client := infrastructures.GetRedisClient()
	iRedisRepository := deck.NewRedisRepository(client)
	iUseCase := deck.NewUseCase(iRepository, iRedisRepository)
	meiliSearch := meilisearch.NewMeiliSearch(iUseCase)
	return meiliSearch
}
