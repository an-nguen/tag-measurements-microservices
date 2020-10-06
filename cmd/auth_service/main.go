package main

import (
	"fmt"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"

	"tag-measurements-microservices/internal/auth_service/controllers"
	"tag-measurements-microservices/internal/auth_service/structures"
	"tag-measurements-microservices/pkg/datasource"
	"tag-measurements-microservices/pkg/repository"
)

var appConfig = structures.ReadAppConfig("/configs/config_auth.json")
var db = datasource.InitDatabaseConnection(appConfig.Host, appConfig.Port,
	appConfig.User, appConfig.Password, appConfig.DbName)

//var amqpChannel = datasource.InitAmqpConn(appConfig.AmqpURI)

var authController = controllers.AuthController{UserRepo: repository.UserRepository{DataSource: db},
	Secret: appConfig.HmacSecret}

func main() {
	router := gin.Default()
	corsConfig := cors.DefaultConfig()
	corsConfig.AddAllowHeaders("authorization")
	corsConfig.AddExposeHeaders("authorization")
	corsConfig.AllowOrigins = []string{appConfig.AllowOrigin}
	router.Use(cors.New(corsConfig))
	defer db.Close()

	auth := router.Group("/auth")
	{
		auth.POST("", authController.AuthUser)
		auth.GET("/verify", authController.VerifyToken)
		auth.POST("/refresh", authController.RefreshToken)
	}

	err := router.Run(fmt.Sprintf(":%s", appConfig.ServerPort))
	if err != nil {
		panic(err)
	}
}
