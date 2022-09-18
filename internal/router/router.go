package router

import (
	"context"
	"fmt"
	"net/http"
	"path"

	"github.com/golang-jwt/jwt"
	coreContext "github.com/khotchapan/KonLakRod-api/internal/core/context"
	"github.com/khotchapan/KonLakRod-api/internal/entities"
	guestEndpoint "github.com/khotchapan/KonLakRod-api/internal/handlers/guest"
	tokenEndpoint "github.com/khotchapan/KonLakRod-api/internal/handlers/token"
	coreMiddleware "github.com/khotchapan/KonLakRod-api/internal/middleware"
	"github.com/khotchapan/KonLakRod-api/services/test"
	"github.com/khotchapan/KonLakRod-api/services/user"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

type Options struct {
	App        context.Context
	Collection context.Context
	Echo       *echo.Echo
}

func Router(options *Options) {
	app := options.App
	collection := options.Collection
	e := options.Echo

	//===============================================================================
	// Configure middleware with the custom claims type
	config := middleware.JWTConfig{
		Claims:        &coreContext.Claims{},
		SigningKey:    []byte("secret"),
		SigningMethod: jwt.SigningMethodHS256.Name,
	}
	checkSessionMiddleware := middleware.JWTWithConfig(config)
	requiredAdmin := coreMiddleware.RequiredRoles(entities.TestRole)

	//===============================================================================
	//home
	e.GET(path.Join("/"), Version)
	api := e.Group("/v1/api")
	api.GET("", func(c echo.Context) error {
		return c.String(http.StatusOK, "Hello, World!")
	})

	// Login route
	//api.POST("/login", login)

	// Unauthenticated route
	api.GET("/", accessible)
	// Restricted group
	r := api.Group("/restricted", checkSessionMiddleware, requiredAdmin)
	{
		//api.Use(checkSessionMiddleware)
		r.GET("", restricted)
	}
	guest := guestEndpoint.NewHandler(guestEndpoint.NewService(app, collection))
	guestGroup := api.Group("/guest")
	{
		guestGroup.POST("/login", guest.LoginUsers)
	}
	token := tokenEndpoint.NewHandler(tokenEndpoint.NewService(app, collection))
	tokens := api.Group("/tokens")
	{
		tokens.POST("/refreshToken", token.RefreshToken)
	}

	//user
	users := user.NewHandler(user.NewService(app, collection))
	usersGroup := api.Group("/users")
	usersGroup.GET("/me", users.GetMe, checkSessionMiddleware)
	usersGroup.GET("", users.GetAllUsers)
	usersGroup.GET("/:id", users.GetOneUsers)
	usersGroup.POST("", users.PostUsers)
	usersGroup.PUT("", users.PutUsers)
	usersGroup.DELETE("", users.DeleteUsers)
	usersGroup.POST("/upload", users.UploadFile)
	usersGroup.POST("/image/upload", users.UploadFileUsers)

	// test zone
	testService := test.NewHandler(test.NewService(app, collection))
	testGroup := api.Group("/tests")
	testGroup.GET("/google-cloud/books", testService.GetFile)
	testGroup.GET("/google-cloud/books/:id", testService.GetOneGoogleCloudBooks)
	testGroup.POST("/google-cloud/books", testService.PostGoogleCloudBooks)
	testGroup.PUT("/google-cloud/books", testService.PutBooks)
	testGroup.DELETE("/google-cloud/books", testService.DeleteBooks)
	testGroup.POST("/google-cloud/image/upload", testService.UploadImage)
}
func Version(c echo.Context) error {
	return c.JSON(http.StatusOK, map[string]interface{}{"version": 2.6})
}

// jwtCustomClaims are custom claims extending default ones.
// See https://github.com/golang-jwt/jwt for more examples
type jwtCustomClaims struct {
	Name  string `json:"name"`
	Admin bool   `json:"admin"`
	jwt.StandardClaims
}

func accessible(c echo.Context) error {
	var mystring = jwt.SigningMethodHS256.Name
	var mystring2 = jwt.SigningMethodHS256.Alg()
	fmt.Printf("%s \n %s", mystring, mystring2)
	return c.String(http.StatusOK, "Accessible")
}

func restricted(c echo.Context) error {
	user := c.Get("user").(*jwt.Token)
	claims := user.Claims.(*coreContext.Claims)
	name := claims.Subject
	return c.String(http.StatusOK, "Welcome "+name+"!")
}
