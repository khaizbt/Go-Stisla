package main

import (
	"fmt"
	"github.com/getsentry/sentry-go"
	sentrygin "github.com/getsentry/sentry-go/gin"
	"github.com/gin-contrib/sessions"
	"github.com/gin-contrib/sessions/cookie"
	"github.com/joho/godotenv"

	"goshop/middleware"
	"goshop/repository"
	"goshop/route"
	"goshop/service"
	"log"
	"os"

	"github.com/gin-gonic/gin"
	// "html/template"
)

func main() {
	err := godotenv.Load()
	// fmt.Println("masuk", os.Getenv("DB_USER"))
	if err != nil {
		log.Fatalf("Error getting env, not comming through %v", err)
	} else {
		fmt.Println("We are getting the env values")
	}

	err = sentry.Init(sentry.ClientOptions{
		Dsn: os.Getenv("SENTRY_API"),
	})
	if err != nil {
		log.Fatalf("sentry.Init: %s", err)
	}

	userRepo := repository.NewUserRepository()
	userService := service.NewUserService(userRepo)

	secureMiddleware := middleware.SecureMiddleware()

	router := gin.Default()
	store := cookie.NewStore([]byte("ini-secret-api"))
	router.Use(sessions.Sessions("session", store))
	router.Static("/assets", "./web/assets")
	router.Use(secureMiddleware)
	router.Use(sentrygin.New(sentrygin.Options{}))
	route.RouteUser(router, userService)

	router.Run(":8000")
}
