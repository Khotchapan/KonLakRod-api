package main

import (
	"fmt"
	"github.com/joho/godotenv"
	"github.com/khotchapan/project-api/configs"
	"github.com/khotchapan/project-api/routes"
	"github.com/labstack/echo/v4"
	"net/http"
	"os"
	"path"
)

func main() {
	// e := echo.New()
	// e.GET("/", func(c echo.Context) error {
	// 	return c.String(http.StatusOK, "Hello, World!")
	// })
	// e.Logger.Fatal(e.Start(":1323"))
	fmt.Println("Go Program")
	e := echo.New()
	e.GET(path.Join("/"), Version)

	godotenv.Load()
	port := os.Getenv("PORT")
	port = "1323"
	address := fmt.Sprintf("%s:%s", "0.0.0.0", port)
	fmt.Println(address)
	//run database
	configs.ConnectDB()
	//routes
    routes.UserRoute(e) //add this
	e.Start(address)

}
func Version(c echo.Context) error {
	return c.JSON(http.StatusOK, map[string]interface{}{"version": 1})
}
