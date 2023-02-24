// Code generated by Wire. DO NOT EDIT.

//go:generate go run github.com/google/wire/cmd/wire
//go:build !wireinject
// +build !wireinject

package internal

import (
	"github.com/memnix/memnix-rest/app/http/controllers"
	"github.com/memnix/memnix-rest/infrastructures"
	"github.com/memnix/memnix-rest/internal/auth"
	"github.com/memnix/memnix-rest/internal/kliento"
	"github.com/memnix/memnix-rest/internal/user"
	"github.com/memnix/memnix-rest/pkg/cacheset"
)

// Injectors from wire.go:

func InitializeKliento() controllers.KlientoController {
	client := infrastructures.GetRedisClient()
	iRedisRepository := kliento.NewRedisRepository(client)
	iUseCase := kliento.NewUseCase(iRedisRepository)
	klientoController := controllers.NewKlientoController(iUseCase)
	return klientoController
}

func InitializeUser() controllers.UserController {
	db := infrastructures.GetDBConn()
	iRepository := user.NewRepository(db)
	iUseCase := user.NewUseCase(iRepository)
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
	iUseCase := user.NewUseCase(iRepository)
	jwtController := controllers.NewJwtController(iUseCase)
	return jwtController
}

func InitializeOAuth() controllers.OAuthController {
	db := infrastructures.GetDBConn()
	iRepository := user.NewRepository(db)
	iUseCase := auth.NewUseCase(iRepository)
	cache := cacheset.New()
	oAuthController := controllers.NewOAuthController(iUseCase, cache)
	return oAuthController
}
