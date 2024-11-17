package main

import (
	"github.com/gin-gonic/gin"
	"log"
	"order-management/config"
	"order-management/handlers"
	"order-management/middleware"
)

func init() {
	config.Init()
	config.InitDB()

	if config.AppConfig.Env != "development" {
		gin.SetMode(gin.ReleaseMode)
	}
}

func main() {
	router := gin.Default()

	router.Use(gin.Recovery())
	if config.AppConfig.Env == "development" {
		router.Use(gin.Logger())
	}

	api := router.Group("/api/v1")
	{
		api.POST("/login", handlers.Login)
		api.POST("/logout", middleware.AuthRequired(), handlers.Logout)

		orders := api.Group("/orders")
		orders.Use(middleware.AuthRequired())
		{
			orders.POST("", handlers.CreateOrder)
			orders.GET("/all", handlers.GetAllOrders)
		}

	}

	serverAddr := config.AppConfig.Server.GetServerAddress()
	log.Printf("Server starting on %s", serverAddr)
	if err := router.Run(serverAddr); err != nil {
		log.Fatalf("Error starting server: %v", err)
	}
}
