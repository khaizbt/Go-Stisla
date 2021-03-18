package route

import (
	"github.com/gin-gonic/gin"
	"goshop/config"
	"goshop/controller"
	"goshop/middleware"
	"goshop/service"
)

func RouteUser(route *gin.Engine, service service.UserService) {
	authService := config.NewServiceAuth()
	userController := controller.NewUserController(service, authService)
	userMiddleware := middleware.AuthMiddlewareUser(authService, service) // middl.AuthMiddlewareManager(authService, service

	//API
	api := route.Group("/api/v1/")
	api.POST("login", userController.Login)
	api.POST("update-account", userMiddleware, userController.UpdateProfile)

	//WEB
	route.LoadHTMLGlob("web/view/**/*")
	route.GET("/login", userController.LoginIndex)
	route.GET("/register", userController.RegisterIndex)
	route.POST("/register", userController.RegisterStore)
	route.POST("/login", userController.LoginBE)
	route.GET("/delete1-session", userController.DeleteSession)
}
