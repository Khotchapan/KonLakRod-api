package routes

import (
	"github.com/khotchapan/KonLakRod-api/controllers" //add this
	"github.com/labstack/echo/v4"
)

func UserRoute(e *echo.Echo) {
	//other routes goes here
	e.GET("/user/:userId", controllers.GetAUser)  //add this
	e.POST("/user", controllers.CreateUser)       //add this
	e.PUT("/user/:userId", controllers.EditAUser) //add this
	e.DELETE("/user/:userId", controllers.DeleteAUser) //add this
	e.GET("/users", controllers.GetAllUsers) //add this
}
