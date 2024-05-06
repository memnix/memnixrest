// Code generated by Wire. DO NOT EDIT.

//go:generate go run -mod=mod github.com/google/wire/cmd/wire gen -tags "prod"
//go:build !wireinject
// +build !wireinject

package services

import (
	"github.com/memnix/memnix-rest/app/v2/handlers"
	"github.com/memnix/memnix-rest/infrastructures"
	"github.com/memnix/memnix-rest/services/auth"
	"github.com/memnix/memnix-rest/services/user"
)

// Injectors from wire.go:

func InitializeAuthHandler() handlers.AuthController {
	pool := infrastructures.GetPgxConn()
	iRepository := user.NewRepository(pool)
	iUseCase := auth.NewUseCase(iRepository)
	authController := handlers.NewAuthController(iUseCase)
	return authController
}
