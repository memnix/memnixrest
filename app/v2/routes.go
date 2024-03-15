package v2

import (
	"github.com/labstack/echo/v4"
	"github.com/memnix/memnix-rest/app/v2/handlers"
	"github.com/memnix/memnix-rest/services"
)

func (i *InstanceSingleton) registerStaticRoutes(e *echo.Echo) {
	g := e.Group("/static", StaticAssetsCacheControlMiddleware)
	g.Static("/", "assets/static")
	g.Static("/img", "assets/img")
}

func (i *InstanceSingleton) registerRoutes(e *echo.Echo) {
	serviceContainer := services.DefaultServiceContainer()
	authController := serviceContainer.AuthHandler()
	pageController := handlers.NewPageController()

	e.GET("/", pageController.GetIndex, StaticPageCacheControlMiddleware)
	e.GET("/login", pageController.GetLogin)
	e.GET("/register", pageController.GetRegister)
	e.POST("/register", authController.PostRegister)
	e.POST("/logout", authController.PostLogout)
	e.POST("/login", authController.PostLogin)
	e.POST("/clicked", pageController.PostClicked)
	e.POST("/register/password", authController.ValidatePassword)
}
