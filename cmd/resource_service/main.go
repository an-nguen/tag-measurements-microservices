package main

import (
	"fmt"

	"github.com/gin-contrib/cors"
	"github.com/gin-contrib/gzip"
	"github.com/gin-gonic/gin"

	"tag-measurements-microservices/internal/resource_service/controllers"
	structures "tag-measurements-microservices/internal/resource_service/structures"
	"tag-measurements-microservices/pkg"
	"tag-measurements-microservices/pkg/datasource"
	"tag-measurements-microservices/pkg/repository"
	"tag-measurements-microservices/pkg/utils"
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
	Repository: repository.RoleRepository{DataSource: userDb},
}

var privilegeController = controllers.PrivilegeController{
	Repository: repository.PrivilegeRepository{DataSource: userDb},
}

func main() {
	gin.SetMode(gin.ReleaseMode)
	router := gin.Default()
	authMiddleware := &utils.JWTAuthMiddleware{
		Secret: appConfig.HmacSecret,
		UserDB: userDb,
	}
	router.Use(gzip.Gzip(gzip.BestCompression))
	corsConfig := cors.DefaultConfig()
	corsConfig.AddAllowHeaders("authorization")
	corsConfig.AddExposeHeaders("authorization")
	corsConfig.AllowOrigins = []string{appConfig.AllowOrigin}
	router.Use(cors.New(corsConfig))

	tagController.ProxyService = new(pkg.ProxyService)
	tagController.ProxyService.CreateProxyService(appConfig)
	tagController.ProxyService.Start()

	api := router.Group("/api")
	ws := router.Group("/ws")
	tagsWS := ws.Group("/tags")
	{
		tagsWS.GET("", tagController.GetLatestTagDetails)
	}
	api.Use(authMiddleware.New())

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
	}
	measurementAPI := api.Group("/measurements")
	{
		measurementAPI.GET("", tagTempDataController.GetMeasurementsByUUID)
		measurementAPI.GET("/csv", tagTempDataController.GetMeasurementsCSVByUUID)
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
	userAPI := api.Group("/user")
	userAPI.Use(authMiddleware.NewWithRole("ADMIN"))
	{
		userAPI.GET("", userController.GetUsers)
		userAPI.GET("/:id", userController.GetUser)
		userAPI.POST("", userController.CreateUser)
		userAPI.PUT("/:id", userController.UpdateUser)
		userAPI.DELETE("/:id", userController.DeleteUser)
	}
	rolesAPI := api.Group("/roles")
	{
		rolesAPI.GET("/token", roleController.GetRolesByToken)
	}
	roleAPI := api.Group("/role")
	roleAPI.Use(authMiddleware.NewWithRole("ADMIN"))
	{
		roleAPI.GET("", roleController.GetRoles)
		roleAPI.POST("", roleController.CreateRole)
		roleAPI.PUT("/:id", roleController.UpdateRole)
		roleAPI.DELETE("/:id", roleController.DeleteRole)
	}
	privilegeAPI := api.Group("/privilege")
	privilegeAPI.Use(authMiddleware.NewWithRole("ADMIN"))
	{
		privilegeAPI.GET("", privilegeController.GetPrivileges)
		privilegeAPI.POST("", privilegeController.CreatePrivilege)
		privilegeAPI.PUT("/:id", privilegeController.UpdatePrivilege)
	}

	err := router.Run(fmt.Sprintf(":%s", appConfig.ServerPort))
	if err != nil {
		panic(err)
	}
}
