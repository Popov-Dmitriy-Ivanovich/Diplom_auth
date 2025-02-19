package main

import (
	_ "github.com/Popov-Dmitriy-Ivanovich/Diplom_auth/docs"
	"github.com/Popov-Dmitriy-Ivanovich/Diplom_auth/models"
	"github.com/Popov-Dmitriy-Ivanovich/Diplom_auth/routes"
	"github.com/Popov-Dmitriy-Ivanovich/Diplom_auth/routes/user"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

// @title           Diplom API
// @version         1.0

// @termsOfService  http://swagger.io/terms/

// @contact.name   API Support
// @contact.url    http://www.swagger.io/support
// @contact.email  support@swagger.io

// @license.name  Apache 2.0
// @license.url   http://www.apache.org/licenses/LICENSE-2.0.html

// @host      localhost
// @BasePath  /auth

// @securityDefinitions.basic  BasicAuth

// @externalDocs.description  OpenAPI
// @externalDocs.url          https://swagger.io/resources/open-api/
func main() {
	err := godotenv.Load()
	if err != nil {
		panic(err)
	}
	models.GetDb()

	r := gin.Default()



	apiGroup := r.Group("/auth")
	routes.WriteRoutes(apiGroup, &user.User{})

	apiGroup.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
	apiGroup.Static("/static", "static")

	r.Run()

}