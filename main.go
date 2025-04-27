package main

import (
	"net/http"

	"go-project/configs"
	"go-project/controllers"
	"go-project/middlewares"
	"go-project/utils"

	"github.com/gin-gonic/gin"
)

func main() {
	gin.SetMode(gin.DebugMode)
	r := gin.Default()
	r.StaticFile("/favicon.ico", "./static/favicon.ico")
	r.Use(middlewares.ErrorHandler(), middlewares.ValidatorMiddleware())

	configs.ConnectMysql()

	r.POST("/login", controllers.Login)
	r.POST("/register", controllers.Register)

	v1 := r.Group("/api/v1")
	v1.Use(middlewares.JWTAuthMiddleware())
	{

		v1.GET("/teams", controllers.GetAllTeams)
		v1.GET("/team/:id", controllers.GetTeamByID)

		v1.GET("/books", controllers.GetAllBooks)
		v1.GET("/books/search", controllers.SearchBooks)
		v1.GET("/book/:id", controllers.GetBookByID)
		v1.POST("/book", controllers.CreateBook)
		v1.PUT("book/:id", controllers.UpdateBook)
		v1.DELETE("/book/:id", controllers.DeleteBook)

		v1.GET("/pets", controllers.GetAllPets)

		v1.GET("/albums", controllers.GetAllAlbums)
		v1.GET("albums/search", controllers.SearchAlbums)
		v1.GET("/album/:id", controllers.GetAlbumByID)
		v1.POST("/album", controllers.CreateAlbum)
		v1.PUT("/album/:id", controllers.UpdateAlbum)
		v1.DELETE("/album/:id", controllers.DeleteAlbum)

		v1.GET("/songs", controllers.GetAllSongs)
	}

	r.NoRoute(func(c *gin.Context) {
		utils.ErrorResponse(
			c,
			http.StatusNotFound,
			utils.ErrorNotFound,
			gin.H{"path": c.Request.URL.Path},
		)
	})

	r.Run()
}
