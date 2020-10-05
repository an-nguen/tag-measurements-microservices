package main

import (
	"Thermo-WH/pkg"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"

	"Thermo-WH/internal/resource_service/controllers"
	structures "Thermo-WH/internal/resource_service/structures"
	"Thermo-WH/pkg/datasource"
	"Thermo-WH/pkg/repository"
	"Thermo-WH/pkg/utils"
)

var appConfig = structures.ReadAppConfig("/configs/config_resource.json")
var db = datasource.InitDatabaseConnection(appConfig.Host, appConfig.Port, appConfig.User, appConfig.Password, appConfig.DbName)
var userDb = datasource.InitDatabaseConnection(appConfig.Host, appConfig.Port, appConfig.User, appConfig.Password, appConfig.DbNameUsers)

var tagController = controllers.TagController{
	Repository: repository.TagRepository{DataSource: db},
}
var tagManagerController = controllers.TagManagerController{
	Repository: repository.TagManagerRepository{DataSource: db},
}
var tagTempDataController = controllers.MeasurementController{
	Repository: repository.MeasurementRepository{DataSource: db},
}
var signalMeasurementController = controllers.SignalMeasurementController{
	Repository: repository.SignalTagDataRepository{DataSource: db},
}
var voltageMeasurementController = controllers.VoltageMeasurementController{
	Repository: repository.VoltageTagDataRepository{DataSource: db},
}
var temperatureZoneController = controllers.TemperatureZoneController{
	Repository: repository.WarehouseGroupRepository{DataSource: db},
}
var wirelessTagAccountController = controllers.WirelessTagAccountController{
	Repository: repository.WstAccountRepository{DataSource: db},
}
var userController = controllers.UserController{
	Repository: repository.UserRepository{DataSource: userDb},
}
var roleController = controllers.RoleController{
	Secret:     appConfig.HmacSecret,
	Repository: repository.UserRepository{DataSource: userDb},
}

func main() {
	if tz := os.Getenv("TZ"); tz != "" {
		var err error
		time.Local, err = time.LoadLocation(tz)
		if err != nil {
			log.Printf("error loading location '%s': %v\n", tz, err)
		}
	}

	router := gin.Default()
	authMiddleware := &utils.JWTAuthMiddleware{
		Secret: appConfig.HmacSecret,
		UserDB: userDb,
	}
	corsConfig := cors.DefaultConfig()
	corsConfig.AddAllowHeaders("authorization")
	corsConfig.AddExposeHeaders("authorization")
	corsConfig.AllowOrigins = []string{appConfig.AllowOrigin}
	router.Use(cors.New(corsConfig))

	defer db.Close()
	db.DB().SetMaxOpenConns(90)

	api := router.Group("/api")
	api.Use(authMiddleware.New())

	tagController.ProxyService = new(pkg.ProxyService)
	tagController.ProxyService.CreateProxyService(appConfig)
	tagController.ProxyService.Start()

	tagManagersAPI := api.Group("/tagManagers")
	{
		tagManagersAPI.GET("", tagManagerController.GetTagManagers)
		tagManagersAPI.GET("/:id", tagManagerController.GetTagManager)
		tagManagersAPI.PUT("/:id", tagManagerController.UpdateTagManager)
	}
	tagsAPI := api.Group("/tags")
	{
		tagsAPI.GET("", tagController.GetTags)
		tagsAPI.PUT("/:id", tagController.UpdateTag)
		tagsAPI.GET("/latest", tagController.GetLatestTagDetails)
	}
	temperatureTagsAPI := api.Group("/measurements")
	{
		temperatureTagsAPI.GET("", tagTempDataController.GetTempByUUID)
	}
	voltageDataAPI := api.Group("/voltageTagData")
	{
		voltageDataAPI.GET("", voltageMeasurementController.GetVoltageByUUID)
	}
	signalDataAPI := api.Group("/signalTagData")
	{
		signalDataAPI.GET("", signalMeasurementController.GetSignalByUUID)
	}
	temperatureZoneAPI := api.Group("/temperatureZones")
	{
		temperatureZoneAPI.GET("", temperatureZoneController.GetTemperatureZones)
		temperatureZoneAPI.GET("/:id", temperatureZoneController.GetTemperatureZone)
		temperatureZoneAPI.POST("", temperatureZoneController.CreateTemperatureZone).
			Use(authMiddleware.NewWithRole("ADMIN"))
		temperatureZoneAPI.PUT("/:id", temperatureZoneController.UpdateTemperatureZone).
			Use(authMiddleware.NewWithRole("ADMIN"))
		temperatureZoneAPI.DELETE("/:id", temperatureZoneController.DeleteTemperatureZone).
			Use(authMiddleware.NewWithRole("ADMIN"))
	}
	wstAccountsAPI := api.Group("/wstAccounts")
	wstAccountsAPI.Use(authMiddleware.NewWithRole("ADMIN"))
	{
		wstAccountsAPI.GET("", wirelessTagAccountController.GetAccounts)
		wstAccountsAPI.POST("", wirelessTagAccountController.AddAccount)
		wstAccountsAPI.PUT("/:id", wirelessTagAccountController.UpdateAccount)
		wstAccountsAPI.DELETE("/:id", wirelessTagAccountController.DeleteAccount)
	}
	userAPI := api.Group("/users")
	{
		userAPI.POST("", userController.CreateUser)
	}
	roleAPI := api.Group("/roles")
	{
		roleAPI.GET("/token", roleController.GetRolesByToken)
	}

	err := router.Run(fmt.Sprintf(":%s", appConfig.ServerPort))
	if err != nil {
		panic(err)
	}
}
